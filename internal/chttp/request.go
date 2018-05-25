/*
 * Copyright 2017 Adobe.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *          http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package chttp

import (
	"net/http"
)

type RequestOption func(req *http.Request) error

func WithHeader(key, value string) RequestOption {
	return func(req *http.Request) error {
		req.Header.Set(key, value)
		return nil
	}
}

func WithParameters(params map[string]string) RequestOption {
	return func(req *http.Request) error {
		q := req.URL.Query()
		if params != nil {
			for k, v := range params {
				q.Add(k, v)
			}
		}
		req.URL.RawQuery = q.Encode()
		return nil
	}
}

const HeaderAuth = "X-User-Token"

func WithAuthToken(token string) RequestOption {
	return WithHeader(HeaderAuth, "Bearer "+token)
}

const HeaderRequestID = "X-Request-Id"
func WithRequestID() RequestOption {
	return func(req *http.Request) error {
		reqId, err := uuid()
		if err != nil {
			return err
		}

		req.Header.Set(HeaderRequestID, reqId)
		return nil
	}
}
