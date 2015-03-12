// Package pb provides encrpyt and descrypt for wechat message request and response
package pb

import (
	"bytes"
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

func aesDecrypt(cipherData, AESKey []byte) ([]byte, error) {
	k := len(AESKey) //PKCS#7
	if len(cipherData)%k != 0 {
		return nil, errors.New("cipherData size is not multiple of AESKey length")
	}

	block, err := aes.NewCipher(AESKey)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(cipherData))
	blockMode.CryptBlocks(origData, cipherData)

	origDataLen := len(origData)
	tailPadElemNum := int(origData[origDataLen-1])
	return origData[:origDataLen-tailPadElemNum], nil
}

func aesEncrypt(origData, AESKey []byte) ([]byte, error) {
	k := len(AESKey)
	if len(origData)%k != 0 {
		origData = pkcs7pad(origData, k)
	}

	block, err := aes.NewCipher(AESKey)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cipherData := make([]byte, len(origData))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherData, origData)

	return cipherData, nil
}

// padLength calculates padding length.
// from github.com/vgorin/cryptogo.
func padLength(sliceLength, blocksize int) (padlen int) {
	padlen = blocksize - sliceLength%blocksize
	if padlen == 0 {
		padlen = blocksize
	}
	return padlen
}

// from github.com/vgorin/cryptogo.
func pkcs7pad(message []byte, blocksize int) (padded []byte) {
	// block size must be bigger or equal 2
	if blocksize < 1<<1 {
		panic("block size is too small (minimum is 2 bytes)")
	}

	// block size up to 255 requires 1 byte padding
	if blocksize < 1<<8 {
		// calculate padding length
		padlen := padLength(len(message), blocksize)

		// define PKCS7 padding block
		padding := bytes.Repeat([]byte{byte(padlen)}, padlen)

		// apply padding
		padded = append(message, padding...)
		return padded
	}

	// block size bigger or equal 256 is not currently supported
	panic("unsupported block size")
}
