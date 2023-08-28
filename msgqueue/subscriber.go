package msgqueue

import (
	"fmt"
	handler2 "github.com/woshidajj/woshidajj-mq/msgqueue/handler"
)

type subscriber struct {
	topic   string
	msgC    chan interface{}
	ExitC   chan struct{}
	handler handler2.MsgHandler
}

func NewSubscriber(topic string, handler handler2.MsgHandler) *subscriber {

	msgC := make(chan interface{}, 10)
	exitC := make(chan struct{})
	s := &subscriber{msgC: msgC, ExitC: exitC, topic: topic, handler: handler}

	go s.start()

	return s
}

func (s *subscriber) start() {

	for {
		select {
		case msg := <-s.msgC:
			s.handler.HandleMsg(msg)
		case <-s.ExitC:
			fmt.Printf("subscriber(%p) is closed \n", s)
			return
		}
	}

}

func (s *subscriber) close() {
	close(s.ExitC)
}
