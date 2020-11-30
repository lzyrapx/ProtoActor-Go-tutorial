package main

import (
	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/lzyrapx/ProtoActor-Go-tutorial/examples/actor-cluster-metrics-consul/messages"
	"log"
	"fmt"
	"time"
)

func setupLogger(c *cluster.Cluster) {
	system := c.ActorSystem
	// Subscribe
	system.EventStream.Subscribe(func(event interface{}) {
		switch msg := event.(type) {
		case *cluster.MemberJoinedEvent:
			log.Printf("Member Joined " + msg.Name())
		case *cluster.MemberLeftEvent:
			log.Printf("Member Left " + msg.Name())
		case *cluster.MemberRejoinedEvent:
			log.Printf("Member Rejoined " + msg.Name())
		case *cluster.MemberUnavailableEvent:
			log.Printf("Member Unavailable " + msg.Name())
		case *cluster.MemberAvailableEvent:
			log.Printf("Member Available " + msg.Name())
		case cluster.TopologyEvent:
			log.Printf("Cluster Topology Poll")
		}
	})
}

func doRequests(c *cluster.Cluster, callopts *cluster.GrainCallOptions) {
	msg := &messages.HelloRequest{Name: "GAM"}
	helloGrain := messages.GetHelloGrainClient(c, "abc")
	// with default callopts
	resp, err := helloGrain.SayHello(msg)
	if err != nil {
		log.Fatalf("SayHello failed. err:%v", err)
	}

	// with custom callopts
	resp, err = helloGrain.SayHello(msg, callopts)
	if err != nil {
		log.Fatalf("SayHello failed. err:%v", err)
	}
	log.Printf("Message from SayHello: %v", resp.Message)
	for i := 0; i < 10000; i++ {
		grainId := fmt.Sprintf("hello%v", i)
		x := messages.GetHelloGrainClient(c, grainId)
		x.SayHello(&messages.HelloRequest{Name: grainId})
	}
	log.Println("Done")
}

func doRequestsAsync(c *cluster.Cluster, callopts *cluster.GrainCallOptions) {
	// sorry, golang has not magic, just use goroutine.
	go func() {
		doRequests(c, callopts)
	}()
}

func main() {
	system := actor.NewActorSystem()
	config := remote.Configure("localhost", 0)

	provider, _ := consul.New()
	clusterConfig := cluster.Configure("my-cluster", provider, config)
	c := cluster.New(system, clusterConfig)
	setupLogger(c)
	c.Start()

	callopts := cluster.NewGrainCallOptions(c).WithTimeout(5 * time.Second).WithRetry(5)
	doRequests(c, callopts)
	doRequestsAsync(c, callopts)
	console.ReadLine()
}
