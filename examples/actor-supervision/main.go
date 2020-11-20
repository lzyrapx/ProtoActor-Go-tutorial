package main

import (
	"fmt"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type hello struct{ Who string }
type world struct {
	What string
}
type parentActor struct{}

// parent actor
func (state *parentActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		props := actor.PropsFromProducer(newChildActor)
		child := context.Spawn(props)
		context.Send(child, msg)
	case *world:
		props := actor.PropsFromProducer(newChildActor)
		child := context.Spawn(props)
		context.Send(child, msg)
	}
}

func newParentActor() actor.Actor {
	return &parentActor{}
}

type childActor struct{}

// child actor
func (state *childActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started: // actor 自带的类型
		fmt.Println("Starting, initialize actor here")
	case *actor.Stopping: // actor 自带的类型
		fmt.Println("Stopping, actor is about to shut down")
	case *actor.Stopped: // actor 自带的类型
		fmt.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting: // actor 自带的类型
		fmt.Println("Restarting, actor is about to restart")
	case *hello: // 自己定义的类型
		fmt.Printf("Hello %v\n", msg.Who)
		time.Sleep(3 * time.Second)
	case *world: // 这里需要在 parent actor 里实现
		fmt.Printf("World %v\n", msg.What)
		panic("Ouch") // panic 会造成触发 actor.stopping 和 actor.stopped
	}
}

func newChildActor() actor.Actor {
	return &childActor{}
}

func main() {
	system := actor.NewActorSystem()
	decider := func(reason interface{}) actor.Directive {
		fmt.Println("handling failure for child")
		return actor.StopDirective
	}
	supervisor := actor.NewOneForOneStrategy(10, 1000, decider)
	rootContext := system.Root
	props := actor.
		PropsFromProducer(newParentActor).
		WithSupervisor(supervisor)

	pid := rootContext.Spawn(props)
	rootContext.Send(pid, &hello{Who: "Roger"})

	rootContext.Send(pid, &world{
		What: "google",
	})
	_, _ = console.ReadLine()
}
