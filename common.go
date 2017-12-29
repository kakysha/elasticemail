package elasticemail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
)

// Config includes all information necessary to make an API request.
type Config struct {
	BaseURL    string
	APIKey     string
	APIVersion int
}

// Client contains connection, configuration, and authentication information.
// Specifying your own http.Client gives you lots of control over how connections are made.
// Clients are safe for concurrent (read-only) reuse by multiple goroutines.
// Headers is useful to set custom headers.
// All changes to Headers must happen before Client is exposed to possible concurrent use.
type Client struct {
	Config  *Config
	Client  *http.Client
	Headers *http.Header
}

// Response contains information about the last HTTP response.
// Helpful when an error message doesn't necessarily give the complete picture.
// Also contains any messages emitted as a result of the Verbose config option.
type Response struct {
	httpResponse *http.Response
	Body         []byte
	Success      bool        `json:"success,omitempty"`
	Error        error       `json:"error,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

// UnmarshalJSON to handle string Error and convert it into "error" type
func (r *Response) UnmarshalJSON(data []byte) error {
	var tmpResponse struct {
		Success bool
		Error   string
		Data    interface{}
	}
	err := json.Unmarshal(data, &tmpResponse)
	if err != nil {
		return err
	}

	r.Success = tmpResponse.Success
	r.Data = tmpResponse.Data
	if tmpResponse.Error != "" {
		r.Error = errors.New(tmpResponse.Error)
	}

	return nil
}

// Init pulls together everything necessary to make an API request.
// Caller may provide their own http.Client by setting it in the provided API object.
func (c *Client) Init(cfg *Config) error {
	// Set default values
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.elasticemail.com"
	} else if !strings.HasPrefix(cfg.BaseURL, "https://") {
		return errors.New("API base url must be https")
	}
	if cfg.APIVersion == 0 {
		cfg.APIVersion = 2
	}
	c.Config = cfg
	c.Headers = &http.Header{}
	if c.Client == nil {
		c.Client = http.DefaultClient
	}

	return nil
}

// HTTPPost sends a Post request with the provided JSON payload to the specified url.
// Query params are converted to via net/url.Values
// Authenticate using the configured API key.
func (c *Client) HTTPPost(ctx context.Context, path string, params interface{}) *Response {
	return c.doRequest(ctx, "POST", path, params)
}

// HTTPGet sends a Get request to the specified url.
// Query params are converted to via net/url.Values
// Authenticate using the configured API key.
func (c *Client) HTTPGet(ctx context.Context, path string, params interface{}) *Response {
	return c.doRequest(ctx, "GET", path, params)
}

func (c *Client) doRequest(ctx context.Context, method string, path string, params interface{}) *Response {
	if c == nil {
		return &Response{Error: errors.New("Client must be non-nil")}
	} else if c.Client == nil {
		return &Response{Error: errors.New("Client.Client (http.Client) must be non-nil")}
	} else if c.Config == nil {
		return &Response{Error: errors.New("Client.Config must be non-nil")}
	}

	jsonBytes, _ := json.Marshal(params)
	var rawParams map[string]string
	json.Unmarshal(jsonBytes, &rawParams)

	requestParams := url.Values{}
	for k, v := range rawParams {
		requestParams.Set(k, v)
	}

	if requestParams.Get("apikey") == "" {
		if apikey, ok := ctx.Value("apikey").(string); ok {
			requestParams.Set("apikey", apikey)
		} else {
			requestParams.Set("apikey", c.Config.APIKey)
		}
	}

	var getParams, postParams string
	if method == "GET" {
		getParams = requestParams.Encode()
	} else {
		postParams = requestParams.Encode()
	}

	var urlStr = fmt.Sprintf("%s/v%d/%s?%s", c.Config.BaseURL, c.Config.APIVersion, path, getParams)

	req, err := http.NewRequest(method, urlStr, strings.NewReader(postParams))
	if err != nil {
		return &Response{Error: errors.Wrap(err, "building request")}
	}

	req.Header.Set("User-Agent", fmt.Sprintf("ElasticEmail Go API Client v%d", c.Config.APIVersion))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Forward additional headers set in client to request
	if c.Headers != nil {
		for header, values := range map[string][]string(*c.Headers) {
			for _, value := range values {
				req.Header.Add(header, value)
			}
		}
	}

	if ctx == nil {
		ctx = context.Background()
	}
	// set any headers provided in context
	if header, ok := ctx.Value("http.Header").(http.Header); ok {
		for key, vals := range map[string][]string(header) {
			req.Header.Del(key)
			for _, val := range vals {
				req.Header.Add(key, val)
			}
		}
	}
	req = req.WithContext(ctx)

	res := &Response{}
	res.httpResponse, res.Error = c.Client.Do(req)

	if res.Error != nil {
		return res
	}
	res.parseResponse()
	return res
}

// parseResponse pulls info from JSON http responses into api.Response object.
// It's helpful to call Response.AssertJson before calling this function.
func (r *Response) parseResponse() {
	if r.Body != nil {
		return
	}

	defer r.httpResponse.Body.Close()

	bodyBytes, err := ioutil.ReadAll(r.httpResponse.Body)
	if err != nil {
		r.Error = errors.Wrap(err, "reading http body")
		return
	}
	r.Body = bodyBytes

	// Don't try to unmarshal an empty response
	if bytes.Compare(bodyBytes, []byte("")) == 0 {
		r.Error = errors.New("empty response body")
		return
	}

	ctype := r.httpResponse.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(ctype)
	if err != nil {
		r.Error = errors.Wrap(err, "parsing content-type")
		return
	}
	// allow things like "application/json; charset=utf-8" in addition to the bare content type
	if mediaType != "application/json" {
		r.Error = errors.Errorf("Expected json, got [%s]", mediaType)
		return
	}

	err = json.Unmarshal(r.Body, r)
	if err != nil {
		r.Error = errors.Wrap(err, "parsing json api response")
		return
	}
}
