package blockchain

type WsReqMessage struct {
	Id     int64
	Code   int64
	Params map[string]string
}

type WsRespMessage struct {
	Id     int64
	Code   int64
	Status int //0:成功 1：失败
	Err    string
	Params map[string]string
	Resp   interface{}
}
