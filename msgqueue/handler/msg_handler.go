package handler

type MsgHandler interface {
	HandleMsg(interface{}) error
}
