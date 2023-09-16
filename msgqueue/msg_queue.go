package msgqueue

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	publishDuration = time.Second * 5
)

type MsgQueue struct {
	exitC  chan bool
	topics sync.Map
}

func NewMsgQueue() (*MsgQueue, error) {

	exitC := make(chan bool)
	var topics sync.Map
	mq := &MsgQueue{exitC: exitC, topics: topics}

	return mq, nil
}

func (mq *MsgQueue) Subscribe(topic string, suber *actor) (*actor, error) {

	select {
	// wait for exit signal
	case <-mq.exitC:
		return nil, errors.New("MQ EXIT")
	default:
	}

	fmt.Printf("(%p) subscirbe %s \n", suber, topic)

	llist, _ := mq.topics.LoadOrStore(topic, newTopiclist(topic))
	tlist := llist.(*topiclist)
	tlist.subscribe(suber)

	go suber.runActor()

	return suber, nil

}

func (mq *MsgQueue) Unsubscribe(topic string, suber *actor) error {

	if suber == nil {
		return nil
	}

	// 先关闭订阅者，退出协程
	suber.close()

	select {
	// wait for exit signal
	case <-mq.exitC:
		return errors.New("MQ EXIT")
	default:
	}

	llist, ok := mq.topics.Load(topic)
	if !ok {
		return nil
	}
	tlist := llist.(*topiclist)
	tlist.unsubscribe(suber)
	return nil
}

func (mq *MsgQueue) Publish(topic string, msg interface{}) error {

	select {
	// wait for exit signal
	case <-mq.exitC:
		return errors.New("MQ EXIT")
	default:
	}

	llist, ok := mq.topics.Load(topic)
	if !ok {
		return nil
	}

	tlist := llist.(*topiclist)

	// 超时退出协程
	idleDuration := publishDuration
	idleDelay := time.NewTimer(idleDuration)
	defer idleDelay.Stop()

	msgP := PublishMsg{msg: msg, subers: tlist.subers}

	select {
	case tlist.publisher.Mailbox <- msgP:

	case <-idleDelay.C:
		return errors.New("MQ PUBLISH TIMEOUT")
	}

	return nil
}

func (mq *MsgQueue) Close() {
	close(mq.exitC)

	// 发布者和订阅者的协程
	mq.topics.Range(func(key, value interface{}) bool {

		tlist := value.(*topiclist)

		tlist.publisher.close()
		for _, suber := range tlist.subers {
			suber.close()
		}

		return true
	})

}

func (mq *MsgQueue) PrintTopiclist(topic string) {
	llist, ok := mq.topics.Load(topic)
	if !ok {
		return
	}
	tlist := llist.(*topiclist)

	for _, v := range tlist.subers {
		fmt.Printf("%p \n", v)
	}
	fmt.Println("------------------------")
}
