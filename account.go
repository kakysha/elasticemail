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
	DailySendLimit         int               `json:",omitempty"`
	EmailSizeLimit         int               `json:",omitempty"` // MB
	EnableContactFeatures  bool              `json:",omitempty"`
	EnableLitmusTest       bool              `json:",omitempty"`
	EnablePrivateIPRequest bool              `json:",omitempty"`
	MaxContacts            int               `json:",omitempty"`
	PoolName               string            `json:",omitempty"`
	RequiresEmailCredits   bool              `json:",omitempty"`
	RequiresLitmusCredits  bool              `json:",omitempty"`
	SendActivation         bool              `json:",omitempty"`
	SendingPermission      sendingPermission `json:",omitempty"`
	APIKey                 string            `json:",omitempty"`
}

type sendingPermission byte

// sending permission enum
// https://api.elasticemail.com/public/help#classes_SendingPermission
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

// AddSubAccount create new subaccount and provide most important data about it.
// https://api.elasticemail.com/public/help#Account_AddSubAccount
func (c *Client) AddSubAccount(s *Subaccount) *Response {
	return c.AddSubAccountContext(context.Background(), s)
}

// AddSubAccountContext is the same as AddSubAccount, and it allows the caller to pass in a context
func (c *Client) AddSubAccountContext(ctx context.Context, s *Subaccount) *Response {
	if s.ConfirmPassword == "" {
		s.ConfirmPassword = s.Password
	}

	res := c.HTTPGet(ctx, "account/addsubaccount", s)
	if res.Success {
		s.APIKey = res.Data.(string)
	}
	return res
}

// DeleteSubAccount deletes specified Subaccount
// one of subAccountEmail or publicAccountID must be provided
// https://api.elasticemail.com/public/help#Account_DeleteSubAccount
func (c *Client) DeleteSubAccount(params map[string]string) *Response {
	return c.DeleteSubAccountContext(context.Background(), params)
}

// DeleteSubAccountContext is the same as DeleteSubAccount, and it allows the caller to pass in a context
func (c *Client) DeleteSubAccountContext(ctx context.Context, params map[string]string) *Response {
	_, emailPrs := params["subAccountEmail"]
	_, IDPrs := params["publicAccountID"]
	if !emailPrs && !IDPrs {
		return &Response{Error: errors.New("DeleteSubAccount called without email or ID")}
	}

	return c.HTTPGet(ctx, "account/deletesubaccount", params)
}

// UpdateSubAccountSettings updates fields from Subaccount object of specified subaccount referenced by ID from params map
// https://api.elasticemail.com/public/help#Account_UpdateSubAccountSettings
func (c *Client) UpdateSubAccountSettings(s *Subaccount, params map[string]string) *Response {
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

	return c.HTTPGet(ctx, "account/updatesubaccountsettings", rawParams)
}

// GetSubAccountAPIKey attempts to retrieve subaccount API Key
// one of subAccountEmail or publicAccountID must be provided
// https://api.elasticemail.com/public/help#Account_GetSubAccountApiKey
func (c *Client) GetSubAccountAPIKey(params map[string]string) *Response {
	return c.GetSubAccountAPIKeyContext(context.Background(), params)
}

// GetSubAccountAPIKeyContext is the same as GetSubAccountAPIKey, and it allows the caller to pass in a context
func (c *Client) GetSubAccountAPIKeyContext(ctx context.Context, params map[string]string) *Response {
	_, emailPrs := params["subAccountEmail"]
	_, IDPrs := params["publicAccountID"]
	if !emailPrs && !IDPrs {
		return &Response{Error: errors.New("GetSubAccountAPIKey called without email or ID")}
	}

	return c.HTTPGet(ctx, "account/getsubaccountapikey", params)
}
