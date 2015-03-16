package qy_test

import (
	"encoding/xml"
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

func TestDecryptMsg1(t *testing.T) {
	corpID := "wx2f6d0a549c129f06"
	/*
		<xml>
		<ToUserName><![CDATA[wx2f6d0a549c129f06]]></ToUserName>
		<FromUserName><![CDATA[baim]]></FromUserName>
		<CreateTime>1426498001</CreateTime>
		<MsgType><![CDATA[text]]></MsgType>
		<Content><![CDATA[hello body]]></Content>
		<MsgId>000001</MsgId>
		<AgentID>3</AgentID>
		</xml>
	*/
	msgEncrypt := "sq1d1sgR6C39QKNRJk21zIwWZrVY4EJrpX3cVJznqSqeNJjbzbjUOMnrFAHGREBizLgVU68/IOWNE5VVzQH7cYG9CVHVtS10SJepGDXhvPjXxdsyRkoxX9YJcEsxQkV4u8niGDDSfUW69d93u2V1/gMfnkxo+0yHMZcS6rvRhMYA0O8TiE2W3K092ELdfWsLxNy2Gd/+Uv9D6IcyQ8uO/1Vu6x0KhuG9EtVSooEfdqqdpkOKyiaXn4bf/Umn0PQTurrO6Fh6ghgxPxpMIcSEzhfAMMCn14pojlt113yjrh6x1vYj3gElWGeMiOm3fpjuplOVwoDSVzcPaR5zLPgizAO3WQj0ho0JQh4RJ6ZRpmaDaPlHHBX7hAiOFOyc3bScUQQfk6tOfwAOAn4x44og+INaKJqtFhsJ6Wavr+H5mYo="
	encodingAESKey := "jRwY6v82amVaTB4eXdjG775NH8ubF6AwauNed88UfGK"

	msg, _, corpIDDecrypted, err := qy.DecryptMsg(msgEncrypt, encodingAESKey)
	if err != nil {
		t.Fatal("DecryptMsg error:", err)
	}

	var recvMsg = &qy.RecvTextDataPkg{}
	err = xml.Unmarshal(msg, recvMsg)
	if err != nil {
		t.Fatal("Xml decoding error:", err)
	}

	if corpIDDecrypted != corpID {
		t.Errorf("Corpid: want[%s], but actually[%s]", corpID, corpIDDecrypted)
	}

	if recvMsg.Content != "hello body" {
		t.Errorf("Msg: want[%s], but actually[%s]", "hello body", recvMsg)
	}
}

func TestEcryptMsg(t *testing.T) {
	corpID := "wx2f6d0a549c129f06"
	msgText := `<xml>
		<ToUserName><![CDATA[wx2f6d0a549c129f06]]></ToUserName>
		<FromUserName><![CDATA[baim]]></FromUserName>
		<CreateTime>1426498001</CreateTime>
		<MsgType><![CDATA[text]]></MsgType>
		<Content><![CDATA[hello body]]></Content>
		<MsgId>000001</MsgId>
		<AgentID>3</AgentID>
		</xml>`
	encodingAESKey := "jRwY6v82amVaTB4eXdjG775NH8ubF6AwauNed88UfGK"
	msgEncrypt, err := qy.EncryptMsg([]byte(msgText), corpID, encodingAESKey)
	if err != nil {
		t.Fatal("EcryptMsg error:", err)
	}

	msg, msgLen, corpIDDecrypted, err := qy.DecryptMsg(msgEncrypt, encodingAESKey)
	if err != nil {
		t.Fatal("DecryptMsg error:", err)
	}

	var recvMsg = &qy.RecvTextDataPkg{}
	err = xml.Unmarshal(msg, recvMsg)
	if err != nil {
		t.Fatal("Xml decoding error:", err)
	}

	if corpIDDecrypted != corpID {
		t.Errorf("Corpid: want[%s], but actually[%s]", corpID, corpIDDecrypted)
	}

	if recvMsg.Content != "hello body" {
		t.Errorf("Msg: want[%s], but actually[%s]", "hello body", recvMsg)
	}

	if msgLen != len(msgText) {
		t.Errorf("MsgLen: want[%d], but actually[%d]", len(msgText), msgLen)
	}
}
