package main

import (
	"fmt"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type hello struct{ Who string }
type helloActor struct{}

func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started: // actor 自带的结构体
		fmt.Println("Started, initialize actor here")
	case *actor.Stopping: // actor 自带的结构体
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Stopped: // actor 自带的结构体
		fmt.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting: // actor 自带的结构体
		fmt.Println("Restarting, actor is about restart")
	case *hello: // 自己定义的结构体
		for i := 0; i < 100; i++ {
			fmt.Println("i = ", i)
			time.Sleep(1 * time.Second)
		}
		fmt.Printf("Hello %v\n", msg.Who)
	}
}

func main() {
	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(func() actor.Actor { return &helloActor{} })
	pid := system.Root.Spawn(props)
	system.Root.Send(pid, &hello{Who: "Roger"}) // 执行完这个才会往下运行

	// why wait?
	// Stop is a system message and is not processed through the user message mailbox
	// thus, it will be handled _before_ any user message
	// we only do this to show the correct order of events in the console
	time.Sleep(1 * time.Second)
	system.Root.Stop(pid)

	_, _ = console.ReadLine()
}
