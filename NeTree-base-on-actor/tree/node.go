package tree

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/lzyrapx/ProtoActor-Go-tutorial/NeTree-base-on-actor/messages"
	"log"
	"sort"
	"time"
)

// 每个节点的 actor
type NodeActor struct {
	Parent, Left, Right  *actor.PID
	LeftMaxKey, LeafSize int32
	Values               map[int32]string
}

func (node *NodeActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Add: // 在树上增加节点(key-value)
		err := node.add(msg, context)
		if err != nil {
			// TODO: send failed message back
		}
		// TODO: send success message back
	case *messages.Find: // 根据 key 在树上找节点
		node.find(msg, context)
	case *messages.Traverse: // 遍历一棵树
		node.traverse(context)
	case *messages.DeleteKey: // 根据 key 删除树上的节点
		node.delete(msg, context)
	// 以下函数通过 actor 供自己调用
	case *messages.SetParentID:
		node.setParentID(context)
	case *messages.Cleanup:
		node.cleanup(msg, context)
	case *messages.FinalCleanup:
		node.finalCleanup(msg, context)
	case *messages.SetParentIDCustom:
		node.Parent.Address = msg.Address
		node.Parent.Id = msg.Pid
	case *messages.InfoSetLeftMaxKey:
		context.Send(context.Sender(), &messages.InfoSetLeftMaxKeyResponse{Key: node.LeftMaxKey})
	case *messages.SetLeftMaxKey:
		log.Printf("Setting largest key of left subtree to %v\n", msg.Key)
		node.LeftMaxKey = msg.Key
	case *messages.SetNewRoot:
		node.Parent = nil
		context.Self().Address = msg.Address
		context.Self().Id = msg.Pid
		context.Stop(actor.NewPID(msg.CustomAddressToDelete, msg.CustomPidToDelete))
	}
}

func (node *NodeActor) add(msg *messages.Add, context actor.Context) error {
	if node.Left == nil {
		log.Printf("Node is a leaf!\n")
		if int32(len(node.Values)) < node.LeafSize {
			node.Values[msg.Key] = msg.Value
			log.Printf("Key: %v was successfully added!\n", msg.Key)
			context.Send(context.Sender(), &messages.SuccessAdd{Key: msg.Key, Value: msg.Value})
		} else {
			log.Printf("Leaf is full, we need to split!\n")
			leftHalf, rightHalf := sortByKeyAndSplitMap(node.Values)
			log.Printf("LeftHalf: %v, RightHalf: %v\n", leftHalf, rightHalf)

			props := actor.PropsFromProducer(func() actor.Actor {
				return &NodeActor{
					Left:     nil,
					Right:    nil,
					LeafSize: node.LeafSize,
					Values:   leftHalf,
				}
			})
			node.Left = context.Spawn(props)
			props = actor.PropsFromProducer(func() actor.Actor {
				return &NodeActor{
					Left:     nil,
					Right:    nil,
					LeafSize: node.LeafSize,
					Values:   rightHalf,
				}
			})
			node.Right = context.Spawn(props)
			context.Send(node.Left, &messages.SetParentID{})
			context.Send(node.Right, &messages.SetParentID{})
			node.LeftMaxKey = lastKeyFromMap(leftHalf)
			if msg.Key <= node.LeftMaxKey {
				log.Printf("key-value pair: %v:%v will be send to the left subtree!\n", msg.Key, msg.Value)
				context.Send(node.Left, msg)
			} else {
				log.Printf("key-value pair: %v:%v will be send to the right subtree!\n", msg.Key, msg.Value)
				context.Send(node.Right, msg)
			}
			node.Values = nil
		}
	} else {
		log.Printf("Node is NOT a leaf!\n")
		if msg.Key <= node.LeftMaxKey {
			context.Send(node.Left, msg)
		} else {
			context.Send(node.Right, msg)
		}
	}
	log.Printf("LeafSize: %v, Key we wanted to add: %v, Current Values: %v\n", node.LeafSize, msg.Key, node.Values)
	return nil
}

func (node *NodeActor) delete(msg *messages.DeleteKey, context actor.Context) {
	if node.Values != nil {
		if node.Parent == nil {
			log.Printf("Trying to delete from root element!\n")
			val, ok := node.Values[msg.Key]
			if !ok {
				log.Printf("Nothing will be deleted from root as specific key: %v does not exist!\n", msg.Key)
				context.Send(context.Sender(), &messages.CouldNotFindKey{})
			} else {
				log.Printf("Deleting element from root with key: %v\n", msg.Key)
				context.Send(context.Sender(), &messages.SuccessDeleteKey{
					Key:   msg.Key,
					Value: val,
				})
				delete(node.Values, msg.Key)
			}
		} else {
			log.Printf("Reached leaf for delete of key: %v\n", msg.Key)
			_, ok := node.Values[msg.Key]
			if !ok {
				log.Printf("Key: %v does not exist\n", msg.Key)
				context.Send(context.Sender(), &messages.ErrorKeyDoesNotExist{Key: msg.Key})
				// Send Error message back
			} else {
				log.Printf("Found Key: %v for deletion\n", msg.Key)
				delete(node.Values, msg.Key)
				context.Send(context.Sender(), &messages.SuccessDeleteKey{
					Key:   msg.Key,
					Value: node.Values[msg.Key],
				})
				if len(node.Values) == 0 {
					log.Printf("Key: %v was deleted but values are empty!\n", msg.Key)
					if msg.IsLeft {
						context.Send(node.Parent, &messages.Cleanup{WasLeft: true})
					} else {
						context.Send(node.Parent, &messages.Cleanup{WasLeft: false})
					}
				} else if msg.IsLeft {
					context.Send(node.Parent, &messages.SetLeftMaxKey{Key: lastKeyFromMap(node.Values)})
				}
			}
		}
	} else {
		log.Printf("Got to move further down to find the correct key\n")
		if msg.Key <= node.LeftMaxKey {
			log.Printf("Going down left!\n")
			context.Send(node.Left, &messages.DeleteKey{
				IsLeft: true,
				Key:    msg.Key,
			})
		} else {
			log.Printf("Going down right!\n")
			context.Send(node.Right, &messages.DeleteKey{
				IsLeft: false,
				Key:    msg.Key,
			})
		}
	}
}

