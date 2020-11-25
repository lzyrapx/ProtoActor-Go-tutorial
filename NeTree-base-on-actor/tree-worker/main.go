package main

import(
	"flag"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/lzyrapx/ProtoActor-Go-tutorial/NeTree-base-on-actor/messages"
	"log"
	"strconv"
	"sync"
	"time"
	"errors"
)

type WorkerActor struct {
	waitGroup *sync.WaitGroup
}

func (worker WorkerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.CreateNewTreeResponse:
		log.Printf("Created tree with token: %v and id: %v\n", msg.Token, msg.Id)
	case *messages.TreeTokenOrIDInvalid:
		log.Printf("The input token or id is invalid!\n")
	case *messages.SuccessFindValue:
		log.Printf("Found key: %v with value: %v\n", msg.Key, msg.Value)
	case *messages.TraverseResponse:
		log.Printf("Traversal tree is: %v\n", msg.KvPair)
	case *messages.ErrorFindingValue:
		log.Printf("Key could not be found for find operation!\n")
	case *messages.CouldNotFindKey:
		log.Printf("Key could not be found!\n")
	case *messages.SuccessDeleteKey:
		log.Printf("Successfully deleted key: %v with value: %v from tree\n", msg.Key, msg.Value)
	case *messages.ErrorKeyDoesNotExist:
		log.Printf("The key: %v does not exist!\n", msg.Key)
	case *messages.SuccessfulTreeDelete:
		log.Printf("Successfully deleted tree with token: %v and id: %v\n", msg.Token, msg.Id)
	case *messages.SuccessAdd:
		log.Printf("SuccessfSetParentIDully added key: %v with value: %v to tree", msg.Key, msg.Value)
	case *messages.SendBackTreeResponse:
		log.Printf("[special] send back tree response: %v\n", msg.Trees)
		for _, tree := range msg.Trees {
			log.Printf("tree id: %v\n", tree.Id)
			log.Printf("tree token: %v\n", tree.Token)
		}
	}
}

func main()  {
	log.Printf("net tree worker started\n")

	// remote actor
	flagRemoteHostName := flag.String("remote hostname", "localhost", "the remote hostname to connect")
	flagRemotePort := flag.Int("remote port", 8080, "the remote port to connect")

	// local actor
	flagBindHostName := flag.String("local bind hostname", "localhost", "bind to local hostname")
	flagBindPort := flag.Int("local bind port", 8081, "bind to local port")

	// tree operation
	flagCreateTree := flag.Bool("createtree", false, "create tree. default not creating tree")
	flagInsertValue := flag.Bool("insertvalue", false, "flag for inserting a value to the tree")
	flagFindValue := flag.Bool("findvalue", false, "flag for find a value in the tree")
	flagDeleteKey := flag.Bool("deletekey", false, "flag for deleting a key/value in the tree")
	flagTraverseTree := flag.Bool("traversetree", false, "flag for traversing the tree")
	flagDeleteTree := flag.Bool("deletetree", false, "flag for deleting a tree")

	// 如果你想对一棵树进行操作，必须要有 token 和 id
	flagToken := flag.String("token", "", "flag for token. necessary for all operations")
	flagID := flag.Int("id", 0, "flag for id. necessary for all operations")

	//
	flagKey := flag.Int("key", 0, "key when inserting/deleting/finding values")
	flagValue := flag.String("value", "", "value when inserting a key")
	flagLeafSize := flag.Int("leafsize", 0, "leafsize value when create a tree")

	flag.Parse()

	system := actor.NewActorSystem()
	config := remote.Configure(*flagBindHostName, *flagBindPort)
	remote := remote.NewRemote(system, config)
	remote.Start()

	var wg sync.WaitGroup

	props := actor.PropsFromProducer(func() actor.Actor {
		wg.Add(1)
		return &WorkerActor{&wg}
	})
	context := system.Root
	pid := context.Spawn(props)

	var msg interface{}
	switch {
	case *flagCreateTree:
		if *flagLeafSize < 2 {
			panic(errors.New("leaf size must be greater than 1"))
		}
		msg = &messages.CreateNewTreeForCLI{
			LeafSize: int32(*flagLeafSize),
		}
	case *flagInsertValue:
		msg = &messages.InsertCLI{
			Id: int32(*flagID),
			Token: *flagToken,
			Key: int32(*flagKey),
			Value: *flagValue,
		}
	case *flagFindValue:
		msg = &messages.SearchCLI{
			Id: int32(*flagID),
			Token: *flagToken,
			Key: int32(*flagKey),
		}
	case *flagDeleteKey:
		msg = &messages.DeleteCLI{
			Id: int32(*flagID),
			Token: *flagToken,
			Key: int32(*flagKey),
		}
	case *flagTraverseTree:
		msg = &messages.TraverseCLI{
			Id:int32(*flagID),
			Token: *flagToken,
		}
	case *flagDeleteTree:
		msg = &messages.DeleteTreeCLI{
			Id: int32(*flagID),
			Token: *flagToken,
		}
	default:
		msg = &messages.Trees{}
	}
	remoteAddress := *flagRemoteHostName + ":" + strconv.Itoa(*flagRemotePort)
	remotePid, err := remote.SpawnNamed(remoteAddress,"remote", "hello", 5*time.Second)
	if err != nil {
		log.Printf("could not connect to remote actor tree server")
		panic(err)
	}
	// send message from worker to serer
	context.RequestWithCustomSender(remotePid.Pid, msg, pid)
	wg.Wait()
}
