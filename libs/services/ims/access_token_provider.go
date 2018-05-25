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

package ims

import (
	"github.com/adobe/ccd/libs/context"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type AccessTokenProvider interface {
	// Return the currently used access token.
	AccessToken() (string, error)
}

func RunAccessTokenProvider(ctx *context.Context) (AccessTokenProvider, error) {
	atp := &accessTokenProvider{ctx: ctx}
	go atp.run()
	return atp, atp.authenticate()
}

type accessTokenProvider struct {
	ctx *context.Context
	cur *authToken
	m   sync.Mutex
}

func (atp *accessTokenProvider) AccessToken() (string, error) {
	switch {
	case atp.cur == nil:
		return "", errors.Errorf("token not yet set")
	case !atp.cur.IsValid():
		return "", errors.Errorf("current token is invalid")
	default:
		return atp.cur.Token, nil
	}
}

func (atp *accessTokenProvider) run() {
	time.Sleep(10*time.Second) // let 1. auth call from constructor progress first
	for {
		err := atp.authenticate()
		if err != nil {
			atp.ctx.L.Printf("failed to authenticate: %s", err)
			time.Sleep(60*time.Second)
			continue
		}

		t := time.NewTimer(atp.cur.Expiry.Sub(time.Now()))
		<-t.C
	}
}

func (atp *accessTokenProvider) authenticate() error {
	atp.m.Lock()
	defer atp.m.Unlock()

	if atp.cur != nil && atp.cur.IsValid() {
		return nil
	}

	r, err := authenticate(atp.ctx)
	if err != nil {
		return err
	}

	atp.cur = r
	atp.ctx.L.Printf("authenticated %s", atp.ctx.Cfg.Username)
	return nil
}
