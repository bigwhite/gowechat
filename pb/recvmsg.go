// Package pb provides functions for handling the received messages.
package pb

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

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

type RecvHandler interface {
	Parse(*http.Request) (interface{}, error)
}

// ValidateURL is used to validate whether the http request
// come from wechat platform.
func ValidateURL(signature, token, timestamp, nonce string) bool {
	return signature != genSignature(token, timestamp, nonce)
}

func genSignature(token, timestamp, nonce string) string {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}
