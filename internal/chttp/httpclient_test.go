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

// +build all

package chttp

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jgroeneveld/trial/assert"
)

func TestHttpClient_Get(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(HealthCheck))
	defer s.Close()

	client := NewClient()
	resp, err := client.Get(s.URL)
	assert.Nil(t, err, "error must be nil")
	assert.Equal(t, 200, resp.StatusCode(), "expected status code 200")
}

func TestHttpClient_GetWithParam(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(ParamCheck))
	defer s.Close()

	client := NewClient()
	paramsA := map[string]string{"param": "a"}
	paramsB := map[string]string{"param": "b"}

	{
		resp, err := client.Get(s.URL+"/", WithParameters(paramsA))
		assert.Nil(t, err, "error must be nil")
		defer resp.Body().Close()
		assert.Equal(t, "a", readResponseBody(t, resp), "expected param to have value 'a'")
	}
	{
		resp, err := client.Get(s.URL+"/", WithParameters(paramsB))
		assert.Nil(t, err, "error must be nil")
		defer resp.Body().Close()
		assert.Equal(t, "b", readResponseBody(t, resp), "expected param to have value 'b'")
	}
	{
		resp, err := client.Get(s.URL + "/")
		assert.Nil(t, err, "error must be nil")
		defer resp.Body().Close()
		assert.Equal(t, "", readResponseBody(t, resp), "expected param to be empty")
	}
}

func TestHttpClient_GetWithHeader(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(HeaderCheck))
	defer s.Close()

	client := NewClient()
	{
		resp, err := client.Get(s.URL+"/", WithHeader("X-HeaderTest", "a"))
		assert.Nil(t, err, "error must be nil")
		defer resp.Body().Close()
		assert.Equal(t, "a", readResponseBody(t, resp), "expected header X-HeaderTest to have value 'a'")
	}
	{
		resp, err := client.Get(s.URL+"/", WithHeader("X-HeaderTest", "b"))
		assert.Nil(t, err, "error must be nil")
		defer resp.Body().Close()
		assert.Equal(t, "b", readResponseBody(t, resp), "expected header X-HeaderTest to have value 'b'")
	}
	{
		resp, err := client.Get(s.URL + "/")
		assert.Nil(t, err, "error must be nil")
		defer resp.Body().Close()
		assert.Equal(t, "", readResponseBody(t, resp), "expected header X-HeaderTest to be empty")
	}
}

func readResponseBody(t *testing.T, r Response) string {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, r.Body()); err != nil {
		t.Fatalf("failed to read response body: %s", err)
	}
	return buf.String()
}

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("{ 'healthy' : true }"))
}

func ParamCheck(w http.ResponseWriter, r *http.Request) {
	val := r.URL.Query().Get("param")
	io.WriteString(w, val)
}

func HeaderCheck(w http.ResponseWriter, r *http.Request) {
	val := r.Header.Get("X-HeaderTest")
	io.WriteString(w, val)
}
