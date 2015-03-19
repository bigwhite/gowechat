package pb

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type TextContent struct {
	Content string `json:"content"`
}

type SendMsgTextPkg struct {
	ToUser  string      `json:"touser,omitempty"`
	MsgType string      `json:"msgtype"`
	Text    TextContent `json:"text"`
}

type MediaID struct {
	MediaID string `json:"media_id"`
}

type SendMsgImagePkg struct {
	ToUser  string  `json:"touser, omitempty"`
	MsgType string  `json:"msgtype"`
	Image   MediaID `json:"image"`
}

type SendMsgRespPkg struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func SendMsg(requestLine string, pkg interface{}) error {
	reqBody, err := json.MarshalIndent(pkg, " ", "  ")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", requestLine, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	rPkg := &SendMsgRespPkg{}
	err = json.Unmarshal(respBody, rPkg)
	if err != nil {
		return err
	}

	if rPkg.Errcode != 0 {
		return errors.New(rPkg.Errmsg)
	}

	return nil
}
