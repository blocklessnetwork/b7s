// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	// ExecuteFunctionWithBody request with any body
	ExecuteFunctionWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	ExecuteFunction(ctx context.Context, body ExecuteFunctionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// InstallFunctionWithBody request with any body
	InstallFunctionWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	InstallFunction(ctx context.Context, body InstallFunctionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ExecutionResultWithBody request with any body
	ExecutionResultWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	ExecutionResult(ctx context.Context, body ExecutionResultJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Health request
	Health(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) ExecuteFunctionWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewExecuteFunctionRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ExecuteFunction(ctx context.Context, body ExecuteFunctionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewExecuteFunctionRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) InstallFunctionWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewInstallFunctionRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) InstallFunction(ctx context.Context, body InstallFunctionJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewInstallFunctionRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ExecutionResultWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewExecutionResultRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ExecutionResult(ctx context.Context, body ExecutionResultJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewExecutionResultRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Health(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewHealthRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewExecuteFunctionRequest calls the generic ExecuteFunction builder with application/json body
func NewExecuteFunctionRequest(server string, body ExecuteFunctionJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewExecuteFunctionRequestWithBody(server, "application/json", bodyReader)
}

// NewExecuteFunctionRequestWithBody generates requests for ExecuteFunction with any type of body
func NewExecuteFunctionRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/functions/execute")
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

// NewInstallFunctionRequest calls the generic InstallFunction builder with application/json body
func NewInstallFunctionRequest(server string, body InstallFunctionJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewInstallFunctionRequestWithBody(server, "application/json", bodyReader)
}

// NewInstallFunctionRequestWithBody generates requests for InstallFunction with any type of body
func NewInstallFunctionRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/functions/install")
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

// NewExecutionResultRequest calls the generic ExecutionResult builder with application/json body
func NewExecutionResultRequest(server string, body ExecutionResultJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewExecutionResultRequestWithBody(server, "application/json", bodyReader)
}

// NewExecutionResultRequestWithBody generates requests for ExecutionResult with any type of body
func NewExecutionResultRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/functions/requests/result")
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

// NewHealthRequest generates requests for Health
func NewHealthRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v1/health")
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
	// ExecuteFunctionWithBodyWithResponse request with any body
	ExecuteFunctionWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ExecuteFunctionResponse, error)

	ExecuteFunctionWithResponse(ctx context.Context, body ExecuteFunctionJSONRequestBody, reqEditors ...RequestEditorFn) (*ExecuteFunctionResponse, error)

	// InstallFunctionWithBodyWithResponse request with any body
	InstallFunctionWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*InstallFunctionResponse, error)

	InstallFunctionWithResponse(ctx context.Context, body InstallFunctionJSONRequestBody, reqEditors ...RequestEditorFn) (*InstallFunctionResponse, error)

	// ExecutionResultWithBodyWithResponse request with any body
	ExecutionResultWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ExecutionResultResponse, error)

	ExecutionResultWithResponse(ctx context.Context, body ExecutionResultJSONRequestBody, reqEditors ...RequestEditorFn) (*ExecutionResultResponse, error)

	// HealthWithResponse request
	HealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthResponse, error)
}

type ExecuteFunctionResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ExecutionResponse
}

// Status returns HTTPResponse.Status
func (r ExecuteFunctionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ExecuteFunctionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type InstallFunctionResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *FunctionInstallResponse
}

// Status returns HTTPResponse.Status
func (r InstallFunctionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r InstallFunctionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ExecutionResultResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *FunctionResultResponse
}

// Status returns HTTPResponse.Status
func (r ExecutionResultResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ExecutionResultResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type HealthResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *HealthStatus
}

// Status returns HTTPResponse.Status
func (r HealthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r HealthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ExecuteFunctionWithBodyWithResponse request with arbitrary body returning *ExecuteFunctionResponse
func (c *ClientWithResponses) ExecuteFunctionWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ExecuteFunctionResponse, error) {
	rsp, err := c.ExecuteFunctionWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseExecuteFunctionResponse(rsp)
}

func (c *ClientWithResponses) ExecuteFunctionWithResponse(ctx context.Context, body ExecuteFunctionJSONRequestBody, reqEditors ...RequestEditorFn) (*ExecuteFunctionResponse, error) {
	rsp, err := c.ExecuteFunction(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseExecuteFunctionResponse(rsp)
}

// InstallFunctionWithBodyWithResponse request with arbitrary body returning *InstallFunctionResponse
func (c *ClientWithResponses) InstallFunctionWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*InstallFunctionResponse, error) {
	rsp, err := c.InstallFunctionWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseInstallFunctionResponse(rsp)
}

func (c *ClientWithResponses) InstallFunctionWithResponse(ctx context.Context, body InstallFunctionJSONRequestBody, reqEditors ...RequestEditorFn) (*InstallFunctionResponse, error) {
	rsp, err := c.InstallFunction(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseInstallFunctionResponse(rsp)
}

// ExecutionResultWithBodyWithResponse request with arbitrary body returning *ExecutionResultResponse
func (c *ClientWithResponses) ExecutionResultWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*ExecutionResultResponse, error) {
	rsp, err := c.ExecutionResultWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseExecutionResultResponse(rsp)
}

func (c *ClientWithResponses) ExecutionResultWithResponse(ctx context.Context, body ExecutionResultJSONRequestBody, reqEditors ...RequestEditorFn) (*ExecutionResultResponse, error) {
	rsp, err := c.ExecutionResult(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseExecutionResultResponse(rsp)
}

// HealthWithResponse request returning *HealthResponse
func (c *ClientWithResponses) HealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthResponse, error) {
	rsp, err := c.Health(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseHealthResponse(rsp)
}

// ParseExecuteFunctionResponse parses an HTTP response from a ExecuteFunctionWithResponse call
func ParseExecuteFunctionResponse(rsp *http.Response) (*ExecuteFunctionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ExecuteFunctionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ExecutionResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseInstallFunctionResponse parses an HTTP response from a InstallFunctionWithResponse call
func ParseInstallFunctionResponse(rsp *http.Response) (*InstallFunctionResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &InstallFunctionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest FunctionInstallResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseExecutionResultResponse parses an HTTP response from a ExecutionResultWithResponse call
func ParseExecutionResultResponse(rsp *http.Response) (*ExecutionResultResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ExecutionResultResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest FunctionResultResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseHealthResponse parses an HTTP response from a HealthWithResponse call
func ParseHealthResponse(rsp *http.Response) (*HealthResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &HealthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest HealthStatus
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}
