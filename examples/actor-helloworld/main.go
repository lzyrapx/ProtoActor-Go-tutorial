package main

import (
	"fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type hello struct{
	Who string
}

// 自己加的
type world struct {
	what string
}

type helloActor struct{}

func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		fmt.Printf("Hello %v\n", msg.Who)
	case *world:
		fmt.Printf("%v world\n", msg.what)
	}

}

func main() {
	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(func() actor.Actor { return &helloActor{} })

	pid := system.Root.Spawn(props)
	system.Root.Send(pid, &hello{Who: "Roger"})
	system.Root.Send(pid, &hello{
		"very good",
	})
	system.Root.Send(pid, &world{
		"big",
	})
	time.Sleep(3 * time.Second)
	//_, _ = console.ReadLine()
}
