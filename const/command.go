package _const

const (
	CommandConn   = iota + 0x01 // 连接请求包
	CommandSubmit               // 消息请求包
)

const (
	CommandConnAck   = iota + 0x08 // 连接请求的响应包
	CommandSubmitAck               //消息请求的响应包
)
