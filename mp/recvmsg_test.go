package mp_test

import (
	"testing"

	"github.com/bigwhite/gowechat/mp"
)

func TestValidateSignatureOk(t *testing.T) {
	signature := "78d6123977c8e5ecb255b74ecef385c5a1b5823f"
	token := "wechat4go"
	timestamp := "1426139593"
	nonce := "1326298654"

	if !mp.ValidateSignature(signature, token, timestamp, nonce) {
		t.Error("want ValidateSignture return true, but actually it returns false")
	}
}

func TestValidateSignatureFailed(t *testing.T) {
	signature := "78d6123977c8e5ecb255b74ecef385c5a1b5823e"
	token := "wechat4go"
	timestamp := "1426139593"
	nonce := "1326298654"

	if mp.ValidateSignature(signature, token, timestamp, nonce) {
		t.Error("want ValidateSignature return true, but actually it returns false")
	}
}
