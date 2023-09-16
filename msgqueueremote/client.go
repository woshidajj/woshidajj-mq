package msgqueueremote

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type client struct {
	conn net.Conn
}

func NewClient(topic string, conn net.Conn) (*client, error) {

	path := routeSub + "?t=" + topic

	_, err := io.WriteString(conn, "CONNECT "+path+" HTTP/1.0\n\n")

	if err != nil {
		return nil, err
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})

	if err != nil {
		return nil, err
	}

	if resp.Status != connectSuccMSg {
		return nil, errors.New("unexpected HTTP response: " + resp.Status)
	}

	return &client{conn: conn}, nil

}

func (c *client) Run() {

	reader := bufio.NewReader(c.conn)
	payloadC := make(chan *Payload)
	stopC := make(chan struct{})
	go ParseStream(reader, payloadC, stopC)

	// 定时器，定时检测上次接收到PONG的时间
	idleDuration := pongDuration
	idleDelay := time.NewTimer(idleDuration)

	for {

		idleDelay.Reset(idleDuration)

		select {
		case p, ok := <-payloadC:
			// chan被关闭，退出
			if !ok {
				fmt.Printf("plC chan close \n")
				close(stopC)
				return
			}

			m, err := p.ToBytes()

			if err != nil {
				fmt.Printf("err pl %s \n", err)
			} else {
				fmt.Println("-----------print msg start---------")
				fmt.Printf("recv %s \n", m)
				fmt.Println("-----------print msg end---------")
			}

		case <-idleDelay.C:
			fmt.Printf("PONG \n")

			sp := Payload{Command: cmdPong}
			b, _ := sp.ToBytes()
			_, err := c.conn.Write(b)
			if err != nil {
				close(stopC)
				fmt.Printf("sth wrong %s \n", err)
			}

		}

	}
}
