package main

import (
	"log"

	console "github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/mailbox"
)

type mailboxLogger struct{}

func (m *mailboxLogger) MailboxStarted() {
	log.Print("Mailbox started")
}
func (m *mailboxLogger) MessagePosted(msg interface{}) {
	log.Printf("Message posted %v", msg)
}
func (m *mailboxLogger) MessageReceived(msg interface{}) {
	log.Printf("Message received %v", msg)
}
func (m *mailboxLogger) MailboxEmpty() {
	log.Print("No more messages")
}

func main() {
	system := actor.NewActorSystem()
	rootContext := system.Root
	props := actor.PropsFromFunc(func(ctx actor.Context) {

	}).WithMailbox(mailbox.Unbounded(&mailboxLogger{}))
	pid := rootContext.Spawn(props)
	rootContext.Send(pid, "Hello")
	_, _ = console.ReadLine()
}
