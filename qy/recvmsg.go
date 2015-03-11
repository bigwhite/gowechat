// Package qy provides functions for handling the received messages.
package qy

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

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
	MenuClickEvent       = "CLICK"
	MenuViewEvent        = "VIEW"
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
	Content pb.CDATAText
	MsgID   uint64
	AgentID int
}

// RecvImageDataPkg is a Image Message received from wechat platform.
type RecvImageDataPkg struct {
	pb.RecvBaseDataPkg
	PicURL  pb.CDATAText
	MediaID pb.CDATAText
	MsgID   uint64
	AgentID int
}

// RecvVoiceDataPkg is a Voice Message received from wechat platform.
type RecvVoiceDataPkg struct {
	pb.RecvBaseDataPkg
	MediaID pb.CDATAText
	Format  pb.CDATAText
	MsgID   uint64
	AgentID int
}

// RecvVideoDataPkg is a Video Message received from wechat platform.
type RecvVideoDataPkg struct {
	pb.RecvBaseDataPkg
	MediaID      pb.CDATAText
	ThumbMediaID pb.CDATAText
	MsgID        uint64
	AgentID      int
}

// RecvLocationDataPkg is a Location Message received from wechat platform.
type RecvLocationDataPkg struct {
	pb.RecvBaseDataPkg
	LocX    float64 `xml:"Location_X"`
	LocY    float64 `xml:"Location_Y"`
	Scale   int
	Label   pb.CDATAText
	MsgID   uint64
	AgentID int
}

// RecvSubscribeEventDataPkg is a Subscribe/Unsubscribe event
// Message received from wechat platform.
type RecvSubscribeEventDataPkg struct {
	pb.RecvBaseDataPkg
	Event   pb.CDATAText
	AgentID int
}

// RecvLocationEventDataPkg is a Location event Message
// received from wechat platform.
type RecvLocationEventDataPkg struct {
	pb.RecvBaseDataPkg
	Event     pb.CDATAText
	Latitude  float64
	Longitude float64
	Precision float64
	AgentID   int
}

// RecvMenuEventDataPkg is a Menu Click event Message
// received from wechat platform.
type RecvMenuEventDataPkg struct {
	pb.RecvBaseDataPkg
	Event    pb.CDATAText
	EventKey pb.CDATAText
	AgentID  int
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
	ToUserName pb.CDATAText
	AgentID    pb.CDATAText
	Encrypt    pb.CDATAText
}

// RecvHTTPResqBody is a source for marshalling below xml data:
// <xml>
// 	<Encrypt><![CDATA[msg_encrypt]]></Encrypt>
// 	<MsgSignature><![CDATA[msg_signature]]></MsgSignature>
// 	<TimeStamp>timestamp</TimeStamp>
// 	<Nonce><![CDATA[nonce]]></Nonce>
// </xml>
type RecvHTTPRespBody struct {
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
func (h *recvHandler) Parse(r *http.Request) (interface{}, error) {
	var bodyText []byte
	var err error

	// HTTP decoding.
	if bodyText, err = ioutil.ReadAll(r.Body); err != nil {
		return nil, err
	}

	// XML decoding.
	reqBody := &RecvHTTPReqBody{}
	if err = xml.Unmarshal(bodyText, reqBody); err != nil {
		return nil, err
	}

	// Decrpyt the "Encrypt" field.
	var origData []byte
	if origData, err = pb.DecryptMsg(reqBody.Encrypt.Text, h.encodingAESKey); err != nil {
		return nil, err
	}

	// Probe the type of message.
	probePkg := &struct {
		MsgType pb.CDATAText
		Event   pb.CDATAText
	}{}
	if err = xml.Unmarshal(origData, probePkg); err != nil {
		return nil, err
	}

	var dataPkg interface{}
	switch probePkg.MsgType.Text {
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
		switch probePkg.Event.Text {
		case SubscribeEvent, UnsubscribeEvent:
			dataPkg = &RecvSubscribeEventDataPkg{}
		case LocationEvent:
			dataPkg = &RecvLocationEventDataPkg{}
		case MenuClickEvent, MenuViewEvent:
			dataPkg = &RecvMenuEventDataPkg{}
		case ScanCodePushEvent:
		case ScanCodeWaitEvent:
		case PicSysPhotoEvent:
		case PicPhotoOrAlbumEvent:
		case PicWeiXinEvent:
		case LocationSelectEvent:
		case EnterAgentEvent:
		default:
			return nil, fmt.Errorf("unknown event type: %s", probePkg.Event.Text)
		}

	default:
		return nil, fmt.Errorf("unknown msg type: %s", probePkg.MsgType.Text)
	}

	if err = xml.Unmarshal(origData, dataPkg); err != nil {
		return nil, err
	}
	return dataPkg, nil
}
