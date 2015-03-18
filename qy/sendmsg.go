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

func SendMsg(accessToken string, pkg interface{}) error {
	r := strings.Join([]string{sendURL, "?access_token=", accessToken}, "")
	return pb.SendMsg(r, pkg)
}
