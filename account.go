package elasticemail

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
)

type sendingPermission byte

const (
	sendingPermissionAll                 sendingPermission = 255
	sendingPermissionHttpApi                               = 2
	sendingPermissionHttpApiAndInterface                   = 6
	sendingPermissionInterface                             = 4
	sendingPermissionNone                                  = 0
	sendingPermissionSmtp                                  = 1
	sendingPermissionSmtpAndHttpApi                        = 3
	sendingPermissionSmtpAndInterface                      = 5
)

// Subaccount is the JSON structure accepted by and returned from the SparkPost Subaccounts API.
type Subaccount struct {
	Email                  string
	Password               string
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

	// converting struct to map
	jsonBytes, _ := json.Marshal(s)
	var params map[string]string
	json.Unmarshal(jsonBytes, &params)
	params["confirmPassword"] = params["Password"]

	res = c.HttpGet(ctx, "/Account/AddSubAccount", params)
	return
}
