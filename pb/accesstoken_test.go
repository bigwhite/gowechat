package pb_test

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/bigwhite/gowechat/pb"
)

const (
	accessTokenFetchURL = "http://127.0.0.1:9000/cgi-bin/gettoken"
	corpID              = "wxfd4448417439fd3x"
	secret              = "p_VQovLdSPNXP0caBalViFvAG_mIpR4bECn-fD1F9JRBut471AcJYXK14SOG1Zld"
	accessToken         = "accesstoken000001"
)

func tokenFetchHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("Http request parse form err:", err)
		return
	}

	corpidInReq := strings.Join(r.Form["corpid"], "")
	secretInReq := strings.Join(r.Form["corpsecret"], "")

	// Construct below error package.
	// {
	//	   "errcode": "40013",
	//	   "errmsg": "invalide corpid"
	// }
	if corpidInReq != corpID {
		respData := `{"errcode": "40013", "errmsg": "invalid corpid"}`
		w.Write([]byte(respData))
		return
	}

	if secretInReq != secret {
		respData := `{"errcode": "40014", "errmsg": "invalid corpSecret"}`
		w.Write([]byte(respData))
		return
	}

	// Construct below successful package.
	// {
	//	   "access_token": "accesstoken000001",
	//	   "expires_in": 7200
	// }
	respData := `{
			"access_token": "accesstoken000001",
			"expires_in": 7200
		}`
	w.Write([]byte(respData))
}

func setup() {
	// Create a server for access token fetching test.
	http.HandleFunc("/cgi-bin/gettoken", tokenFetchHandler)
	http.ListenAndServe(":9000", nil)
}

func TestFetchAccessTokenOk(t *testing.T) {
	myCorpID := "wxfd4448417439fd3x"
	mySecret := "p_VQovLdSPNXP0caBalViFvAG_mIpR4bECn-fD1F9JRBut471AcJYXK14SOG1Zld"
	myURL := accessTokenFetchURL + "?corpid=" + myCorpID + "&corpsecret=" + mySecret

	token, expiresIn, err := pb.FetchAccessToken(myURL)
	if err != nil {
		t.Fatal("Fetch accesstoken error:", err)
	}

	if token != accessToken {
		t.Errorf("Token: want[%s], but actually[%s]", accessToken, token)
	}

	if expiresIn != float64(7200) {
		t.Errorf("ExpiresIn: want[%s], but actually[%s]", accessToken, token)
	}
}

func TestCorpidInvalid(t *testing.T) {
	myCorpID := "wxfd4448417439fd3y"
	mySecret := "p_VQovLdSPNXP0caBalViFvAG_mIpR4bECn-fD1F9JRBut471AcJYXK14SOG1Zld"
	myURL := accessTokenFetchURL + "?corpid=" + myCorpID + "&corpsecret=" + mySecret

	_, _, err := pb.FetchAccessToken(myURL)
	errStr := fmt.Sprintf("%s", err)
	if errStr != "invalid corpid" {
		t.Errorf("Errmsg: want[%s], but actually[%s]", "invalid corpid", errStr)
	}
}

func TestSecretInvalid(t *testing.T) {
	myCorpID := "wxfd4448417439fd3x"
	mySecret := "p_VQovLdSPNXP0caBalViFvAG_mIpR4bECn-fD1F9JRBut471AcJYXK14SOG1Zle"
	myURL := accessTokenFetchURL + "?corpid=" + myCorpID + "&corpsecret=" + mySecret

	_, _, err := pb.FetchAccessToken(myURL)
	errStr := fmt.Sprintf("%s", err)
	if errStr != "invalid corpSecret" {
		t.Errorf("Errmsg: want[%s], but actually[%s]", "invalid corpSecret", errStr)
	}
}

func TestMain(m *testing.M) {
	go setup()
	os.Exit(m.Run())
}
