// Package mp provides accesstoken fetch functions for wechat mp dev
package mp

import (
	"strings"

	"github.com/bigwhite/gowechat/pb"
)

const (
	accessTokenFetchURL = "https://api.weixin.qq.com/cgi-bin/token"
)

// FetchAccessToken could be used to fetch access token for wechat qy dev.
func FetchAccessToken(appID, appSecret string) (string, float64, error) {
	requestLine := strings.Join([]string{accessTokenFetchURL,
		"?grant_type=client_credential&appid=",
		appID,
		"&secret=",
		appSecret}, "")

	return pb.FetchAccessToken(requestLine)
}
