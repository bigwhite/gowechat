package pb_test

import (
	"encoding/xml"
	"testing"

	"github.com/bigwhite/gowechat/pb"
)

func TestString2CDATA(t *testing.T) {
	want := `<![CDATA[toUser]]>`
	field := pb.String2CDATA("toUser")

	if field.Text != want {
		t.Errorf("Want [%s], but actual[%s]", want, field.Text)
	}
}

func TestParseRecvBaseDataPkg(t *testing.T) {
	var data = &pb.RecvBaseDataPkg{}
	var pkg = `
	<xml>
		<ToUserName><![CDATA[toUser]]></ToUserName>
		<FromUserName><![CDATA[fromUser]]></FromUserName> 
		<CreateTime>1348831860</CreateTime>
		<MsgType><![CDATA[text]]></MsgType>
	</xml>`

	err := xml.Unmarshal([]byte(pkg), data)
	if err != nil {
		t.Error("Xml unmarshal error:", err)
	}

	if data.ToUserName != "toUser" {
		t.Errorf("ToUserName: want[%s], actual[%s]", "toUser", data.ToUserName)
	}
	if data.FromUserName != "fromUser" {
		t.Errorf("FromUserName: want[%s], actual[%s]", "fromUser", data.FromUserName)
	}
	if data.CreateTime != 1348831860 {
		t.Errorf("CreateTime: want[%d], actual[%d]", 1348831860, data.CreateTime)
	}
	if data.MsgType != "text" {
		t.Errorf("MsgType: want[%s], actual[%s]", "text", data.MsgType)
	}
}

func TestGenRecvRespBaseDataPkg(t *testing.T) {
	var data = &pb.RecvRespBaseDataPkg{
		ToUserName:   pb.String2CDATA("toUser"),
		FromUserName: pb.String2CDATA("fromUser"),
		CreateTime:   1348831860,
		MsgType:      pb.String2CDATA("text"),
	}

	var want = `<xml>
  <ToUserName><![CDATA[toUser]]></ToUserName>
  <FromUserName><![CDATA[fromUser]]></FromUserName>
  <CreateTime>1348831860</CreateTime>
  <MsgType><![CDATA[text]]></MsgType>
</xml>`
	pkg, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Error("Xml marshalling error:", err)
	}

	if string(pkg) != want {
		t.Errorf("Want [%s], but actual[%s]", want, string(pkg))
	}
}
