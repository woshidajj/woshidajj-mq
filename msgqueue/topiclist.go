package msgqueue

import (
	"fmt"
	"sync"
	"time"
)

const (
	msgHandleDuration = time.Second * 5
)

type topicList struct {
	subers []*subscriber
	sync.RWMutex
}

func newTopicList() (*topicList, error) {

	subers := make([]*subscriber, 0)
	tl := &topicList{subers: subers}

	return tl, nil
}

func (tlist *topicList) subscribe(suber *subscriber) {

	tlist.Lock()
	tlist.subers = append(tlist.subers, suber)
	tlist.Unlock()

}

func (tlist *topicList) unsubscribe(suber *subscriber) {

	tlist.Lock()
	defer tlist.Unlock()

	tempList := make([]*subscriber, 0)
	found := false

	for _, suberInList := range tlist.subers {

		if suberInList == suber {
			found = true
			continue
		}
		tempList = append(tempList, suberInList)
	}

	if found {
		tlist.subers = tempList
	}

}

func (tlist *topicList) publish(msg interface{}) error {

	for _, suber := range tlist.subers {

		go func(s *subscriber, m interface{}) {
			idleDuration := msgHandleDuration
			idleDelay := time.NewTimer(idleDuration)
			defer idleDelay.Stop()

			select {
			case s.msgC <- m:
				fmt.Printf("publish to : %p \n", s)
			case <-idleDelay.C:
				return
			}
		}(suber, msg)

	}

	return nil

}
