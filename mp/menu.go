// Package mp provides menu create opertations.
package mp

import (
	"strings"

	"github.com/bigwhite/gowechat/pb"
)

const (
	url = "https://api.weixin.qq.com/cgi-bin/menu/create"
)

func CreateMenu(menuLayout []byte, accessToken, agentID string) error {
	reqLine := strings.Join([]string{url, "?access_token=", accessToken}, "")
	return pb.CreateMenu(req, menuLayout)
}
