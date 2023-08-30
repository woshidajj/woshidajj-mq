package main

import (
	"fmt"
	"github.com/woshidajj/woshidajj-mq/msgqueueremote"
)

func main() {
	/*
		mq, _ := msgqueue.NewMsgQueue()

		sh := &handler.SimpleHandler{}
		suber := msgqueue.NewSubscriber("t", sh)
		suber1 := msgqueue.NewSubscriber("t", sh)

		mq.Subscribe(suber)
		mq.Subscribe(suber1)

		mq.Publish("t", "123")
		mq.Publish("t", "456")
		mq.Publish("t", "789")

		time.Sleep(time.Second * 5)

		mq.Unsubscribe(suber)

		mq.Publish("t", "123")
		mq.Publish("t", "456")
		mq.Publish("t", "789")

		time.Sleep(time.Second * 10)
	*/
	s, err := msgqueueremote.NewServer()

	if err != nil {
		fmt.Println(err)
	}

	err = s.Start(":8000")

	if err != nil {
		fmt.Println(err)
	}

}
