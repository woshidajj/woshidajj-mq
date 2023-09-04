package msgqueue

import (
	"errors"
	"fmt"
	"sync"
)

type MsgQueue struct {
	ExitC  chan bool
	topics sync.Map
}

func NewMsgQueue() (*MsgQueue, error) {

	exitC := make(chan bool)
	var topics sync.Map
	s := &MsgQueue{ExitC: exitC, topics: topics}

	return s, nil
}

func (mq *MsgQueue) Subscribe(suber *subscriber) (*subscriber, error) {

	select {
	// wait for exit signal
	case <-mq.ExitC:
		return nil, errors.New("MQ EXIT")
	default:
	}

	fmt.Printf("(%p) subscirbe %s \n", suber, suber.topic)

	nlist, _ := newTopicList()
	llist, _ := mq.topics.LoadOrStore(suber.topic, nlist)
	tlist := llist.(*topicList)
	tlist.subscribe(suber)

	go suber.start()

	return suber, nil
}

func (mq *MsgQueue) Publish(topic string, msg interface{}) error {

	select {
	// wait for exit signal
	case <-mq.ExitC:
		return errors.New("MQ EXIT")
	default:
	}

	llist, ok := mq.topics.Load(topic)
	if !ok {
		return nil
	}

	tlist := llist.(*topicList)
	tlist.publish(msg)

	return nil
}

func (mq *MsgQueue) Close() {
	close(mq.ExitC)

	mq.topics.Range(func(key, value interface{}) bool {

		tlist := value.(*topicList)

		for _, suber := range tlist.subers {
			mq.Unsubscribe(suber)
		}

		return true
	})

}

func (mq *MsgQueue) Unsubscribe(suber *subscriber) error {

	if suber == nil {
		return nil
	}

	// 先关闭订阅者，退出协程
	suber.close()

	select {
	// wait for exit signal
	case <-mq.ExitC:
		return errors.New("MQ EXIT")
	default:
	}

	llist, ok := mq.topics.Load(suber.topic)
	if !ok {
		return nil
	}
	tlist := llist.(*topicList)
	tlist.unsubscribe(suber)
	return nil
}
