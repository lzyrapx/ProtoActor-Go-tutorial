package main

import (
	"fmt"
	"log"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
)

type NoInfluence string

func (NoInfluence) NotInfluenceReceiveTimeout() {}


/*
type producer struct {
	c int
}

type timeout struct {
	 c int
}
var cnt int = 0
func (p *producer) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		ctx.SetReceiveTimeout(1 * time.Second)

	case *actor.ReceiveTimeout:
		cnt++
		log.Printf("ReceiveTimeout: %d", cnt)

	case string:
		log.Printf("received '%s'", msg)
		if msg == "cancel" {
			fmt.Println("Cancelling")
			ctx.CancelReceiveTimeout()
		}

	case NoInfluence:
		log.Println("received a no-influence message")
	}
}
*/

func main() {
	log.Println("Receive timeout test")

	system := actor.NewActorSystem()
	c := 0

	//props := actor.PropsFromProducer(func() actor.Actor{ return &producer{}})
	rootContext := system.Root
	props := actor.PropsFromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *actor.Started: // actor 自带的类型
			context.SetReceiveTimeout(1 * time.Second) // 利用 SetReceiveTimeout 进行超时控制

		case *actor.ReceiveTimeout: // actor 自带的类型
			c++
			log.Printf("ReceiveTimeout: %d", c)

		case string: // 自己定义的类型
			log.Printf("received '%s'", msg)
			if msg == "cancel" {
				fmt.Println("Cancelling")
				context.CancelReceiveTimeout() // 取消超时控制
			}

		case NoInfluence: // 自己定义的类型
			log.Println("received a no-influence message: ", msg)
		}
	})

	pid := rootContext.Spawn(props)
	for i := 0; i < 6; i++ {
		rootContext.Send(pid, "hello")
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("hit [return] to send no-influence messages")
	resp, _ := console.ReadLine() // 从终端读取数据
	log.Println("resp = ", resp)

	for i := 0; i < 6; i++ {
		rootContext.Send(pid, NoInfluence("hello"))
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("hit [return] to send a message to cancel the timeout")
	_, _ = console.ReadLine()
	rootContext.Send(pid, "cancel")

	log.Println("hit [return] to finish")
	_, _ = console.ReadLine()

	rootContext.Stop(pid)
}
