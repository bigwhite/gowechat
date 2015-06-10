// Package mp provides functions for handling the received messages.
package mp

import (
	"encoding/xml"
	"errors"
	"fmt"
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
	LinkMsg       = "link"
	EventMsg      = "event"

	// Event type
	SubscribeEvent   = "subscribe"
	UnsubscribeEvent = "unsubscribe"
	LocationEvent    = "LOCATION"
	MenuClickEvent   = "CLICK"
	MenuViewEvent    = "VIEW"
	ScanEvent        = "SCAN"
)

// RecvTextDataPkg is a Text Message received from wechat platform.
type RecvTextDataPkg struct {
	pb.RecvBaseDataPkg
	Content string
	MsgID   uint64 `xml:"MsgId"`
}

// RecvImageDataPkg is a Image Message received from wechat platform.
type RecvImageDataPkg struct {
	pb.RecvBaseDataPkg
	PicURL  string
	MediaID string
	MsgID   uint64 `xml:"MsgId"`
}

// RecvVoiceDataPkg is a Voice Message received from wechat platform.
type RecvVoiceDataPkg struct {
	pb.RecvBaseDataPkg
	MediaID string
	Format  string
	MsgID   uint64 `xml:"MsgId"`
}

// RecvVideoDataPkg is a Video Message received from wechat platform.
type RecvVideoDataPkg struct {
	pb.RecvBaseDataPkg
	MediaID      string
	ThumbMediaID string
	MsgID        uint64 `xml:"MsgId"`
}

// RecvShortVideoDataPkg is a short video Message received from wechat platform.
type RecvShortVideoDataPkg struct {
	pb.RecvBaseDataPkg
	MediaID      string
	ThumbMediaID string
	MsgID        uint64 `xml:"MsgId"`
}

// RecvLocationDataPkg is a Location Message received from wechat platform.
type RecvLocationDataPkg struct {
	pb.RecvBaseDataPkg
	LocX  float64 `xml:"Location_X"`
	LocY  float64 `xml:"Location_Y"`
	Scale int
	Label string
	MsgID uint64 `xml:"MsgId"`
}

// RecvLinkDataPkg is a Link Message received from wechat platform.
type RecvLinkDataPkg struct {
	pb.RecvBaseDataPkg
	Title       string
	Description string
	URL         string `xml:"Url"`
	MsgID       uint64 `xml:"MsgId"`
}

// RecvSubscribeEventDataPkg is a Subscribe/Unsubscribe event
// Message received from wechat platform.
type RecvSubscribeEventDataPkg struct {
	pb.RecvBaseDataPkg
	Event string
}

// RecvUnsubscribeScanDataPkg is a scan event from unsubscribe user
// from wechat platform.
type RecvUnsubscribeScanEventDataPkg struct {
	pb.RecvBaseDataPkg
	Event    string
	EventKey string
	Ticket   string
}

// RecvScanDataPkg is a scan event from subscribed user
// from wechat platform.
type RecvScanEventDataPkg struct {
	pb.RecvBaseDataPkg
	Event    string
	EventKey string
	Ticket   string
}

// RecvLocationEventDataPkg is a Location event Message
// received from wechat platform.
type RecvLocationEventDataPkg struct {
	pb.RecvBaseDataPkg
	Event     string
	Latitude  float64
	Longitude float64
	Precision float64
}

// RecvMenuEventDataPkg is a Menu Click event Message
// received from wechat platform.
type RecvMenuEventDataPkg struct {
	pb.RecvBaseDataPkg
	Event    string
	EventKey string
}

// RecvVoiceRecognitionDataPkg is a Voice recognition Message received from wechat platform.
type RecvVoiceRecognitionDataPkg struct {
	pb.RecvBaseDataPkg
	MediaID     string `xml:"MediaId"`
	Format      string
	Recognition string
	MsgID       uint64 `xml:"MsgId"`
}

// RecvHTTPEncryptReqBody is a unmarshall result for below xml data:
// <xml>
//  <ToUserName><![CDATA[toUser]]</ToUserName>
//  <Encrypt><![CDATA[msg_encrypt]]</Encrypt>
// </xml>
type RecvHTTPEncryptReqBody struct {
	ToUserName string
	Encrypt    string
}

// RecvHTTPEncryptResqBody is a source for marshalling below xml data:
// <xml>
//  <Encrypt><![CDATA[msg_encrypt]]></Encrypt>
//  <MsgSignature><![CDATA[msg_signature]]></MsgSignature>
//  <TimeStamp>timestamp</TimeStamp>
//  <Nonce><![CDATA[nonce]]></Nonce>
// </xml>
type RecvHTTPEncryptRespBody struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      pb.CDATAText
	MsgSignature pb.CDATAText
	TimeStamp    int
	Nonce        pb.CDATAText
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
	return pb.GenSignature(token, timestamp, nonce)
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
		MsgType     string
		Event       string
		EventKey    string
		Recognition string
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
		if probePkg.Recognition == "" {
			dataPkg = &RecvVoiceDataPkg{}
		} else {
			dataPkg = &RecvVoiceRecognitionDataPkg{}
		}
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
		case SubscribeEvent:
			if strings.Contains(probePkg.EventKey, "qrscene_") {
				dataPkg = &RecvUnsubscribeScanEventDataPkg{}
			} else {
				dataPkg = &RecvSubscribeEventDataPkg{}
			}
		case UnsubscribeEvent:
			dataPkg = &RecvSubscribeEventDataPkg{}
		case ScanEvent:
			dataPkg = &RecvScanEventDataPkg{}
		case LocationEvent:
			dataPkg = &RecvLocationEventDataPkg{}
		case MenuClickEvent, MenuViewEvent:
			dataPkg = &RecvMenuEventDataPkg{}
		default:
			return nil, fmt.Errorf("unknown event type: %s", probePkg.Event)
		}
	default:
		return nil, fmt.Errorf("unknown msg type: %s", probePkg.MsgType)
	}

	if err = xml.Unmarshal(origData, dataPkg); err != nil {
		return nil, err
	}
	return dataPkg, nil
}

// Response returns the response body data for the request from wechat mp platform.
func (h *recvHandler) Response(msg []byte, encryptType string) ([]byte, error) {
	if encryptType == "aes" {
		msgEncrypt, err := EncryptMsg(msg, h.appID, h.encodingAESKey)
		if err != nil {
			return nil, err
		}

		nonce := pb.GenNonce()
		timestamp := pb.GenTimestamp()
		signature := pb.GenSignature(h.token, fmt.Sprintf("%d", timestamp), nonce, msgEncrypt)

		resp := &RecvHTTPEncryptRespBody{
			Encrypt:      pb.String2CDATA(msgEncrypt),
			MsgSignature: pb.String2CDATA(signature),
			TimeStamp:    timestamp,
			Nonce:        pb.String2CDATA(nonce),
		}
		return xml.MarshalIndent(resp, " ", "  ")
	}
	return msg, nil
}
