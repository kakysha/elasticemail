package elasticemail

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
)

// Subaccount is the JSON structure accepted by ElasticEmail Subaccounts API.
// list of fields: https://api.elasticemail.com/public/help#Account_AddSubAccount
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
	APIKey                 string
}

// sending permission enum
// https://api.elasticemail.com/public/help#classes_SendingPermission
type sendingPermission byte

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

// AddSubAccount attempts to create a subaccount using the provided object
// https://api.elasticemail.com/public/help#Account_AddSubAccount
func (c *Client) AddSubAccount(s *Subaccount) (res *Response) {
	return c.AddSubAccountContext(context.Background(), s)
}

// AddSubAccountContext is the same as AddSubAccount, and it allows the caller to pass in a context
func (c *Client) AddSubAccountContext(ctx context.Context, s *Subaccount) *Response {
	if s.ConfirmPassword == "" {
		s.ConfirmPassword = s.Password
	}

	res := c.HTTPGet(ctx, "/account/addsubaccount", s)
	if res.Success {
		s.APIKey = res.Data.(string)
	}
	return res
}

// DeleteSubAccount attempts to delete subaccount using the provided params map
// one of subAccountEmail or publicAccountID must be provided
// https://api.elasticemail.com/public/help#Account_DeleteSubAccount
func (c *Client) DeleteSubAccount(params map[string]string) (res *Response) {
	return c.DeleteSubAccountContext(context.Background(), params)
}

// DeleteSubAccountContext is the same as DeleteSubAccount, and it allows the caller to pass in a context
func (c *Client) DeleteSubAccountContext(ctx context.Context, params map[string]string) *Response {
	_, emailPrs := params["subAccountEmail"]
	_, IDPrs := params["publicAccountID"]
	if !emailPrs && !IDPrs {
		return &Response{Error: errors.New("DeleteSubAccount called without email or ID")}
	}

	return c.HTTPGet(ctx, "/account/deletesubaccount", params)
}

// UpdateSubAccountSettings attempts to update a subaccount
// referenced by ID from params object and with fields from Subaccount object
// https://api.elasticemail.com/public/help#Account_UpdateSubAccountSettings
func (c *Client) UpdateSubAccountSettings(s *Subaccount, params map[string]string) (res *Response) {
	return c.UpdateSubAccountSettingsContext(context.Background(), s, params)
}

// UpdateSubAccountSettingsContext is the same as UpdateSubAccountSettings, and it allows the caller to pass in a context
func (c *Client) UpdateSubAccountSettingsContext(ctx context.Context, s *Subaccount, params map[string]string) *Response {
	_, emailPrs := params["subAccountEmail"]
	_, IDPrs := params["publicAccountID"]
	if !emailPrs && !IDPrs {
		return &Response{Error: errors.New("UpdateSubAccountSettings called without email or ID")}
	}

	jsonBytes, _ := json.Marshal(s)
	var rawParams map[string]string
	json.Unmarshal(jsonBytes, &rawParams)

	for k := range params {
		rawParams[k] = params[k]
	}

	return c.HTTPGet(ctx, "/account/updatesubaccountsettings", rawParams)
}

// GetSubAccountApiKey attempts to retrieve subaccount API Key
// one of subAccountEmail or publicAccountID must be provided
// https://api.elasticemail.com/public/help#Account_GetSubAccountApiKey
func (c *Client) GetSubAccountApiKey(params map[string]string) (res *Response) {
	return c.GetSubAccountApiKeyContext(context.Background(), params)
}

// GetSubAccountApiKeyContext is the same as GetSubAccountApiKey, and it allows the caller to pass in a context
func (c *Client) GetSubAccountApiKeyContext(ctx context.Context, params map[string]string) *Response {
	_, emailPrs := params["subAccountEmail"]
	_, IDPrs := params["publicAccountID"]
	if !emailPrs && !IDPrs {
		return &Response{Error: errors.New("GetSubAccountApiKey called without email or ID")}
	}

	return c.HTTPGet(ctx, "/account/getsubaccountapikey", params)
}
