package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// HTTPTestClient provides utilities for testing HTTP handlers.
type HTTPTestClient struct {
	t       *testing.T
	handler http.Handler
	headers map[string]string
}

// NewHTTPTestClient creates a new HTTP test client with the given handler.
func NewHTTPTestClient(t *testing.T, handler http.Handler) *HTTPTestClient {
	return &HTTPTestClient{
		t:       t,
		handler: handler,
		headers: make(map[string]string),
	}
}

// WithHeader adds a header to all subsequent requests.
func (c *HTTPTestClient) WithHeader(key, value string) *HTTPTestClient {
	c.headers[key] = value
	return c
}

// WithTenant sets the tenant ID header for multi-tenant testing.
func (c *HTTPTestClient) WithTenant(tenantID string) *HTTPTestClient {
	return c.WithHeader("X-Tenant-Id", tenantID)
}

// WithSchool sets the school ID header.
func (c *HTTPTestClient) WithSchool(schoolID string) *HTTPTestClient {
	return c.WithHeader("X-School-Id", schoolID)
}

// WithAuth sets the authorization header.
func (c *HTTPTestClient) WithAuth(token string) *HTTPTestClient {
	return c.WithHeader("Authorization", "Bearer "+token)
}

// Get performs a GET request and returns the response.
func (c *HTTPTestClient) Get(path string) *HTTPResponse {
	c.t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	return c.do(req)
}

// Post performs a POST request with a JSON body.
func (c *HTTPTestClient) Post(path string, body interface{}) *HTTPResponse {
	c.t.Helper()
	jsonBody, err := json.Marshal(body)
	require.NoError(c.t, err, "failed to marshal request body")

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return c.do(req)
}

// Put performs a PUT request with a JSON body.
func (c *HTTPTestClient) Put(path string, body interface{}) *HTTPResponse {
	c.t.Helper()
	jsonBody, err := json.Marshal(body)
	require.NoError(c.t, err, "failed to marshal request body")

	req := httptest.NewRequest(http.MethodPut, path, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return c.do(req)
}

// Patch performs a PATCH request with a JSON body.
func (c *HTTPTestClient) Patch(path string, body interface{}) *HTTPResponse {
	c.t.Helper()
	jsonBody, err := json.Marshal(body)
	require.NoError(c.t, err, "failed to marshal request body")

	req := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return c.do(req)
}

// Delete performs a DELETE request.
func (c *HTTPTestClient) Delete(path string) *HTTPResponse {
	c.t.Helper()
	req := httptest.NewRequest(http.MethodDelete, path, nil)
	return c.do(req)
}

// do executes the request with all configured headers.
func (c *HTTPTestClient) do(req *http.Request) *HTTPResponse {
	c.t.Helper()

	// Add configured headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Record the response
	rr := httptest.NewRecorder()
	c.handler.ServeHTTP(rr, req)

	return &HTTPResponse{
		t:        c.t,
		recorder: rr,
	}
}

// HTTPResponse wraps httptest.ResponseRecorder with helper methods.
type HTTPResponse struct {
	t        *testing.T
	recorder *httptest.ResponseRecorder
}

// StatusCode returns the HTTP status code.
func (r *HTTPResponse) StatusCode() int {
	return r.recorder.Code
}

// Body returns the response body as a string.
func (r *HTTPResponse) Body() string {
	return r.recorder.Body.String()
}

// BodyBytes returns the response body as bytes.
func (r *HTTPResponse) BodyBytes() []byte {
	return r.recorder.Body.Bytes()
}

// AssertStatus asserts that the status code matches the expected value.
func (r *HTTPResponse) AssertStatus(expected int) *HTTPResponse {
	r.t.Helper()
	require.Equal(r.t, expected, r.recorder.Code, "unexpected status code")
	return r
}

// AssertJSON unmarshals the response body into the provided value and asserts no error.
func (r *HTTPResponse) AssertJSON(v interface{}) *HTTPResponse {
	r.t.Helper()
	err := json.Unmarshal(r.recorder.Body.Bytes(), v)
	require.NoError(r.t, err, "failed to unmarshal response body")
	return r
}

// AssertJSONField checks that a specific JSON field matches the expected value.
func (r *HTTPResponse) AssertJSONField(field string, expected interface{}) *HTTPResponse {
	r.t.Helper()
	var result map[string]interface{}
	err := json.Unmarshal(r.recorder.Body.Bytes(), &result)
	require.NoError(r.t, err, "failed to unmarshal response body")
	require.Equal(r.t, expected, result[field], "field %s mismatch", field)
	return r
}

// AssertHeader asserts that a header matches the expected value.
func (r *HTTPResponse) AssertHeader(key, expected string) *HTTPResponse {
	r.t.Helper()
	actual := r.recorder.Header().Get(key)
	require.Equal(r.t, expected, actual, "header %s mismatch", key)
	return r
}

// AssertContentType asserts the Content-Type header.
func (r *HTTPResponse) AssertContentType(expected string) *HTTPResponse {
	r.t.Helper()
	return r.AssertHeader("Content-Type", expected)
}

// AssertBodyContains asserts that the response body contains a substring.
func (r *HTTPResponse) AssertBodyContains(substr string) *HTTPResponse {
	r.t.Helper()
	body := r.recorder.Body.String()
	require.Contains(r.t, body, substr, "response body does not contain expected substring")
	return r
}

// GetJSON unmarshals the response body into the provided value.
func (r *HTTPResponse) GetJSON(v interface{}) {
	r.t.Helper()
	err := json.Unmarshal(r.recorder.Body.Bytes(), v)
	require.NoError(r.t, err, "failed to unmarshal response body")
}

// Recorder returns the underlying httptest.ResponseRecorder.
func (r *HTTPResponse) Recorder() *httptest.ResponseRecorder {
	return r.recorder
}

// NewRequest creates a new HTTP request with optional body for manual testing.
func NewRequest(t *testing.T, method, path string, body io.Reader) *http.Request {
	t.Helper()
	req := httptest.NewRequest(method, path, body)
	return req
}

// NewRequestWithContext creates a new HTTP request with context.
func NewRequestWithContext(t *testing.T, ctx context.Context, method, path string, body io.Reader) *http.Request {
	t.Helper()
	req := httptest.NewRequest(method, path, body)
	return req.WithContext(ctx)
}

// NewJSONRequest creates a new HTTP request with a JSON body.
func NewJSONRequest(t *testing.T, method, path string, body interface{}) *http.Request {
	t.Helper()
	jsonBody, err := json.Marshal(body)
	require.NoError(t, err, "failed to marshal request body")

	req := httptest.NewRequest(method, path, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return req
}
