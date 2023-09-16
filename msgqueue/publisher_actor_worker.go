package msgqueue

import (
	"fmt"
	"time"
)

type PublishMsg struct {
	msg    interface{} // 消息
	subers []*actor    // 订阅者的列表
}

type PublisherActorWorker struct {
}

func (w *PublisherActorWorker) Work(m interface{}) error {
	v, ok := m.(PublishMsg)

	if ok {
		msgHandleDuration := time.Second * 5
		// 对每个订阅者，开启一个协程往他mailbox发送消息
		for _, suber := range v.subers {

			go func(s *actor, m interface{}) {
				// 超时退出协程
				idleDuration := msgHandleDuration
				idleDelay := time.NewTimer(idleDuration)
				defer idleDelay.Stop()

				select {
				case s.Mailbox <- m:
					fmt.Printf("publish to : %p \n", s)
				case <-idleDelay.C:
					return
				}
			}(suber, v.msg)

		}

	} else {
		fmt.Println("INVALID M")
	}

	return nil
}
