// Package pb provides functions for handling the received messages.
package pb

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// CDATAText is a struct whose field won't be seemed as escape sequence
// when doing xml parsing.
type CDATAText struct {
	Text string `xml:",innerxml"`
}

func String2CDATA(v string) CDATAText {
	return CDATAText{"<![CDATA[" + v + "]]>"}
}

// RecvBaseDataPkg is the base msg struct for qy and mp receive message.
// it contains the fields shared by qy and mp receive message.
type RecvBaseDataPkg struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
}

// RecvRespBaseDataPkg is the base msg struct for qy and mp receive response message.
// it contains the fields shared by qy and mp receive response message.
type RecvRespBaseDataPkg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   int
	MsgType      CDATAText
}

// RecvHandler is a interface for qy and mp package to implement.
type RecvHandler interface {
	Parse(bodyText []byte, signature, timestamp, nonce, encryptType string) (interface{}, error)
	Response(msg []byte, encryptType string) ([]byte, error)
}

func GenNonce() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	return fmt.Sprintf("%d", r.Int31())
}

func GenTimestamp() int {
	return int(time.Now().Unix())
}

func GenSignature(args ...string) string {
	sl := args
	if len(sl) == 0 {
		return ""
	}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}
