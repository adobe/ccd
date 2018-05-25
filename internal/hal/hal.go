// Copyright 2017 Adobe.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//          http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hal

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type href string

func (h *href) UnmarshalJSON(b []byte) error {
	r := struct{
		Link string `json:"href"`
	}{}

	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	*h = href(r.Link)
	return nil
}

type Object struct {
	Links map[string]href `json:"_links"`
}

func New(name, link string) *Object {
	o := new(Object)
	o.Links = map[string]href{name: href(link)}
	return o
}

func (h *Object) Link(name string) (string, error) {
	if h == nil || h.Links == nil {
		return "", errors.Errorf("links not initialized")
	}
	link, found := h.Links[name]
	if !found {
		return "", errors.Errorf("no link %q found", name)
	}
	return string(link), nil
}
