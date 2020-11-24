package server

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/lzyrapx/ProtoActor-Go-tutorial/NeTree-base-on-actor/messages"
	"github.com/lzyrapx/ProtoActor-Go-tutorial/NeTree-base-on-actor/tree"

	"log"
)

type NodeService struct {
	Trees  map[int32]TreeItem
	NextID int32
}

// 树上节点定义体
type TreeItem struct {
	token string
	pid   *actor.PID
	id    int32
}

func (treeService *NodeService) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Trees: // 返回所有存储的树
		log.Printf("Currently available trees: \n")
		response := make([]*messages.TreesResponse, 0)
		for _, v := range treeService.Trees {
			response = append(response, &messages.TreesResponse{
				Token: v.token,
				Id:    v.id,
			})
		}
		context.Respond(&messages.SendBackTreeResponse{Trees: response})
	case *messages.CreateNewTreeForCLI: // 创建新树
		log.Printf("New tree-actor(root) will be created!")
		props := actor.PropsFromProducer(func() actor.Actor {
			return &tree.NodeActor{
				Parent:     nil,
				Left:       nil,
				Right:      nil,
				LeftMaxKey: 0,
				LeafSize:   msg.LeafSize,
				Values:     make(map[int32]string),
			}
		})
		pid := context.Spawn(props)
		newTree := TreeItem{
			token: CreateToken(8), // 创建 8 位 token
			pid: pid,
			id: treeService.NextID,
		}
		treeService.Trees[newTree.id] = newTree
		context.Respond(&messages.CreateNewTreeResponse{ // 向 worker 返回新创建的树的 id 和 token
			Id:    newTree.id,
			Token: newTree.token,
		})
		log.Printf("Tree with id: %v and token: %v was successfully created!\n", newTree.id, newTree.token)
		treeService.NextID++
	case *messages.InsertCLI: // 根据树的 id 向一棵树插入新的 key-value
		log.Printf("Trying to insert key: %v with value: %v\n", msg.Key, msg.Value)
		treeItem, ok := treeService.Trees[msg.Id]
		if !ok {
			context.Respond(&messages.TreeTokenOrIDInvalid{Id: msg.Id, Token: msg.Token})
		} else if CheckIDAndToken2(ok, msg.Id, msg.Token, treeItem) { // 检查 id 和 token
			log.Printf("Adding key: %v with value: %v to the tree!\n", msg.Key, msg.Value)
			pid := treeItem.pid
			// server sends the message to the actual Node Actor by pid
			context.RequestWithCustomSender(pid, &messages.Add{Key: msg.Key, Value: msg.Value}, context.Sender())
			// sending success message back
			// TODO: 需要 node actor 创建成功返回后在 send success message back or fail message back
		} else {
			context.Respond(&messages.TreeTokenOrIDInvalid{
				Id: msg.Id,
				Token: msg.Token,
			})
		}
	case *messages.DeleteCLI: // 根据 key 删除树上的一个节点
		log.Printf("Trying to delete key: %v\n", msg.Key)
		treeItem, ok := treeService.Trees[msg.Id]
		if !ok {
			context.Respond(&messages.TreeTokenOrIDInvalid{Id: msg.Id, Token: msg.Token})
		} else if CheckIDAndToken2(ok, msg.Id, msg.Token, treeItem) {
			log.Printf("Deleting key: %v!\n", msg.Key)
			pid := treeItem.pid
			// server sends the message to the actual NodeActor
			context.RequestWithCustomSender(pid, &messages.DeleteKey{Key: msg.Key}, context.Sender())
			// Currently missing sending success message back
			// TODO: 需要 node actor 删除成功返回后 send success message back or fail message back
		} else {
			context.Respond(&messages.TreeTokenOrIDInvalid{Id: msg.Id, Token: msg.Token})
		}
	case *messages.SearchCLI: // 根据 key 在树上找节点
		log.Printf("Trying to find a key: %v\n", msg.Key)
		treeItem, ok := treeService.Trees[msg.Id]
		if !ok {
			context.Respond(&messages.TreeTokenOrIDInvalid{Id: msg.Id, Token: msg.Token})
		} else if CheckIDAndToken2(ok, msg.Id, msg.Token, treeItem) {
			log.Printf("Looking for key: %v now!\n", msg.Key)
			pid := treeItem.pid
			// server sends the message to the actual NodeActor
			context.RequestWithCustomSender(pid, &messages.Find{Key: msg.Key}, context.Sender())
			// Currently sending success message back
			// TODO: 需要 node actor 删除节点成功后返回 send success message back or fail message back
		} else {
			context.Respond(&messages.TreeTokenOrIDInvalid{Id: msg.Id, Token: msg.Token})
		}
	case *messages.TraverseCLI: // 遍历整棵树
		log.Printf("Trying to traverse through a tree!\n")
		treeItem, ok := treeService.Trees[msg.Id]
		if !ok {
			context.Respond(&messages.TreeTokenOrIDInvalid{Id: msg.Id, Token: msg.Token})
		} else if CheckIDAndToken2(ok, msg.Id, msg.Token, treeItem) {
			log.Printf("Traversing through the tree with id: %v now!\n", msg.Id)
			pid := treeItem.pid
			// server sends the message to the actual NodeActor
			context.RequestWithCustomSender(pid, &messages.Traverse{}, context.Sender())
			// Currently missing sending success message back
			// TODO: 需要 node actor 遍历树成功后返回 send success message back or fail message back
		} else {
			context.Respond(&messages.TreeTokenOrIDInvalid{Id: msg.Id, Token: msg.Token})
		}
	case *messages.DeleteTreeCLI: // 删除一棵树
		log.Printf("Trying to delete a tree!\n")
		treeItem, ok := treeService.Trees[msg.Id]
		if !ok {
			context.Respond(&messages.TreeTokenOrIDInvalid{Id: msg.Id, Token: msg.Token})
		} else if CheckIDAndToken2(ok, msg.Id, msg.Token, treeItem) {
			log.Printf("Deleting the tree with id: %v and token: %v now!\n", msg.Id, msg.Token)
			pid := treeItem.pid
			context.Stop(pid)
			delete(treeService.Trees, msg.Id)
			// NodeService directly responds to the remote actor
			context.Respond(&messages.SuccessfulTreeDelete{
				Id:    msg.Id,
				Token: msg.Token,
			})
			// Currently missing sending success message back
			// TODO: 需要 node actor 删除树成功后返回 send success message back or fail message back
		} else {
			context.Respond(&messages.TreeTokenOrIDInvalid{Id: msg.Id, Token: msg.Token})
		}
	}
}
