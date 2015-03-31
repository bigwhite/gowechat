// Package mp provides functions for handling the received messages.
package mp

import (
	"crypto/sha1"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/bigwhite/gowechat/pb"
)

const (
	// Msg type
	TextMsg       = "text"
	ImageMsg      = "image"
	VoiceMsg      = "voice"
	VideoMsg      = "video"
	ShortVideoMsg = "shortvideo"
	LocationMsg   = "location"
	LinkMsg       = "location"
	EventMsg      = "event"

	// Event type
	SubscribeEvent   = "subscribe"
	UnsubscribeEvent = "unsubscribe"
	LocationEvent    = "LOCATION"
	MenuClickEvent   = "CLICK"
	MenuViewEvent    = "VIEW"
	ScanEvent        = "SCAN"
)

// RecvHTTPEncryptReqBody is a unmarshall result for below xml data:
// <xml>
//  <ToUserName><![CDATA[toUser]]</ToUserName>
//  <Encrypt><![CDATA[msg_encrypt]]</Encrypt>
// </xml>
type RecvHTTPEncryptReqBody struct {
	ToUserName string
	Encrypt    string
}

type recvHandler struct {
	appID          string
	token          string
	encodingAESKey string
}

// ValidateSignature is used to validate the signature in request to figure out
// whether the http request come from wechat mp platform.
func ValidateSignature(signature, token, timestamp, nonce string) bool {
	return signature == genSignature(token, timestamp, nonce)
}

func genSignature(token, timestamp, nonce string) string {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

// NewRecvHandler creates an instance of recvHandler
// which implements pb.RecvHandler interface.
func NewRecvHandler(appID, token, encodingAESKey string) pb.RecvHandler {
	return &recvHandler{appID: appID,
		token:          token,
		encodingAESKey: encodingAESKey}
}

// Parse used to parse the receive "post" data request.
// if Parse ok, it return one pkg struct of above; otherwise return error.
//
// Note: We suppose that r.ParseForm() has been invoked before entering this method.
// and we suppose that you have validate the URL in the post request.
func (h *recvHandler) Parse(bodyText []byte, signature, timestamp, nonce, encryptType string) (interface{}, error) {
	var err error
	var appID string
	var origData []byte

	if valid := ValidateSignature(signature, h.token, timestamp, nonce); !valid {
		return nil, errors.New("validate signature error")
	}

	if encryptType == "aes" {
		// Decoding the body.
		pkg := &RecvHTTPEncryptReqBody{}
		err = xml.Unmarshal(bodyText, pkg)
		if err != nil {
			return nil, err
		}
		// Decrypt the Encrypt field.
		origData, _, appID, err = DecryptMsg(pkg.Encrypt, h.encodingAESKey)
		if err != nil {
			return nil, err
		}

		if appID != h.appID {
			return nil, fmt.Errorf("the request is from app[%s], not from app[%s]", appID, h.appID)
		}
	} else {
		origData = bodyText
	}

	// Probe the type of message.
	probePkg := &struct {
		MsgType string
		Event   string
	}{}
	if err = xml.Unmarshal(origData, probePkg); err != nil {
		return nil, err
	}

	var dataPkg interface{}
	switch probePkg.MsgType {
	case TextMsg:
		dataPkg = &RecvTextDataPkg{}
	case ImageMsg:
		dataPkg = &RecvImageDataPkg{}
	case VoiceMsg:
		dataPkg = &RecvVoiceDataPkg{}
	case VideoMsg:
		dataPkg = &RecvVideoDataPkg{}
	case ShortVideoMsg:
		dataPkg = &RecvShortVideoDataPkg{}
	case LinkMsg:
		dataPkg = &RecvLinkDataPkg{}
	case LocationMsg:
		dataPkg = &RecvLocationDataPkg{}
	case EventMsg:
		switch probePkg.Event {
		case SubscribeEvent, UnsubscribeEvent:
			dataPkg = &RecvSubscribeEventDataPkg{}
		case LocationEvent:
			dataPkg = &RecvLocationEventDataPkg{}
		case MenuClickEvent, MenuViewEvent:
			dataPkg = &RecvMenuEventDataPkg{}
		case EnterAgentEvent:
			dataPkg = &RecvEnterAgentDataPkg{}
		case ScanCodePushEvent:
		case ScanCodeWaitEvent:
		case PicSysPhotoEvent:
		case PicPhotoOrAlbumEvent:
		case PicWeiXinEvent:
		case LocationSelectEvent:
		default:
			return nil, fmt.Errorf("unknown event type: %s", probePkg.Event)
		}
	default:
	}
}
