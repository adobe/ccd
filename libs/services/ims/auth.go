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

package ims

import (
	"encoding/json"

	"fmt"
	"io/ioutil"

	"strconv"
	"time"

	"github.com/adobe/ccd/internal/chttp"
	"github.com/adobe/ccd/libs/context"
	"github.com/pkg/errors"
)

type authToken struct {
	Token  string
	Expiry time.Time
}

func newAuthToken(token, expiry string) (*authToken, error) {
	i, err := strconv.Atoi(expiry)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert milliseconds till expiry")
	}

	return &authToken{Token: token, Expiry: time.Now().Add(time.Duration(i) * time.Millisecond)}, nil
}

func (t *authToken) IsValid() bool {
	return time.Now().Before(t.Expiry)
}

// Factory method that creates new AccessToken instances.
func authenticate(ctx *context.Context) (*authToken, error) {
	params := map[string]string{
		"username":  ctx.Cfg.Username,
		"password":  ctx.Cfg.Password,
		"client_id": ctx.Cfg.ClientId,
		"scope":     ctx.Cfg.Scope,
		"locale":    ctx.Cfg.Locale,
	}

	resp, err := ctx.Client.Get(authorizeEndpoint(ctx), chttp.WithParameters(params))
	if err != nil {
		return nil, err
	}
	defer resp.Body().Close()

	if resp.StatusCode() != 200 {
		s, _ := ioutil.ReadAll(resp.Body())
		return nil, errors.Errorf("failed to retrieve auth token from IMS (%d: %s)", resp.StatusCode(), s)
	}

	authResp := struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
	}{}
	if err := json.NewDecoder(resp.Body()).Decode(&authResp); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal IMS response")
	}

	// returns a newly created access token.
	return newAuthToken(authResp.AccessToken, authResp.ExpiresIn)
}

func authorizeEndpoint(ctx *context.Context) string {
	return fmt.Sprintf("%s/ims/login/v1/token", ctx.Cfg.AuthEndpoint)
}
