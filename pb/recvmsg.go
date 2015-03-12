// Package pb provides functions for handling the received messages.
package pb

// CDATAText is a struct whose field won't be seemed as escape sequence
// when doing xml parsing.
type CDATAText struct {
	Text string `xml:",innerxml"`
}

// RecvBaseDataPkg is the base msg struct for qy and mp receive message.
// it contains the fields shared by qy and mp receive message.
type RecvBaseDataPkg struct {
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   int
	MsgType      CDATAText
}

// RecvRespBaseDataPkg is the base msg struct for qy and mp receive response message.
// it contains the fields shared by qy and mp receive response message.
type RecvRespBaseDataPkg struct {
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   int
	MsgType      CDATAText
}

// RecvHandler is a interface for qy and mp package to implement.
type RecvHandler interface {
	Parse([]byte) (interface{}, error)
}
