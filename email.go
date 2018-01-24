package elasticemail

import (
	"context"
	"encoding/json"
)

// Email is the JSON structure accepted by ElasticEmail Email API.
// fields list from https://api.elasticemail.com/public/help#Email_Send
type Email struct {
	BodyHTML          string            `json:",omitempty"`
	BodyText          string            `json:",omitempty"`
	Channel           string            `json:",omitempty"`
	Charset           string            `json:",omitempty"`
	CharsetBodyHTML   string            `json:",omitempty"`
	CharsetBodyText   string            `json:",omitempty"`
	EncodingType      encodingType      `json:",omitempty"`
	From              string            `json:",omitempty"`
	FromName          string            `json:",omitempty"`
	IsTransactional   bool              `json:",omitempty"`
	Lists             string            `json:",omitempty"`
	Merge             string            `json:",omitempty"`
	MsgBCC            string            `json:",omitempty"`
	MsgCC             string            `json:",omitempty"`
	MsgFrom           string            `json:",omitempty"`
	MsgFromName       string            `json:",omitempty"`
	MsgTo             string            `json:",omitempty"`
	PoolName          string            `json:",omitempty"`
	ReplyTo           string            `json:",omitempty"`
	ReplyToName       string            `json:",omitempty"`
	Segments          string            `json:",omitempty"`
	Sender            string            `json:",omitempty"`
	SenderName        string            `json:",omitempty"`
	Subject           string            `json:",omitempty"`
	Template          string            `json:",omitempty"`
	TimeOffSetMinutes string            `json:",omitempty"`
	To                string            `json:",omitempty"`
	CustomHeaders     map[string]string `json:"-"` // postback headers are presented in email headers also
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

// Send submits email. The default, maximum (accepted by us) size of an email is 10 MB in total, with or without attachments included.
// https://api.elasticemail.com/public/help#Email_Send
func (c *Client) Send(e *Email) *Response {
	return c.SendContext(context.Background(), e)
}

// SendContext is the same as Send, and it allows the caller to pass in a context
func (c *Client) SendContext(ctx context.Context, e *Email) *Response {
	jsonBytes, _ := json.Marshal(e)
	var rawParams map[string]string
	json.Unmarshal(jsonBytes, &rawParams)

	for k, v := range e.CustomHeaders {
		rawParams[("headers_postback-" + k)] = "postback-" + k + ": " + v
	}
	return c.HTTPPost(ctx, "email/send", rawParams)
}

// Status retrieves detailed status of a unique email sent through your account. Returns a 'Email has expired and the status is unknown.' error, if the email has not been fully processed yet.
// https://api.elasticemail.com/public/help#Email_Send
func (c *Client) Status(messageID string) *Response {
	return c.StatusContext(context.Background(), messageID)
}

// StatusContext is the same as Status, and it allows the caller to pass in a context
func (c *Client) StatusContext(ctx context.Context, messageID string) *Response {
	return c.HTTPGet(ctx, "email/status", map[string]string{"messageID": messageID})
}

// View retrieves email content (public, no API key required)
// https://api.elasticemail.com/public/help#classes_EmailView
func (c *Client) View(messageID string) *Response {
	return c.ViewContext(context.Background(), messageID)
}

// ViewContext is the same as View, and it allows the caller to pass in a context
func (c *Client) ViewContext(ctx context.Context, messageID string) *Response {
	return c.HTTPGet(ctx, "email/view", map[string]string{"messageID": messageID})
}
