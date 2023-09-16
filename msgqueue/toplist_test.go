package msgqueue

import (
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"
)

// 并发订阅，结束看看订阅列表成员，持续跟踪协程数量对不对（每个订阅者actor有一个协程），一个topic有一个actor负责publish消息）
func TestSubscribe(t *testing.T) {

	tl := newTopiclist("abc")

	a1 := NewActor(10, &HelloActorWorker{})
	a2 := NewActor(10, &HelloActorWorker{})
	a3 := NewActor(10, &HelloActorWorker{})
	a4 := NewActor(10, &HelloActorWorker{})

	startC := make(chan struct{})
	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	go func() {
		<-startC
		fmt.Println("sub a1")
		tl.subscribe(a1)
	}()

	go func() {
		<-startC
		fmt.Println("sub a2")
		tl.subscribe(a2)
	}()

	go func() {
		<-startC
		fmt.Println("sub a3")
		tl.subscribe(a3)
	}()

	go func() {
		<-startC
		fmt.Println("sub a4")
		tl.subscribe(a4)
	}()

	close(startC)

	time.Sleep(time.Second * 5)

	for _, v := range tl.subers {
		fmt.Printf("%p \n", v)
	}

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

}

// 并发订阅，查看订阅者列表，然后取消订阅到列表订阅者为0，持续跟踪列表成员，协程数量对不对（每个订阅者actor有一个协程），一个topic有一个actor负责publish消息）
func TestUnsubscribe(t *testing.T) {

	tl := newTopiclist("abc")

	a1 := NewActor(10, &HelloActorWorker{})
	a2 := NewActor(10, &HelloActorWorker{})
	a3 := NewActor(10, &HelloActorWorker{})
	a4 := NewActor(10, &HelloActorWorker{})

	startC := make(chan struct{})

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	go func() {
		<-startC
		fmt.Println("sub a1")
		tl.subscribe(a1)
	}()

	go func() {
		<-startC
		fmt.Println("sub a2")
		tl.subscribe(a2)
	}()

	go func() {
		<-startC
		fmt.Println("sub a3")
		tl.subscribe(a3)
	}()

	go func() {
		<-startC
		fmt.Println("sub a4")
		tl.subscribe(a4)
	}()

	close(startC)

	time.Sleep(time.Second * 5)

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	for _, v := range tl.subers {
		fmt.Printf("%p \n", v)
	}

	tl.unsubscribe(a2)

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	for _, v := range tl.subers {
		fmt.Printf("%p \n", v)
	}

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	tl.unsubscribe(a4)

	for _, v := range tl.subers {
		fmt.Printf("%p \n", v)
	}

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	tl.unsubscribe(a3)

	for _, v := range tl.subers {
		fmt.Printf("%p \n", v)
	}

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	tl.unsubscribe(a1)

	for _, v := range tl.subers {
		fmt.Printf("%p \n", v)
	}
	time.Sleep(time.Second * 2)

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

}

// 发布消息，并发订阅，然后给订阅者发布消息
func TestToplistPublish(t *testing.T) {
	tl := newTopiclist("abc")

	a1 := NewActor(10, &HelloActorWorker{})
	a2 := NewActor(10, &HelloActorWorker{})
	a3 := NewActor(10, &HelloActorWorker{})
	a4 := NewActor(10, &HelloActorWorker{})
	go a1.runActor()
	go a2.runActor()
	go a3.runActor()
	go a4.runActor()

	startC := make(chan struct{})

	go func() {
		<-startC
		fmt.Println("sub a1")
		tl.subscribe(a1)
	}()

	go func() {
		<-startC
		fmt.Println("sub a2")
		tl.subscribe(a2)
	}()

	go func() {
		<-startC
		fmt.Println("sub a3")
		tl.subscribe(a3)
	}()

	go func() {
		<-startC
		fmt.Println("sub a4")
		tl.subscribe(a4)
	}()

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	close(startC)

	time.Sleep(time.Second * 2)

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	go func() {

		for i := 0; i < 20; i++ {
			msg := PublishMsg{msg: fmt.Sprintf("haha_%d", i), subers: tl.subers}
			tl.publisher.Mailbox <- msg
		}

	}()

	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	time.Sleep(time.Second * 30)
}
