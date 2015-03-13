package qy_test

import (
	"testing"

	"github.com/bigwhite/gowechat/qy"
)

func TestDecryptMsg(t *testing.T) {
	corpID := "wx5823bf96d3bd56c7"
	msgEncrypt := "RypEvHKD8QQKFhvQ6QleEB4J58tiPdvo+rtK1I9qca6aM/wvqnLSV5zEPeusUiX5L5X/0lWfrf0QADHHhGd3QczcdCUpj911L3vg3W/sYYvuJTs3TUUkSUXxaccAS0qhxchrRYt66wiSpGLYL42aM6A8dTT+6k4aSknmPj48kzJs8qLjvd4Xgpue06DOdnLxAUHzM6+kDZ+HMZfJYuR+LtwGc2hgf5gsijff0ekUNXZiqATP7PF5mZxZ3Izoun1s4zG4LUMnvw2r+KqCKIw+3IQH03v+BCA9nMELNqbSf6tiWSrXJB3LAVGUcallcrw8V2t9EL4EhzJWrQUax5wLVMNS0+rUPA3k22Ncx4XXZS9o0MBH27Bo6BpNelZpS+/uh9KsNlY6bHCmJU9p8g7m3fVKn28H3KDYA5Pl/T8Z1ptDAVe0lXdQ2YoyyH2uyPIGHBZZIs2pDBS8R07+qN+E7Q=="
	encodingAESKey := "jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C"

	_, _, corpIDDecrypted, err := qy.DecryptMsg(msgEncrypt, encodingAESKey)
	if err != nil {
		t.Fatal("DecryptMsg error:", err)
	}

	if corpIDDecrypted != corpID {
		t.Errorf("Corpid: want[%s], but actually[%s]", corpID, corpIDDecrypted)
	}
}
