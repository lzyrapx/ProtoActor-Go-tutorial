package main

import (
	"log"
	"sync/atomic"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/mailbox"
)

// sent to producer to request more work
type requestMoreWork struct {
	items int
}
type requestWorkBehavior struct {
	tokens   int64
	producer *actor.PID
}

func (m *requestWorkBehavior) MailboxStarted() {
	m.requestMore()
}
func (m *requestWorkBehavior) MessagePosted(msg interface{}) {

}
func (m *requestWorkBehavior) MessageReceived(msg interface{}) {
	atomic.AddInt64(&m.tokens, -1)
	if m.tokens == 0 {
		m.requestMore()
	}
}
func (m *requestWorkBehavior) MailboxEmpty() {
}

func (m *requestWorkBehavior) requestMore() {
	log.Println("Requesting more tokens")
	m.tokens = 5
	system.Root.Send(m.producer, &requestMoreWork{items: 5})
}

type producer struct {
	requestedWork int
	producedWork  int
	worker        *actor.PID
}

func (p *producer) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started: // 启动 actor
		log.Println("actor started")
		// spawn our worker
		workerProps := actor.PropsFromProducer(func() actor.Actor {
			return &worker{}
		})
		// mailbox: 保存发往某 Actor 的消息. 通常每个 Actor 拥有自己的邮箱
		mb := mailbox.Unbounded(&requestWorkBehavior{
			producer: ctx.Self(),
		})
		p.worker = ctx.Spawn(workerProps.WithMailbox(mb)) // 启动 worker
	case *requestMoreWork: // 获取更多的 work
		log.Println("msg items : ", msg.items) // 5
		p.requestedWork += msg.items
		log.Println("Producer got a new work request：", p.requestedWork)
		ctx.Send(ctx.Self(), &produce{}) // 给 receiver actor 自己发送生产消息
	case *produce:
		// produce more work
		log.Println("Producer is producing work：", p.producedWork)
		p.producedWork++

		ctx.Send(p.worker, &work{p.producedWork}) // 给 worker actor 发送
		log.Println("decrease our word and produce more work")
		// decrease our workload and tell ourselves to produce more work
		if p.requestedWork > 0 {
			p.requestedWork--
			ctx.Send(ctx.Self(), &produce{}) // 给 receiver actor 自己发送生产消息
		}
	}
}

type produce struct{}
type worker struct{}

func (w *worker) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *work:
		log.Printf("Worker is working %+v", msg)
		time.Sleep(100 * time.Millisecond)
	}
}

type work struct {
	id int
}

var system = actor.NewActorSystem()

// 在数据流从上游生产者向下游消费者传输的过程中，上游生产速度大于下游消费速度，导致下游的 Buffer 溢出，这种现象就叫做 Backpressure 出现。
// 这句话的重点不在于「上游生产速度大于下游消费速度」，而在于「Buffer 溢出」
// Backpressure 指的是在 Buffer 有上限的系统中，Buffer 溢出的现象；它的应对措施只有一个：丢弃新事件。

func main() {
	producerProps := actor.PropsFromProducer(func() actor.Actor { return &producer{} })
	system.Root.Spawn(producerProps) // 启动 actor system
	log.Println("Done")
	_, _ = console.ReadLine()
	log.Println("finished")
}
