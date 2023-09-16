package msgqueue

import (
	"fmt"
	"sync"
)

type ActorWorker interface {
	Work(interface{}) error
}

type actor struct {
	id      int
	status  int
	Mailbox chan interface{}
	ExitC   chan struct{}
	once    sync.Once
	worker  ActorWorker
}

func NewActor(mailboxLen int, worker ActorWorker) *actor {

	mailbox := make(chan interface{}, mailboxLen)
	exitC := make(chan struct{})
	a := &actor{Mailbox: mailbox, ExitC: exitC, worker: worker}

	return a
}

func (a *actor) runActor() {
	for {
		select {
		case msg := <-a.Mailbox:
			a.worker.Work(msg)
		case <-a.ExitC:
			fmt.Printf("actor(%p) is closed \n", a)
			return
		}
	}
}

func (a *actor) close() {
	// 把状态设置为停止，不需要并发安全，重复关闭不会有问题
	a.status = 1
	// 只关闭一次
	a.once.Do(func() {
		close(a.ExitC)
	})
}
