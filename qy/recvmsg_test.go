package qy_test

import (
	"testing"

	"github.com/bigwhite/gowechat/qy"
)

func TestValidateSignatureOk(t *testing.T) {
	signature := "61b23841affc32e28a339764e43a9679f38ad17d"
	token := "wechat4go"
	timestamp := "1426129452"
	nonce := "1019369511"

	// For validate url request, the msgEncrpyt is the cipherEchoStr
	msgEncrpyt := "esO4Svu/v89CuQ07sXQVHN9alKpivaxlO++FrgwNaIC+oeFMQa6FC0u5OtiNb+GjRo352TIvlTjiN/xEsRaX0Q=="

	if !qy.ValidateSignature(signature, token, timestamp, nonce, msgEncrpyt) {
		t.Error("want ValidateSignature return true, but actually it returns false")
	}

	// From wechat qy official examples.
	signature = "5c45ff5e21c57e6ad56bac8758b79b1d9ac89fd3"
	token = "QDG6eK"
	timestamp = "1409659589"
	nonce = "263014780"

	//for validate url request, the msgEncrpyt is the cipherEchoStr
	msgEncrpyt = "P9nAzCzyDtyTWESHep1vC5X9xho/qYX3Zpb4yKa9SKld1DsH3Iyt3tP3zNdtp+4RPcs8TgAE7OaBO+FZXvnaqQ=="

	if !qy.ValidateSignature(signature, token, timestamp, nonce, msgEncrpyt) {
		t.Error("want ValidateSignature return true, but actually it returns false")
	}
}

func TestValidateSignatureFailed(t *testing.T) {
	signature := "61b23841affc32e28a339764e43a9679f38ad17x" // Not correct in last position.
	token := "wechat4go"
	timestamp := "1426129452"
	nonce := "1019369511"
	//for validate url request, the msgEncrpyt is the cipherEchoStr
	msgEncrpyt := "esO4Svu/v89CuQ07sXQVHN9alKpivaxlO++FrgwNaIC+oeFMQa6FC0u5OtiNb+GjRo352TIvlTjiN/xEsRaX0Q=="

	if qy.ValidateSignature(signature, token, timestamp, nonce, msgEncrpyt) {
		t.Error("want ValidateSignature return false, but actually it returns true")
	}
}

var encodingAESKey = "jRwY6v82amVaTB4eXdjG775NH8ubF6AwauNed88UfGK"

func TestValidateURLOk(t *testing.T) {
	signature := "61b23841affc32e28a339764e43a9679f38ad17d"
	token := "wechat4go"
	timestamp := "1426129452"
	nonce := "1019369511"

	echoStrWanted := "4362985891886127916"

	cipherEchoStr := "esO4Svu/v89CuQ07sXQVHN9alKpivaxlO++FrgwNaIC+oeFMQa6FC0u5OtiNb+GjRo352TIvlTjiN/xEsRaX0Q=="

	result, echoStr := qy.ValidateURL(signature, token, timestamp, nonce, cipherEchoStr, encodingAESKey)
	if !result {
		t.Error("want ValidateURL return true, but actually it returns false")
	}

	if string(echoStr) != echoStrWanted {
		t.Errorf("want [%s], but actually the echoStr is [%s]", echoStrWanted, string(echoStr))
	}
}
