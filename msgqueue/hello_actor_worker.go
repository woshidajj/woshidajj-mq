package msgqueue

import "fmt"

type HelloActorWorker struct {
}

func (w *HelloActorWorker) Work(m interface{}) error {
	v, ok := m.(string)

	if ok {
		fmt.Printf("HELLO - %s \n", v)
	} else {
		fmt.Println("INVALID M")
	}

	return nil
}
