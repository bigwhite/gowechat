package pb

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type TextContent struct {
	Content string `json:"content"`
}

type SendMsgTextPkg struct {
	ToUser  string      `json:"touser"`
	MsgType string      `json:"msgtype"`
	Text    TextContent `json:"text"`
}

type MediaID struct {
	MediaID string `json:"media_id"`
}

type SendMsgImagePkg struct {
	ToUser  string  `json:"touser"`
	MsgType string  `json:"msgtype"`
	Image   MediaID `json:"image"`
}

func SendMsg(requestLine string, pkg interface{}) error {
	body, err := json.MarshalIndent(pkg, " ", "  ")
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", requestLine, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}
