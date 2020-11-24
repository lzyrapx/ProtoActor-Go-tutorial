package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"flag"
	server "github.com/lzyrapx/ProtoActor-Go-tutorial/NeTree-base-on-actor/tree-server/server"
	"log"
	"sync"
)

func main()  {
	log.Printf("net tree server started\n")
	flagHost := flag.String("host", "localhost", "hostname allow remote actor to connect")
	flagPort := flag.Int("port", 8080, "port allow remote actor to connect")
	flagServiceName := flag.String("serverName", "Tree-Server", "the name of server")

	flag.Parse()

	system := actor.NewActorSystem()
	config := remote.Configure(*flagHost, *flagPort)
	remote := remote.NewRemote(system, config)
	remote.Start()

	context := system.Root
	props := actor.PropsFromProducer(func() actor.Actor {
		return &server.NodeService{
			Trees: make(map[int32]server.TreeItem),
			NextID: 19937,
		}
	})
	pid, err := context.SpawnNamed(props, *flagServiceName)

	remote.Register("hello", props)
	var waitGroup sync.WaitGroup
	if err != nil {
		log.Printf("Something went wrong with creating the named actor\n")
	}
	log.Printf("Successfully created actor with pid %v\n", pid)
	waitGroup.Add(1)
	waitGroup.Wait()
}
