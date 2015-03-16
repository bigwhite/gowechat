// Package qy provides menu create opertations.
package qy

import (
	"fmt"
	"strings"

	"github.com/bigwhite/gowechat/pb"
)

const (
	url = "https://qyapi.weixin.qq.com/cgi-bin/menu/create"
)

func CreateMenu(menuLayout []byte, accessToken, agentID string) error {
	reqLine := strings.Join([]string{url, "?access_token=", accessToken, "&agentid=", agentID}, "")
	fmt.Println(reqLine)
	return pb.CreateMenu(reqLine, menuLayout)

}
