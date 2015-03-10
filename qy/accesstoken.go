// Package qy provides accesstoken fetch functions for wechat qy dev
package qy

import (
	"strings"

	"github.com/bigwhite/gowechat/pb"
)

const (
	accessTokenFetchURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
)

// FetchAccessToken could be used to fetch access token for wechat qy dev.
func FetchAccessToken(corpID, corpSecret string) (string, float64, error) {
	requestLine := strings.Join([]string{accessTokenFetchURL,
		"?corpid=",
		corpID,
		"&corpsecret=",
		corpSecret}, "")

	return pb.FetchAccessToken(requestLine)
}
