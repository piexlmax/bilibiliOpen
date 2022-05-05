package module

type Chan struct {
	MsgType int
	Msg     []byte
	Err     error
}
