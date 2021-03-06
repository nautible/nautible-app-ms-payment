// Package creditclient provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.0 DO NOT EDIT.
package creditclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// Create request with any body
	CreateWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	Create(ctx context.Context, body CreateJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Update request with any body
	UpdateWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	Update(ctx context.Context, body UpdateJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Delete request
	Delete(ctx context.Context, acceptNo string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetByAcceptNo request
	GetByAcceptNo(ctx context.Context, acceptNo string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Healthz request
	Healthz(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) CreateWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Create(ctx context.Context, body CreateJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Update(ctx context.Context, body UpdateJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Delete(ctx context.Context, acceptNo string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteRequest(c.Server, acceptNo)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetByAcceptNo(ctx context.Context, acceptNo string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetByAcceptNoRequest(c.Server, acceptNo)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Healthz(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewHealthzRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewCreateRequest calls the generic Create builder with application/json body
func NewCreateRequest(server string, body CreateJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateRequestWithBody generates requests for Create with any type of body
func NewCreateRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/credit")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewUpdateRequest calls the generic Update builder with application/json body
func NewUpdateRequest(server string, body UpdateJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateRequestWithBody(server, "application/json", bodyReader)
}

// NewUpdateRequestWithBody generates requests for Update with any type of body
func NewUpdateRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/credit")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteRequest generates requests for Delete
func NewDeleteRequest(server string, acceptNo string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "acceptNo", runtime.ParamLocationPath, acceptNo)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/credit/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetByAcceptNoRequest generates requests for GetByAcceptNo
func NewGetByAcceptNoRequest(server string, acceptNo string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "acceptNo", runtime.ParamLocationPath, acceptNo)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/credit/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewHealthzRequest generates requests for Healthz
func NewHealthzRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/healthz")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// Create request with any body
	CreateWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateResponse, error)

	CreateWithResponse(ctx context.Context, body CreateJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateResponse, error)

	// Update request with any body
	UpdateWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateResponse, error)

	UpdateWithResponse(ctx context.Context, body UpdateJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateResponse, error)

	// Delete request
	DeleteWithResponse(ctx context.Context, acceptNo string, reqEditors ...RequestEditorFn) (*DeleteResponse, error)

	// GetByAcceptNo request
	GetByAcceptNoWithResponse(ctx context.Context, acceptNo string, reqEditors ...RequestEditorFn) (*GetByAcceptNoResponse, error)

	// Healthz request
	HealthzWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthzResponse, error)
}

type CreateResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *RestCreditPayment
}

// Status returns HTTPResponse.Status
func (r CreateResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *RestCreditPayment
}

// Status returns HTTPResponse.Status
func (r UpdateResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r DeleteResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetByAcceptNoResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *RestCreditPayment
}

// Status returns HTTPResponse.Status
func (r GetByAcceptNoResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetByAcceptNoResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type HealthzResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r HealthzResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r HealthzResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// CreateWithBodyWithResponse request with arbitrary body returning *CreateResponse
func (c *ClientWithResponses) CreateWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateResponse, error) {
	rsp, err := c.CreateWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateResponse(rsp)
}

func (c *ClientWithResponses) CreateWithResponse(ctx context.Context, body CreateJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateResponse, error) {
	rsp, err := c.Create(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateResponse(rsp)
}

// UpdateWithBodyWithResponse request with arbitrary body returning *UpdateResponse
func (c *ClientWithResponses) UpdateWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateResponse, error) {
	rsp, err := c.UpdateWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateResponse(rsp)
}

func (c *ClientWithResponses) UpdateWithResponse(ctx context.Context, body UpdateJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateResponse, error) {
	rsp, err := c.Update(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateResponse(rsp)
}

// DeleteWithResponse request returning *DeleteResponse
func (c *ClientWithResponses) DeleteWithResponse(ctx context.Context, acceptNo string, reqEditors ...RequestEditorFn) (*DeleteResponse, error) {
	rsp, err := c.Delete(ctx, acceptNo, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteResponse(rsp)
}

// GetByAcceptNoWithResponse request returning *GetByAcceptNoResponse
func (c *ClientWithResponses) GetByAcceptNoWithResponse(ctx context.Context, acceptNo string, reqEditors ...RequestEditorFn) (*GetByAcceptNoResponse, error) {
	rsp, err := c.GetByAcceptNo(ctx, acceptNo, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetByAcceptNoResponse(rsp)
}

// HealthzWithResponse request returning *HealthzResponse
func (c *ClientWithResponses) HealthzWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthzResponse, error) {
	rsp, err := c.Healthz(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseHealthzResponse(rsp)
}

// ParseCreateResponse parses an HTTP response from a CreateWithResponse call
func ParseCreateResponse(rsp *http.Response) (*CreateResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest RestCreditPayment
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseUpdateResponse parses an HTTP response from a UpdateWithResponse call
func ParseUpdateResponse(rsp *http.Response) (*UpdateResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest RestCreditPayment
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseDeleteResponse parses an HTTP response from a DeleteWithResponse call
func ParseDeleteResponse(rsp *http.Response) (*DeleteResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetByAcceptNoResponse parses an HTTP response from a GetByAcceptNoWithResponse call
func ParseGetByAcceptNoResponse(rsp *http.Response) (*GetByAcceptNoResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetByAcceptNoResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest RestCreditPayment
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseHealthzResponse parses an HTTP response from a HealthzWithResponse call
func ParseHealthzResponse(rsp *http.Response) (*HealthzResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &HealthzResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}
