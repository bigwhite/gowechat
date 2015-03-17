// Package qy provides encrpytmsg and decrpytmsg for qy wechat message.
package qy

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"

	"github.com/bigwhite/gowechat/pb"
)

// DecryptMsg is used to decrpyt msg_encrypt in wechat qy request.
// it returns msg, msgLen, corpid, error.
// msg_encrypt = Base64_Encode( AES_Encrypt[random(16B) + msg_len(4B) + msg + $CorpID]).
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
	var corpID = origData[20+msgLen:]

	return msg, int(msgLen), string(corpID), nil
}

// EncryptMsg is used to encrpyt msg in wechat qy response or custom message.
// it returns msg_encrypt.
// msg_encrypt = Base64_Encode( AES_Encrypt[random(16B) + msg_len(4B) + msg + $CorpID]).
func EncryptMsg(msg []byte, corpID string, encodingAESKey string) (string, error) {
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

	origData := bytes.Join([][]byte{randBytes, msgLen, msg, []byte(corpID)}, nil)

	return pb.EncryptMsg(origData, encodingAESKey)
}
