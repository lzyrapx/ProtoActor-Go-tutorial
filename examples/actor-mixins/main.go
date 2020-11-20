package main

import (
	"fmt"
	console "github.com/AsynkronIT/goconsole"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/actor/middleware"
	"github.com/AsynkronIT/protoactor-go/plugin"
)

type myActor struct {
	NameAwareHolder
}

func (state *myActor) Receive(context actor.Context) {
	switch context.Message().(type) {
	case *actor.Started:
		// this actor have been initialized by the receive pipeline
		fmt.Printf("My name is %v\n", state.name)
	}
}

type NameAware interface {
	SetName(name string)
}

type NameAwareHolder struct {
	name string
}

func (state *NameAwareHolder) SetName(name string) {
	state.name = name
}

type NamerPlugin struct{}

// OnStart 和 OnOtherMessage 都是 protoactot-go plugin 里已经实现的
func (p *NamerPlugin) OnStart(ctx actor.ReceiverContext) {
	if p, ok := ctx.Actor().(NameAware); ok {
		p.SetName("Proto.Actor")
	}
}
func (p *NamerPlugin) OnOtherMessage(ctx actor.ReceiverContext, env *actor.MessageEnvelope) {}

func main() {
	system := actor.NewActorSystem()
	rootContext := system.Root
	props := actor.
		PropsFromProducer(func() actor.Actor { return &myActor{} }).
		WithReceiverMiddleware(
			plugin.Use(&NamerPlugin{}),
			middleware.Logger,
		)

	pid := rootContext.Spawn(props)
	rootContext.Send(pid, "bar")
	_, _ = console.ReadLine()
}
