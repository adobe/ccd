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

package sync

import (
	"encoding/json"

	"github.com/adobe/ccd/internal/chttp"
	"github.com/adobe/ccd/libs/context"
	"github.com/adobe/ccd/libs/services/ims"
	"github.com/pkg/errors"
)

func sendGET(ctx *context.Context, atp ims.AccessTokenProvider, link string, res interface{}, opts ...chttp.RequestOption) error {
	token, err := atp.AccessToken()
	if err != nil {
		return err
	}

	baseOpts := []chttp.RequestOption{
		chttp.WithAuthToken(token),
		chttp.WithRequestID(),
		chttp.WithHeader("Accept", "application/hal+json"),
	}
	opts = append(baseOpts, opts...)
	resp, err := ctx.Client.Get(link, opts...)
	if err != nil {
		return err
	}
	defer resp.Body().Close()

	if resp.StatusCode() != 200 {
		return errors.Errorf("wrong status code %d for request", resp.StatusCode())
	}

	err = json.NewDecoder(resp.Body()).Decode(&res)
	return errors.Wrap(err, "failed to unmarshal response")
}
