package pb

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type MenuCreateOpResp struct {
	Errcode int
	Errmsg  string
}

func CreateMenu(requestLine string, menuLayout []byte) error {
	req, err := http.NewRequest("POST",
		requestLine,
		bytes.NewReader(menuLayout))
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	opResp := &MenuCreateOpResp{}
	err = json.Unmarshal(body, opResp)
	if err != nil {
		return err
	}

	if opResp.Errcode != 0 {
		return errors.New(opResp.Errmsg)
	}

	return nil
}
