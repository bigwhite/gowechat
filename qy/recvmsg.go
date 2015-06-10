// Package qy provides functions for handling the received messages.
package qy

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/bigwhite/gowechat/pb"
)

const (
	// Msg type
	TextMsg     = "text"
	ImageMsg    = "image"
	VoiceMsg    = "voice"
	VideoMsg    = "video"
	LocationMsg = "location"
	EventMsg    = "event"

	// Event type
	SubscribeEvent       = "subscribe"
	UnsubscribeEvent     = "unsubscribe"
	LocationEvent        = "LOCATION"
	MenuClickEvent       = "click"
	MenuViewEvent        = "view"
	ScanCodePushEvent    = "scancode_push"
	ScanCodeWaitEvent    = "scancode_waitmsg"
	PicSysPhotoEvent     = "pic_sysphoto"
	PicPhotoOrAlbumEvent = "pic_photo_or_album"
	PicWeiXinEvent       = "pic_weixin"
	LocationSelectEvent  = "location_select"
	EnterAgentEvent      = "enter_agent"
)

// RecvTextDataPkg is a Text Message received from wechat platform.
type RecvTextDataPkg struct {
	pb.RecvBaseDataPkg
	Content string
	MsgID   uint64 `xml:"MsgId"`
	AgentID int
}

// RecvImageDataPkg is a Image Message received from wechat platform.
type RecvImageDataPkg struct {
	pb.RecvBaseDataPkg
	PicURL  string `xml:"PicUrl"`
	MediaID string `xml:"MediaId"`
	MsgID   uint64 `xml:"MsgId"`
	AgentID int
}

// RecvVoiceDataPkg is a Voice Message received from wechat platform.
type RecvVoiceDataPkg struct {
	pb.RecvBaseDataPkg
	MediaID string `xml:"MediaId"`
	Format  string
	MsgID   uint64 `xml:"MsgId"`
	AgentID int
}

// RecvVideoDataPkg is a Video Message received from wechat platform.
type RecvVideoDataPkg struct {
	pb.RecvBaseDataPkg
	MediaID      string `xml:"MediaId"`
	ThumbMediaID string `xml:"ThumbMediaId"`
	MsgID        uint64 `xml:"MsgId"`
	AgentID      int
}

// RecvLocationDataPkg is a Location Message received from wechat platform.
type RecvLocationDataPkg struct {
	pb.RecvBaseDataPkg
	LocX    float64 `xml:"Location_X"`
	LocY    float64 `xml:"Location_Y"`
	Scale   int
	Label   string
	MsgID   uint64 `xml:"MsgId"`
	AgentID int
}

// RecvSubscribeEventDataPkg is a Subscribe/Unsubscribe event
// Message received from wechat platform.
type RecvSubscribeEventDataPkg struct {
	pb.RecvBaseDataPkg
	Event   string
	AgentID int
}

// RecvLocationEventDataPkg is a Location event Message
// received from wechat platform.
type RecvLocationEventDataPkg struct {
	pb.RecvBaseDataPkg
	Event     string
	Latitude  float64
	Longitude float64
	Precision float64
	AgentID   int
}

// RecvMenuEventDataPkg is a Menu Click event Message
// received from wechat platform.
type RecvMenuEventDataPkg struct {
	pb.RecvBaseDataPkg
	Event    string
	EventKey string
	AgentID  int
}

type RecvEnterAgentDataPkg struct {
	pb.RecvBaseDataPkg
	Event    string
	EventKey string
}

type RecvRespTextDataPkg struct {
	pb.RecvRespBaseDataPkg
	Content pb.CDATAText
}

type recvHandler struct {
	corpID         string
	token          string
	encodingAESKey string
}

// RecvHTTPReqBody is a unmarshall result for below xml data:
// <xml>
// 	<ToUserName><![CDATA[toUser]]</ToUserName>
// 	<AgentID><![CDATA[toAgentID]]</AgentID>
// 	<Encrypt><![CDATA[msg_encrypt]]</Encrypt>
// </xml>
type RecvHTTPReqBody struct {
	ToUserName string
	AgentID    string
	Encrypt    string
}

