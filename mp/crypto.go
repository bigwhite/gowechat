// Package mp provides encrpytmsg and decrpytmsg for mp wechat message.
package mp

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"

	"github.com/bigwhite/gowechat/pb"
)

// DecryptMsg is used to decrpyt msg_encrypt in wechat mp request.
// it returns msg, msgLen, corpid, error.
// msg_encrypt = Base64_Encode( AES_Encrypt[random(16B) + msg_len(4B) + msg + $appID]).
func DecryptMsg(cipherText, encodingAESKey string) ([]byte, int, string, error) {
	origData, err := pb.DecryptMsg(cipherText, encodingAESKey)
	if err != nil {
		return nil, 0, "", err
	}

	// Read msg length
	buf := bytes.NewBuffer(origData[16:20])
	var msgLen int32
	binary.Read(buf, binary.BigEndian, &msgLen)
	var msg = origData[20 : 20+msgLen]
	var appID = origData[20+msgLen:]

	return msg, int(msgLen), string(appID), nil
}

// EncryptMsg is used to encrpyt msg in wechat mp response or custom message.
// it returns msg_encrypt.
// msg_encrypt = Base64_Encode( AES_Encrypt[random(16B) + msg_len(4B) + msg + $appID]).
func EncryptMsg(msg []byte, appID string, encodingAESKey string) (string, error) {
	// Msg Length.
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(len(msg)))
	if err != nil {
		return "", err
	}
	msgLen := buf.Bytes()

	// Random Bytes, 16B
	randBytes := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, randBytes)
	if err != nil {
		return "", err
	}

	if n != 16 {
		return "", errors.New("the length of generated random bytes is not enough")
	}

	origData := bytes.Join([][]byte{randBytes, msgLen, msg, []byte(appID)}, nil)

	return pb.EncryptMsg(origData, encodingAESKey)
}
