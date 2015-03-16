package pb

import (
	"bytes"
	"net/http"
)

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
	resp.Body.Close()
	return nil
}
