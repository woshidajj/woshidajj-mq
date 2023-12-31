package msgqueueremote

import (
	"bufio"
	"fmt"
	"github.com/woshidajj/woshidajj-mq/msgqueue"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	pongDuration   = time.Second * 5
	pongLimit      = 3
	routeSub       = "/sub"
	routePub       = "/pub"
	queryTopic     = "t"
	queryMsg       = "m"
	connectSuccMSg = "200 Connected to Woshidajj Msg Queue"
)

type Server struct {
	address string
	mq      *msgqueue.MsgQueue
	exitC   chan struct{}
}

func NewServer(address string, mq *msgqueue.MsgQueue) (*Server, error) {

	s := &Server{address: address, mq: mq}

	return s, nil
}

func (server *Server) handleHttpSub(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	select {
	// wait for exit signal
	case <-server.exitC:
		io.WriteString(w, "503 MQ EXIT\n")
		return
	default:
	}

	if req.Method != "CONNECT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "405 NOT CONNECT METHOD\n")
		return
	}
	conn, _, err := w.(http.Hijacker).Hijack()

	if err != nil {
		log.Print("hijacking ", req.RemoteAddr, ": ", err.Error())
		io.WriteString(w, "503 CAN'T CONNECT\n")
		return
	}

	defer conn.Close()

	query := req.URL.Query()
	topic := query.Get("t")
	if topic == "" {
		io.WriteString(w, "503 NEED QUERY t\n")
		return
	}

	actor := msgqueue.NewActor(10, &RemoteActorWorker{conn: conn})

	_, err = server.mq.Subscribe(topic, actor)

	if err != nil {
		io.WriteString(w, fmt.Sprintf("503 SUB FAIL\n"))
		return
	}

	defer server.mq.Unsubscribe(topic, actor)

	io.WriteString(conn, "HTTP/1.0 "+connectSuccMSg+"\n\n")

	runSuber(conn, actor.ExitC)

}

func (server *Server) handleHttpPub(w http.ResponseWriter, req *http.Request) {

	query := req.URL.Query()
	topic := query.Get(queryTopic)
	msg := query.Get(queryMsg)

	if topic == "" || msg == "" {
		fmt.Fprintf(w, "NEED QUERY t AND m")
		return
	}

	msgByte := []byte(msg)
	msgPl := Payload{Command: cmdMsg, Bodylen: int64(len(msgByte)), Body: msgByte}

	server.mq.Publish(topic, msgPl)
	fmt.Fprintf(w, "publish-"+topic+"-"+msg)
}

func (server *Server) Start() error {
	http.HandleFunc(routeSub, server.handleHttpSub)
	http.HandleFunc(routePub, server.handleHttpPub)

	listener, err := net.Listen("tcp", server.address)

	if err != nil {
		log.Fatal("启动服务监听失败:", err)
		return err
	}
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("启动 HTTP 服务失败:", err)
		return err
	}

	return nil
}

func runSuber(conn net.Conn, exitC chan struct{}) {

	reader := bufio.NewReader(conn)
	payloadC := make(chan *Payload)
	stopC := make(chan struct{})
	// 开启协程，从conn里读取数据，组装成payload，放进payloadC
	go ParseStream(reader, payloadC, stopC)

	// 定时器，定时检测上次接收到PONG的时间
	idleDuration := pongDuration
	idleDelay := time.NewTimer(idleDuration)
	var pong int

	for {

		idleDelay.Reset(idleDuration)

		select {
		// 监听从conn读取的payload，进行相应处理，
		case p, ok := <-payloadC:
			// chan被关闭，退出
			if !ok {
				return
			}

			if p.Err != nil {
				fmt.Printf("RECV CLIENT ERR %s \n", p.Err)
			} else if p.Command == cmdPong {
				// 如果是心跳消息PONG，要重置pong的数值
				pong = 0
				fmt.Printf("PONG \n")
			}
		case <-idleDelay.C:

			// 超过5次没收到PONG，关闭parse协程
			if pong > pongLimit {
				fmt.Printf("CLIENT TIME UP \n")
				// 关闭parse协程
				close(stopC)
				return
			} else {
				pong++
			}
		case <-exitC:
			// mq服务主动关闭了订阅者actor的协程，要主动退出
			fmt.Println("SUBER EXIT")
			// 关闭parse协程
			close(stopC)
			return
		}
	}

}
