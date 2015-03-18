package pb_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/bigwhite/gowechat/pb"
)

const (
	postURL = "http://127.0.0.1:9001/cgi-bin/testpost"
)

func postHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("Http request parse form err:", err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("read all error:", err)
		return
	}

	pkg := &pb.SendMsgTextPkg{}
	err = json.Unmarshal(body, pkg)
	if err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}
	fmt.Println(pkg)

}

func setup() {
	// Create a server for msg send test.
	http.HandleFunc("/cgi-bin/testpost", postHandler)
	http.ListenAndServe(":9001", nil)
}

func TestSendMsg(t *testing.T) {
	pkg := &pb.SendMsgTextPkg{
		ToUser:  "tonybai",
		MsgType: "text",
		Text:    pb.TextContent{"hello body"},
	}

	err := pb.SendMsg(postURL, pkg)
	if err != nil {
		t.Fatal("SendMsg error:", err)
	}
}

func TestMain(m *testing.M) {
	go setup()
	os.Exit(m.Run())
}
