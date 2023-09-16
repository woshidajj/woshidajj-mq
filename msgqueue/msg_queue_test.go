package msgqueue

import (
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"
)

// 并发订阅，检查订阅者列表
func TestMQSubscribe(t *testing.T) {
	mq, _ := NewMsgQueue()

	topic := "aaa"

	topic2 := "bbb"

	a1 := NewActor(10, &HelloActorWorker{})
	a2 := NewActor(10, &HelloActorWorker{})
	a3 := NewActor(10, &HelloActorWorker{})
	a4 := NewActor(10, &HelloActorWorker{})

	startC := make(chan struct{})

	go func() {
		<-startC
		mq.Subscribe(topic, a1)
	}()

	go func() {
		<-startC
		mq.Subscribe(topic, a2)
	}()

	go func() {
		<-startC
		mq.Subscribe(topic2, a3)
	}()

	go func() {
		<-startC
		mq.Subscribe(topic2, a4)
	}()

	close(startC)

	time.Sleep(time.Second * 5)

	mq.PrintTopiclist(topic)
	mq.PrintTopiclist(topic2)

}

// 取消订阅，并发订阅，然后一个一个取消订阅，跟踪列表的成员和协程数量
func TestMQUnsubscribe(t *testing.T) {
	mq, _ := NewMsgQueue()

	topic := "aaa"

	a1 := NewActor(10, &HelloActorWorker{})
	a2 := NewActor(10, &HelloActorWorker{})
	a3 := NewActor(10, &HelloActorWorker{})
	a4 := NewActor(10, &HelloActorWorker{})

	startC := make(chan struct{})

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	go func() {
		<-startC
		mq.Subscribe(topic, a1)
	}()

	go func() {
		<-startC
		mq.Subscribe(topic, a2)
	}()

	go func() {
		<-startC
		mq.Subscribe(topic, a3)
	}()

	go func() {
		<-startC
		mq.Subscribe(topic, a4)
	}()

	close(startC)

	time.Sleep(time.Second * 5)

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	mq.PrintTopiclist(topic)

	mq.Unsubscribe(topic, a2)

	time.Sleep(time.Second * 1)

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	mq.PrintTopiclist(topic)

	mq.Unsubscribe(topic, a4)

	time.Sleep(time.Second * 1)

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	mq.PrintTopiclist(topic)

	mq.Unsubscribe(topic, a3)

	time.Sleep(time.Second * 1)

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	mq.PrintTopiclist(topic)

	mq.Unsubscribe(topic, a1)

	time.Sleep(time.Second * 1)

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	mq.PrintTopiclist(topic)

}

// 发布消息，并发订阅，然后发布消息
func TestMQPublish(t *testing.T) {
	mq, _ := NewMsgQueue()

	topic := "aaa"

	topic2 := "bbb"

	a1 := NewActor(10, &HelloActorWorker{})
	a2 := NewActor(10, &HelloActorWorker{})
	a3 := NewActor(10, &HelloActorWorker{})
	a4 := NewActor(10, &HelloActorWorker{})

	startC := make(chan struct{})

	go func() {
		<-startC
		mq.Subscribe(topic, a1)
	}()

	go func() {
		<-startC
		mq.Subscribe(topic, a2)
	}()

	go func() {
		<-startC
		mq.Subscribe(topic, a3)
	}()

	go func() {
		<-startC
		mq.Subscribe(topic2, a4)
	}()

	close(startC)

	time.Sleep(time.Second * 5)

	mq.PrintTopiclist(topic)
	mq.PrintTopiclist(topic2)

	time.Sleep(time.Second * 1)

	go func() {

		for i := 0; i < 20; i++ {
			msg := fmt.Sprintf("gege_%d", i)
			mq.Publish(topic, msg)
		}

	}()

	time.Sleep(time.Second * 20)

}
