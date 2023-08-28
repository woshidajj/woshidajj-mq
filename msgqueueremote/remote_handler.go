package msgqueueremote

import (
	"fmt"
	"net"
	"sync"
)

type RemoteHandler struct {
	conn net.Conn
	sync.RWMutex
}

func (h *RemoteHandler) HandleMsg(m interface{}) error {

	h.Lock()
	defer h.Unlock()

	v, ok := m.(Payload)

	if ok {
		b, err := v.ToBytes()
		if err != nil {
			fmt.Println("----------start-----------")
			fmt.Println("CHANGE BYTE FAIL \n")
			fmt.Println("----------end-----------")
			return err
		}
		h.conn.Write(b)
	} else {
		fmt.Println("----------start-----------")
		fmt.Println("INVALID M \n")
		fmt.Println("----------end-----------")
	}

	return nil

}
