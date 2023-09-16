package msgqueueremote

import (
	"fmt"
	"net"
	"testing"
)

func TestClient(t *testing.T) {

	conn, err := net.Dial("tcp", "127.0.0.1:8000")

	if err != nil {
		fmt.Printf("dial err %s", err)
		return
	}

	c, err := NewClient("t", conn)

	if err != nil {
		fmt.Printf("new client err %s", err)
		return
	}

	c.Run()

}
