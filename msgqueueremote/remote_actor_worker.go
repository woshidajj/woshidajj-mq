package msgqueueremote

import (
	"fmt"
	"net"
)

type RemoteActorWorker struct {
	conn net.Conn
}

func (w *RemoteActorWorker) Work(m interface{}) error {

	v, ok := m.(Payload)

	if ok {
		b, err := v.ToBytes()
		if err != nil {
			fmt.Println("----------start-----------")
			fmt.Println("CHANGE BYTE FAIL")
			fmt.Println("----------end-----------")
			return err
		}
		w.conn.Write(b)
	} else {
		fmt.Println("----------start-----------")
		fmt.Println("INVALID M")
		fmt.Println("----------end-----------")
	}

	return nil

}
