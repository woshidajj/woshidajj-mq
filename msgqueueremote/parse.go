package msgqueueremote

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"strconv"
)

const (
	cmdPong = "PONG"
	cmdMsg  = "MSG"
	cmdErr  = "ERR"
)

type Payload struct {
	Command string
	Bodylen int64
	Body    []byte
	Err     error
}

func (p *Payload) ToBytes() ([]byte, error) {

	var comandLine []byte
	var bodylenLine []byte
	var bodyLine []byte
	hasBody := false

	if p.Err != nil {
		comandLine = []byte(cmdErr + "\r\n")
		bodyLine = []byte(p.Err.Error() + "\r\n")
		bodylenLine = []byte(fmt.Sprintf("%d\r\n", len(bodyLine)-2))
		hasBody = true

	} else {
		comandLine = []byte(p.Command + "\r\n")
		if p.Bodylen > 0 {
			hasBody = true
			bodylenLine = []byte(fmt.Sprintf("%d\r\n", p.Bodylen))

			if int64(len(p.Body)) != p.Bodylen {
				return nil, errors.New("body's len != bodylen")
			}

			bodyLine = make([]byte, p.Bodylen+2)
			copy(bodyLine, p.Body)
			copy(bodyLine[p.Bodylen:], []byte{'\r', '\n'})
		}

	}

	if hasBody {
		res := make([]byte, len(comandLine)+len(bodylenLine)+len(bodyLine))

		copy(res, comandLine)
		copy(res[len(comandLine):], bodylenLine)
		copy(res[len(comandLine)+len(bodylenLine):], bodyLine)

		return res, nil

	} else {
		return comandLine, nil
	}

}

func ParseStream(rawReader io.Reader, ch chan<- *Payload, stopC chan struct{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, string(debug.Stack()))
		}
	}()

	reader := bufio.NewReader(rawReader)

	for {

		select {
		// wait for exit signal
		case <-stopC:
			return
		default:
		}

		payload, hasbody, err := parseCommand(reader)

		if err != nil {
			fmt.Printf("parse command err close %s \n", err)
			close(ch)
			return
		}

		if hasbody {
			err := parseBody(reader, payload)

			if err != nil {
				fmt.Printf("parse body err close %s \n", err)
				continue
			}

		}

		ch <- payload

	}

}

func parseCommand(reader *bufio.Reader) (*Payload, bool, error) {

	hasBody := false
	rawline, err := reader.ReadBytes('\n')

	if err != nil {
		return &Payload{Err: err}, false, err
	}

	length := len(rawline)
	if length <= 2 || rawline[length-2] != '\r' || rawline[length-1] != '\n' {
		return &Payload{Err: errors.New("invalid command end")}, false, nil
	}
	commandBytes := bytes.TrimSuffix(rawline, []byte{'\r', '\n'})

	command := string(commandBytes)
	switch command {
	case cmdPong:
	case cmdMsg:
		hasBody = true
	default:
		return &Payload{Err: errors.New("invalid command")}, false, nil
	}

	return &Payload{Command: command}, hasBody, nil

}

func parseBody(reader *bufio.Reader, payload *Payload) error {

	rawline, err := reader.ReadBytes('\n')

	if err != nil {
		payload.Err = err
		return err
	}

	length := len(rawline)
	if length <= 2 || rawline[length-2] != '\r' || rawline[length-1] != '\n' {
		payload.Err = errors.New("invalid line")
		return nil
	}

	fixline := bytes.TrimSuffix(rawline, []byte{'\r', '\n'})
	nStrs, err := strconv.ParseInt(string(fixline), 10, 64)

	if err != nil {
		payload.Err = err
		return nil
	} else if nStrs <= 0 {
		payload.Err = errors.New("body len <= 0")
		return nil
	}

	// 结尾\r,\n
	payload.Body = make([]byte, nStrs+2)

	_, err = io.ReadFull(reader, payload.Body)

	if err != nil {
		payload.Err = err
		return err
	}

	if payload.Body[nStrs] != '\r' || payload.Body[nStrs+1] != '\n' {
		payload.Err = errors.New("invalid body ending")
		return nil
	}

	payload.Body = payload.Body[:nStrs]
	payload.Bodylen = nStrs

	return nil

}
