// Package mp provides accesstoken fetch functions for wechat mp dev
package mp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bigwhite/gowechat/pb"
)

const (
	accessTokenFetchURL    = "https://api.weixin.qq.com/cgi-bin/token"
	webAccessTokenFetchUrl = "https://api.weixin.qq.com/sns/oauth2/access_token"
)

type WebAccessTokenResponse struct {
	AccessToken  string  `json:"access_token"`
	ExpiresIn    float64 `json:"expires_in"`
	RefreshToken string  `json:"refresh_token"`
	OpenID       string  `json:"openid"`
	Scope        string  `json:"scope"`
}

// FetchAccessToken could be used to fetch access token for wechat qy dev.
func FetchAccessToken(appID, appSecret string) (string, float64, error) {
	requestLine := strings.Join([]string{accessTokenFetchURL,
		"?grant_type=client_credential&appid=",
		appID,
		"&secret=",
		appSecret}, "")

	return pb.FetchAccessToken(requestLine)
}

// url:https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code
func FetchWebAuthInfo(appID, appSecret, code string) (*WebAccessTokenResponse, error) {
	requestLine := strings.Join([]string{webAccessTokenFetchUrl,
		"?appid=", appID, "&secret=", appSecret, "&code=", code,
		"&grant_type=authorization_code"}, "")

	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//Json Decoding
	if bytes.Contains(body, []byte("access_token")) {
		atr := WebAccessTokenResponse{}
		err = json.Unmarshal(body, &atr)
		if err != nil {
			return nil, err
		}
		return &atr, nil
	} else {
		fmt.Println("return err")
		ater := pb.AccessTokenErrorResponse{}
		err = json.Unmarshal(body, &ater)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", ater.Errmsg)
	}
}
