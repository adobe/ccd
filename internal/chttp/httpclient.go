//
//  Copyright 2017 Adobe.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//          http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package chttp

import (
	"net/http"

	"github.com/pkg/errors"
)

// The `HttpClient` is used in operations executed against the backend service.
// It hides the actual HTTP client implementation details by adding an
// abstraction layer between its clients and the underlying implementation.
type HttpClient interface {
	Get(url string, opts ... RequestOption) (Response, error)
	Head(url string, opts ... RequestOption) (Response, error)
}

// Underlying implementation of the HTTP client. HttpClient type depends
// on the Go standard library.
type httpClient struct {
	Client *http.Client
}

// A factory to create a new HTTP client instance. The configuration provided
// is the daemon configuration. The client should be configurable.
func NewClient() HttpClient {
	c := new(httpClient)
	// This should use a custom client configurable as required.
	c.Client = http.DefaultClient
	return c
}

// Get Request method.
func (h *httpClient) Get(url string, opts ...RequestOption) (Response, error) {
	return h.do(http.MethodGet, url, opts)
}

// Head Request method.
func (h *httpClient) Head(url string, opts ...RequestOption) (Response, error) {
	return h.do(http.MethodHead, url, opts)
}

// Execute the provided request.
func (h *httpClient) do(method, url string, opts []RequestOption) (Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create base request")
	}
	for i := range opts {
		if err := opts[i](req); err != nil {
			return nil, err
		}
	}
	return h.executeRequest(req)
}

// Executes the request that is fully initialised so far.
func (h *httpClient) executeRequest(req *http.Request) (Response, error) {
	res, err := h.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}
	return &response{resp: res}, nil
}
