package elasticemail

import (
	"context"
)

// Email is the JSON structure accepted by ElasticEmail Email API.
// fields list from https://api.elasticemail.com/public/help#Email_Send
type Email struct {
	BodyHTML          string       `json:",omitempty"`
	BodyText          string       `json:",omitempty"`
	Channel           string       `json:",omitempty"`
	Charset           string       `json:",omitempty"`
	CharsetBodyHTML   string       `json:",omitempty"`
	CharsetBodyText   string       `json:",omitempty"`
	EncodingType      encodingType `json:",omitempty"`
	From              string       `json:",omitempty"`
	FromName          string       `json:",omitempty"`
	Headers           string       `json:",omitempty"`
	IsTransactional   bool         `json:",omitempty"`
	Lists             string       `json:",omitempty"`
	Merge             string       `json:",omitempty"`
	MsgBCC            string       `json:",omitempty"`
	MsgCC             string       `json:",omitempty"`
	MsgFrom           string       `json:",omitempty"`
	MsgFromName       string       `json:",omitempty"`
	MsgTo             string       `json:",omitempty"`
	PoolName          string       `json:",omitempty"`
	PostBack          string       `json:",omitempty"`
	ReplyTo           string       `json:",omitempty"`
	ReplyToName       string       `json:",omitempty"`
	Segments          string       `json:",omitempty"`
	Sender            string       `json:",omitempty"`
	SenderName        string       `json:",omitempty"`
	Subject           string       `json:",omitempty"`
	Template          string       `json:",omitempty"`
	TimeOffSetMinutes string       `json:",omitempty"`
	To                string       `json:",omitempty"`
	File              []byte       `json:"-"`
}

type encodingType byte

// email encoding type
// https://api.elasticemail.com/public/help#classes_EncodingType
const (
	EncodingTypeBase64          encodingType = 4
	EncodingTypeNone                         = 0
	EncodingTypeQuotedPrintable              = 3
	EncodingTypeRaw7bit                      = 1
	EncodingTypeRaw8bit                      = 2
	EncodingTypeUserProvided                 = -1
	EncodingTypeUue                          = 5
)

// Send attempts to send provided email object
// https://api.elasticemail.com/public/help#Email_Send
func (c *Client) Send(e *Email) (res *Response) {
	return c.SendContext(context.Background(), e)
}

// SendContext is the same as Send, and it allows the caller to pass in a context
func (c *Client) SendContext(ctx context.Context, e *Email) *Response {
	return c.HTTPPost(ctx, "/email/send", e, e.File)
}
