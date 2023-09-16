package msgqueue

import "sync"

type topiclist struct {
	topic        string
	subers       []*actor // 订阅者的列表
	sync.RWMutex          // 订阅和取消订阅的时候应该锁住
	publisher    *actor   // publisherActor
}

func newTopiclist(topic string) *topiclist {

	subers := make([]*actor, 0)
	tl := &topiclist{topic: topic, subers: subers}

	return tl
}

// 订阅，只负责添加actor
func (tlist *topiclist) subscribe(suber *actor) {

	tlist.Lock()
	defer tlist.Unlock()
	// 第一个订阅者，需要配备publisher
	if len(tlist.subers) == 0 {
		tlist.publisher = NewActor(10, &PublisherActorWorker{})
		go tlist.publisher.runActor()
	}

	tlist.subers = append(tlist.subers, suber)

}

// 取消订阅，只负责从订阅列表移除，当订阅者列表为0，关闭publisher的协程
func (tlist *topiclist) unsubscribe(suber *actor) {

	tlist.Lock()
	defer tlist.Unlock()

	// 先关闭协程
	suber.close()

	tempList := make([]*actor, 0)
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

	// 取消订阅以后列表为空，要关闭publisher的协程，下次再有订阅会重新new publisher然后启动协程
	if len(tlist.subers) == 0 {
		tlist.publisher.close()
	}

}
