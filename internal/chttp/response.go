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
	"io"
	"net/http"
)

type Response interface {
	Headers() map[string][]string
	Body() io.ReadCloser
	StatusCode() int
}

// Response context object represents the response of a service request. The
// operations should return instances of this type.
type response struct {
	resp *http.Response
}

func (r *response) Headers() map[string][]string {
	return r.resp.Header
}

func (r *response) Body() io.ReadCloser {
	return r.resp.Body
}

func (r *response) StatusCode() int {
	return r.resp.StatusCode
}
