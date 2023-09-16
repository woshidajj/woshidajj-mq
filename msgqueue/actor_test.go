package msgqueue

import (
	"fmt"
	"testing"
	"time"
)

func TestHelloActor(t *testing.T) {

	a := NewActor(10, &HelloActorWorker{})
	go a.runActor()
	time.Sleep(time.Second * 1)

	go func() {

		for i := 0; i < 20; i++ {
			a.Mailbox <- fmt.Sprintf("haha_%d", i)
		}

	}()

	time.Sleep(time.Second * 5)

}

func TestPublishActor(t *testing.T) {

	a1 := NewActor(10, &HelloActorWorker{})
	go a1.runActor()
	a2 := NewActor(10, &HelloActorWorker{})
	go a2.runActor()
	a3 := NewActor(10, &HelloActorWorker{})
	go a3.runActor()
	a4 := NewActor(10, &HelloActorWorker{})
	go a4.runActor()

	subers := []*actor{a1, a2, a3, a4}

	p := NewActor(10, &PublisherActorWorker{})

	go p.runActor()

	time.Sleep(time.Second * 1)

	go func() {

		for i := 0; i < 20; i++ {
			msg := PublishMsg{msg: fmt.Sprintf("haha_%d", i), subers: subers}
			p.Mailbox <- msg
		}

	}()

	time.Sleep(time.Second * 20)

}
