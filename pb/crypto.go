// Package pb provides encrpyt and descrypt for wechat message request and response
package pb

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// DecryptMsg is used to descrpyt the encrypted msg from wechat.
// it returns the origData of the cipherText. cipherText is msg_encrypt.
// origData = AES_Decrypt(Base64_Decode[cipherText])
func DecryptMsg(cipherText, encodingAESKey string) ([]byte, error) {
	var cipherData []byte
	var err error
	if cipherData, err = base64.StdEncoding.DecodeString(cipherText); err != nil {
		return nil, err
	}

	var AESKey []byte
	if AESKey, err = encodingAESKey2AESKey(encodingAESKey); err != nil {
		return nil, err
	}

	return aesDecrypt(cipherData, AESKey)
}

// EncryptMsg is used to encrpyt the msg being sent to wechat.
// it returns the cipherText of the origData.
// cipherText = Base64_Encode(AES_Encrypt [origData])
func EncryptMsg(origData []byte, encodingAESKey string) (string, error) {
	var AESKey []byte
	var err error
	if AESKey, err = encodingAESKey2AESKey(encodingAESKey); err != nil {
		return "", err
	}

	var cipherData []byte
	if cipherData, err = aesEncrypt(origData, AESKey); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(cipherData), nil
}

// AESKey = Base64_Decode(EncodingAESKey + "=").
func encodingAESKey2AESKey(encodingAESKey string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodingAESKey + "=")
}

// unpad strips the PKCS #7 padding on a buffer. If the padding is
// invalid, nil is returned.
func unpad(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}

	padding := in[len(in)-1]
	if int(padding) > len(in) {
		return nil
	} else if padding == 0 {
		return nil
	}

	for i := len(in) - 1; i > len(in)-int(padding)-1; i-- {
		if in[i] != padding {
			return nil
		}
	}
	return in[:len(in)-int(padding)]
}

func aesDecrypt(in, k []byte) ([]byte, error) {
	l := len(k) //PKCS#7
	if len(in) == 0 || len(in)%l != 0 {
		return nil, errors.New("cipher data size is zero or not multiple of AESKey length")
	}

	c, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	iv := generateIV()
	if iv == nil {
		return nil, errors.New("generateIV error")
	}

	cbc := cipher.NewCBCDecrypter(c, iv)
	cbc.CryptBlocks(in, in)

	out := unpad(in)
	if out == nil {
		return nil, errors.New("unpad error")
	}
	return out, nil
}

func generateIV() []byte {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil
	}
	return iv
}

// pad applies the PKCS #7 padding scheme on the buffer.
func pad(in []byte, length int) []byte {
	padding := length - (len(in) % length)
	if padding == 0 {
		padding = length
	}
	for i := 0; i < padding; i++ {
		in = append(in, byte(padding))
	}
	return in
}

// aesEncrypt applies the necessary padding to the message and encrypts it
// with AES-CBC.
func aesEncrypt(in, k []byte) ([]byte, error) {
	in = pad(in, len(k))

	iv := generateIV()
	if iv == nil {
		return nil, errors.New("generateIV error")
	}

	c, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(c, iv)
	cbc.CryptBlocks(in, in)
	return in, nil //do not return iv ahead of in
}
