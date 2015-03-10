// Package pb provides underlying implementation for qy and mp
package pb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// AccessTokenResponse stores the normal result of access token fetching.
type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
}

// AccessTokenErrorResponse stores the error result of access token fetching.
type AccessTokenErrorResponse struct {
	Errcode float64
	Errmsg  string
}

// FetchAccessToken provides underlying access token fetching implementation.
func FetchAccessToken(requestLine string) (string, float64, error) {
	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", 0.0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0.0, err
	}

	//Json Decoding
	if bytes.Contains(body, []byte("access_token")) {
		atr := AccessTokenResponse{}
		err = json.Unmarshal(body, &atr)
		if err != nil {
			return "", 0.0, err
		}
		return atr.AccessToken, atr.ExpiresIn, nil
	}

	ater := AccessTokenErrorResponse{}
	err = json.Unmarshal(body, &ater)
	if err != nil {
		return "", 0.0, err
	}
	return "", 0.0, fmt.Errorf("%s", ater.Errmsg)
}
