package msgqueue

import (
	"fmt"
	"github.com/woshidajj/woshidajj-mq/msgqueue/handler"
	"sync"
)

const (
	msgCacheLen = 10
)

type subscriber struct {
	topic      string
	msgC       chan interface{}
	ExitC      chan struct{}
	msgHandler handler.MsgHandler
	once       sync.Once
}

func NewSubscriber(topic string, msgHandler handler.MsgHandler) *subscriber {

	msgC := make(chan interface{}, msgCacheLen)
	exitC := make(chan struct{})
	s := &subscriber{msgC: msgC, ExitC: exitC, topic: topic, msgHandler: msgHandler}

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
	// 只关闭一次
	s.once.Do(func() {
		close(s.ExitC)
	})
}