// RecvHTTPResqBody is a source for marshalling below xml data:
// <xml>
// 	<Encrypt><![CDATA[msg_encrypt]]></Encrypt>
// 	<MsgSignature><![CDATA[msg_signature]]></MsgSignature>
// 	<TimeStamp>timestamp</TimeStamp>
// 	<Nonce><![CDATA[nonce]]></Nonce>
// </xml>
type RecvHTTPRespBody struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      pb.CDATAText
	MsgSignature pb.CDATAText
	TimeStamp    int
	Nonce        pb.CDATAText
}

// NewRecvHandler creates an instance of recvHandler
// which implements pb.RecvHandler interface.
func NewRecvHandler(corpID, token, encodingAESKey string) pb.RecvHandler {
	return &recvHandler{corpID: corpID,
		token:          token,
		encodingAESKey: encodingAESKey}
}

// Parse used to parse the receive "post" data request.
// if Parse ok, it return one pkg struct of above; otherwise return error.
//
// Note: We suppose that r.ParseForm() has been invoked before entering this method.
// and we suppose that you have validate the URL in the post request.
func (h *recvHandler) Parse(bodyText []byte, signature, timestamp, nonce, encryptType string /* not used */) (interface{}, error) {
	var err error

	// XML decoding.
	reqBody := &RecvHTTPReqBody{}
	if err = xml.Unmarshal(bodyText, reqBody); err != nil {
		return nil, err
	}

	// Validate signature.
	if !ValidateSignature(signature, h.token, timestamp, nonce, reqBody.Encrypt) {
		return nil, errors.New("validate signature error")
	}

	// Decrpyt the "Encrypt" field.
	var origData []byte
	var corpID string
	origData, _, corpID, err = DecryptMsg(reqBody.Encrypt, h.encodingAESKey)
	if err != nil {
		return nil, err
	}

	if corpID != h.corpID {
		return nil, fmt.Errorf("the request is from corp[%s], not from corp[%s]", corpID, h.corpID)
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
		return nil, fmt.Errorf("unknown msg type: %s", probePkg.MsgType)
	}

	if err = xml.Unmarshal(origData, dataPkg); err != nil {
		return nil, err
	}
	return dataPkg, nil
}

// Response returns the response body data for the request from wechat qy platform.
func (h *recvHandler) Response(msg []byte, encryptType string /* not used */) ([]byte, error) {
	msgEncrypt, err := EncryptMsg(msg, h.corpID, h.encodingAESKey)
	if err != nil {
		return nil, err
	}

	nonce := pb.GenNonce()
	timestamp := pb.GenTimestamp()

	signature := genSignature(h.token, fmt.Sprintf("%d", timestamp), nonce, msgEncrypt)
	resp := &RecvHTTPRespBody{
		Encrypt:      pb.String2CDATA(msgEncrypt),
		MsgSignature: pb.String2CDATA(signature),
		TimeStamp:    timestamp,
		Nonce:        pb.String2CDATA(nonce),
	}

	return xml.MarshalIndent(resp, " ", "  ")
}

// ValidateSignature is used to validate the signature in request to figure out
// whether the http request come from wechat qy platform.
func ValidateSignature(signature, token, timestamp, nonce, msgEncrypt string) bool {
	return signature == genSignature(token, timestamp, nonce, msgEncrypt)
}

// dev_msg_signature=sha1(sort(token、timestamp、nonce、msg_encrypt))
func genSignature(token, timestamp, nonce, msgEncrypt string) string {
	return pb.GenSignature(token, timestamp, nonce, msgEncrypt)
}

func ValidateURL(signature, token, timestamp, nonce, cipherEchoStr, encodingAESKey string) (bool, []byte) {
	if !ValidateSignature(signature, token, timestamp, nonce, cipherEchoStr) {
		return false, nil
	}

	echostr, _, _, err := DecryptMsg(cipherEchoStr, encodingAESKey)
	if err != nil {
		return false, nil
	}
	return true, echostr
}
