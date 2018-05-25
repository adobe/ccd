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

package ccd

import (
	"log"

	"github.com/adobe/ccd/libs/config"
	"github.com/adobe/ccd/libs/context"
	"github.com/adobe/ccd/libs/services/sync"
	"github.com/pkg/errors"
	"github.com/adobe/ccd/internal/chttp"
)

func Run(l *log.Logger, cfg *config.Config) error {
	if cfg.Username == "" {
		return errors.Errorf("no username configured (in %s)", cfg.Path())
	}

	if cfg.Password == "" {
		return errors.Errorf("no password configured (in %s)", cfg.Path())
	}

	ctx := context.New(l, cfg, chttp.NewClient())

	remoteEvents, err := sync.RunEventLoop(ctx)
	if err != nil {
		return err
	}

	for ev := range remoteEvents {
		ctx.L.Printf("got event: %s", ev)
	}

	return nil
}
