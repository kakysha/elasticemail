package elasticemail

import (
	"context"
	"github.com/pkg/errors"
)

type sendingPermission byte

// sending permission enum
const (
	SendingPermissionAll                 sendingPermission = 255
	SendingPermissionHTTPAPI                               = 2
	SendingPermissionHTTPAPIAndInterface                   = 6
	SendingPermissionInterface                             = 4
	SendingPermissionNone                                  = 0
	SendingPermissionSMTP                                  = 1
	SendingPermissionSMTPAndHTTPAPI                        = 3
	SendingPermissionSMTPAndInterface                      = 5
)

// Subaccount is the JSON structure accepted by and returned from the SparkPost Subaccounts API.
type Subaccount struct {
	Email                  string
	Password               string
	ConfirmPassword        string
	DailySendLimit         int
	EmailSizeLimit         int // MB
	EnableContactFeatures  bool
	EnableLitmusTest       bool
	EnablePrivateIPRequest bool
	MaxContacts            int
	PoolName               string
	RequiresEmailCredits   bool
	RequiresLitmusCredits  bool
	SendActivation         bool
	SendingPermission      sendingPermission
}

// AddSubAccount attempts to create a subaccount using the provided object
func (c *Client) AddSubAccount(s *Subaccount) (res *Response) {
	return c.AddSubAccountContext(context.Background(), s)
}

// AddSubAccountContext is the same as AddSubAccount, and it allows the caller to pass in a context
func (c *Client) AddSubAccountContext(ctx context.Context, s *Subaccount) (res *Response) {
	// enforce required parameters
	if s == nil {
		res.Error = errors.New("Create called with nil Subaccount")
		return
	}

	if s.ConfirmPassword == "" {
		s.ConfirmPassword = s.Password
	}

	res = c.HTTPGet(ctx, "/Account/AddSubAccount", s)
	return
}
