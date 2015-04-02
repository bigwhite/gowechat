package mp

import (
	"strings"

	"github.com/bigwhite/gowechat/pb"
)

const (
	sendURL = "https://api.weixin.qq.com/cgi-bin/message/custom/send"
)

func SendMsg(accessToken string, pkg interface{}) error {
	r := strings.Join([]string{sendURL, "?access_token=", accessToken}, "")
	return pb.SendMsg(r, pkg)
}
