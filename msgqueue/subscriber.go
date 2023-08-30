package msgqueue

import (
	"fmt"
	"github.com/woshidajj/woshidajj-mq/msgqueue/handler"
)

const (
	msgCacheLen = 10
)

type subscriber struct {
	topic      string
	msgC       chan interface{}
	ExitC      chan struct{}
	msgHandler handler.MsgHandler
}

func NewSubscriber(topic string, msgHandler handler.MsgHandler) *subscriber {

	msgC := make(chan interface{}, msgCacheLen)
	exitC := make(chan struct{})
	s := &subscriber{msgC: msgC, ExitC: exitC, topic: topic, msgHandler: msgHandler}

	go s.start()

	return s
}

func (s *subscriber) start() {

	for {
		select {
		case msg := <-s.msgC:
			s.msgHandler.HandleMsg(msg)
		case <-s.ExitC:
			fmt.Printf("subscriber(%p) is closed \n", s)
			return
		}
	}

}

func (s *subscriber) close() {
	close(s.ExitC)
}