func (node *NodeActor) setParentID(context actor.Context) {
	log.Printf("Setting Parent ID of: %v to: %v\n", context.Self(), context.Parent())
	node.Parent = context.Parent()
}

func (node *NodeActor) traverse(context actor.Context) {
	log.Printf("Traversing through the tree!\n")
	if node.Left == nil {
		log.Printf("Current node is a leaf! Sort values and respond!\n")
		context.Respond(sortMap(node.Values))
	} else {
		log.Printf("Current node is not a leaf! Recursively iterating through the tree and waiting for a response\n")
		left, _ := context.RequestFuture(node.Left, &messages.Traverse{}, 1*time.Second).Result()
		right, _ := context.RequestFuture(node.Right, &messages.Traverse{}, 1*time.Second).Result()

		leftHalf := left.(*messages.TraverseResponse)
		rightHalf := right.(*messages.TraverseResponse)

		leftHalf.KvPair = append(leftHalf.KvPair, rightHalf.KvPair...)

		log.Printf("Sending a response back to the initial caller!")
		context.Respond(leftHalf)
	}
}

func (node *NodeActor) find(msg *messages.Find, context actor.Context) {
	if node.Values == nil {
		if msg.Key <= node.LeftMaxKey {
			context.Send(context.Sender(), &messages.LookingForKeyLeft{})
			context.Send(node.Left, msg)
		} else {
			context.Send(node.Right, msg)
			context.Send(context.Sender(), &messages.LookingForKeyRight{})
		}
	} else {
		val, ok := node.Values[msg.Key]
		if ok {
			context.Send(context.Sender(), &messages.SuccessFindValue{
				Key:   msg.Key,
				Value: val,
			})
		} else {
			context.Send(context.Sender(), &messages.ErrorFindingValue{})
		}
	}
}

func (node *NodeActor) finalCleanup(msg *messages.FinalCleanup, context actor.Context) {
	if msg.SetToLeft {
		node.Left.Address = msg.Address
		node.Left.Id = msg.Pid
	} else {
		node.Right.Address = msg.Address
		node.Right.Id = msg.Pid
	}
	context.Stop(context.Sender())
}

func (node *NodeActor) cleanup(msg *messages.Cleanup, context actor.Context) {
	if node.Parent == nil {
		if msg.WasLeft {
			context.Send(node.Right, &messages.SetNewRoot{
				Pid:                   context.Self().Id,
				Address:               context.Self().Address,
				CustomPidToDelete:     "1234",
				CustomAddressToDelete: "2345",
			})
		} else {
			context.Send(node.Left, &messages.SetNewRoot{
				Pid:                   context.Self().Id,
				Address:               context.Self().Address,
				CustomPidToDelete:     "1234",
				CustomAddressToDelete: "2345",
			})
		}
	} else {
		if msg.WasLeft {
			context.Send(node.Right, &messages.SetParentIDCustom{
				Pid:     node.Parent.Id,
				Address: node.Parent.Address,
			})
			context.Send(node.Parent, &messages.FinalCleanup{Pid: node.Right.Id, Address: node.Right.Address, SetToLeft: false})
		} else {
			context.Send(node.Left, &messages.SetParentIDCustom{
				Pid:     node.Parent.Id,
				Address: node.Parent.Address,
			})
			context.Send(node.Parent, &messages.FinalCleanup{Pid: node.Left.Id, Address: node.Left.Address, SetToLeft: true})
		}
		context.Stop(context.Sender())
	}
}

func sortByKeyAndSplitMap(input map[int32]string) (map[int32]string, map[int32]string) {
	resultKeys := make([]int, 0)

	for key := range input {
		resultKeys = append(resultKeys, int(key))
	}
	sort.Ints(resultKeys)
	leftHalf := make(map[int32]string)
	rightHalf := make(map[int32]string)

	for i := 0; i < len(resultKeys)/2; i++ {
		leftHalf[int32(resultKeys[i])] = input[int32(resultKeys[i])]
	}

	for i := len(resultKeys) / 2; i < len(resultKeys); i++ {
		rightHalf[int32(resultKeys[i])] = input[int32(resultKeys[i])]
	}

	return leftHalf, rightHalf
}

func sortMap(input map[int32]string) *messages.TraverseResponse {
	resultKeys := make([]int, 0)

	for key := range input {
		resultKeys = append(resultKeys, int(key))
	}
	sort.Ints(resultKeys)

	resultMap := make([]*messages.KeyValue, 0)

	for i := 0; i < len(resultKeys); i++ {
		resultMap = append(resultMap, &messages.KeyValue{
			Key:   int32(resultKeys[i]),
			Value: input[int32(resultKeys[i])],
		})
	}
	return &messages.TraverseResponse{KvPair: resultMap}
}

func lastKeyFromMap(input map[int32]string) int32 {
	resultKeys := make([]int32, len(input))
	for key := range input {
		resultKeys = append(resultKeys, key)
	}
	return resultKeys[len(resultKeys)-1]
}
