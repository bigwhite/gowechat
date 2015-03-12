package pb_test

import (
	"testing"

	"github.com/bigwhite/gowechat/pb"
)

func TestDecryptMsg(t *testing.T) {
	corpID := "wx2f6d0a549c129f06"
	msgEncrypt := "8xkfcZ4H50bmdMXxUh1fi9sPYboKROAgPN9Obyvjxs/q9CGjRoVJJpGUz3A4XzXSI/faqOZClv0Y+aHlYMqhBpyzO6wd9iIPNcShnPTet/lUisLiz4moEeqEnLrJo25slK5j7zuI0lrLu9EnMArdYFNHd4J/rr+SK3hNh3zyXin+wEC+RuLJG+TG32AizGCcuPfe0db9/jvID8pqWjE/+Q08aaecMWhSDFk2VbWT8I5TKdo/MUsj+NQMg3c5Z4WB0fkSF8JWGz8VDnIofo9FvsFUCc3BvjLqcTldYTHE/65Qn9COdsd9qwAsPZoPdjpFRB5pl3lPjeoSW/WzT+lL5V+Y/5VfcvniZAzVKDoCdtV8Ufzs+H7JRDa/yGGMYT48AY1skYdFA00aUAOJkeTPDEtz8CtZcREYsiSnGMpgFTY="
	encodingAESKey := "jRwY6v82amVaTB4eXdjG775NH8ubF6AwauNed88UfGK"

	origData, err := pb.DecryptMsg(msgEncrypt, encodingAESKey)
	if err != nil {
		t.Errorf("err [%s] is not what we wanted", err)
	}

	origDataLen := len(origData)
	corpIDLen := len(corpID)

	tailCorpID := string(origData[origDataLen-corpIDLen:])
	if corpID != tailCorpID {
		t.Errorf("want corpID [%s], actually it is [%s]", corpID, tailCorpID)
	}

	corpID = "wx5823bf96d3bd56c7"
	msgEncrypt = "RypEvHKD8QQKFhvQ6QleEB4J58tiPdvo+rtK1I9qca6aM/wvqnLSV5zEPeusUiX5L5X/0lWfrf0QADHHhGd3QczcdCUpj911L3vg3W/sYYvuJTs3TUUkSUXxaccAS0qhxchrRYt66wiSpGLYL42aM6A8dTT+6k4aSknmPj48kzJs8qLjvd4Xgpue06DOdnLxAUHzM6+kDZ+HMZfJYuR+LtwGc2hgf5gsijff0ekUNXZiqATP7PF5mZxZ3Izoun1s4zG4LUMnvw2r+KqCKIw+3IQH03v+BCA9nMELNqbSf6tiWSrXJB3LAVGUcallcrw8V2t9EL4EhzJWrQUax5wLVMNS0+rUPA3k22Ncx4XXZS9o0MBH27Bo6BpNelZpS+/uh9KsNlY6bHCmJU9p8g7m3fVKn28H3KDYA5Pl/T8Z1ptDAVe0lXdQ2YoyyH2uyPIGHBZZIs2pDBS8R07+qN+E7Q=="
	encodingAESKey = "jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C"

	origData, err = pb.DecryptMsg(msgEncrypt, encodingAESKey)
	if err != nil {
		t.Errorf("err [%s] is not what we wanted", err)
	}
	origDataLen = len(origData)
	corpIDLen = len(corpID)

	tailCorpID = string(origData[origDataLen-corpIDLen:])
	if corpID != tailCorpID {
		t.Errorf("want corpID [%s], actually it is [%s]", corpID, tailCorpID)
	}
}

func TestEncryptMsg(t *testing.T) {

}
