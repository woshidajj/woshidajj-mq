package main

import (
	"fmt"
	"github.com/woshidajj/woshidajj-mq/msgqueue"
	"github.com/woshidajj/woshidajj-mq/msgqueueremote"
)

func main() {

	mq, err := msgqueue.NewMsgQueue()
	address := ":8000"

	if err != nil {
		fmt.Sprintf("new mq fail, %s \n", err)
	}

	s, err := msgqueueremote.NewServer(address, mq)

	if err != nil {
		fmt.Println(err)
	}

	err = s.Start()

	if err != nil {
		fmt.Println(err)
	}

}
