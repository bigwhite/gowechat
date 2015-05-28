package qy

import (
	"strings"

	"github.com/bigwhite/gowechat/pb"
)

const (
	sendURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
)

type SendMsgTextPkg struct {
	pb.SendMsgTextPkg
	ToParty string `json:"toparty,omitempty"`
	ToTag   string `json:"totag,omitempty"`
	AgentID string `json:"agentid"`
	Safe    string `json:"safe,omitempty"`
}

type SendMsgImagePkg struct {
	pb.SendMsgImagePkg
	ToParty string `json:"toparty,omitempty"`
	ToTag   string `json:"totag,omitempty"`
	AgentID string `json:"agentid"`
	Safe    string `json:"safe,omitempty"`
}

type Article struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
	PicUrl      string `json:"picurl,omitempty"`
}

type Articles struct {
	Articles []Article `json:"articles"`
}
type SendMsgNewsPkg struct {
	ToUserName string   `json:"touser"`
	ToParty    string   `json:"toparty,omitempty"`
	ToTag      string   `json:"totag,omitempty"`
	MsgType    string   `json:"msgtype"`
	AgentID    string   `json:"agentid"`
	News       Articles `json:"news"`
}

func SendMsg(accessToken string, pkg interface{}) error {
	r := strings.Join([]string{sendURL, "?access_token=", accessToken}, "")
	return pb.SendMsg(r, pkg)
}
