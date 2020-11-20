package main

import (
	"fmt"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type Hello struct{ Who string }
type SetBehaviorActor struct {
	behavior actor.Behavior
}

func (state *SetBehaviorActor) Receive(context actor.Context) {
	state.behavior.Receive(context)
}

func (state *SetBehaviorActor) One(context actor.Context) {
	switch msg := context.Message().(type) {
	case Hello:
		fmt.Printf("Hello %v\n", msg.Who)
		state.behavior.Become(state.Other) // 设置当前 actor 的行为
		fmt.Println("done")
	}
}

func (state *SetBehaviorActor) Other(context actor.Context) {
	switch msg := context.Message().(type) {
	case Hello:
		fmt.Printf("%v, ey we are now handling messages in another behavior", msg.Who)
	}
}

func NewSetBehaviorActor() actor.Actor {
	act := &SetBehaviorActor{
		behavior: actor.NewBehavior(),
	}
	act.behavior.Become(act.One) // 设置 actor 的行为
	return act
}

func main() {
	system := actor.NewActorSystem()
	rootContext := system.Root
	// prosFromProducer 传入参数是一个 actor 对象
	props := actor.PropsFromProducer(NewSetBehaviorActor)
	pid := rootContext.Spawn(props)

	rootContext.Send(pid, Hello{Who: "Roger"})
	fmt.Println("Roger sent")

	time.Sleep(3 * time.Second)

	rootContext.Send(pid, Hello{Who: "Roger"})
	console.ReadLine() // 等待终端读入数据
}
