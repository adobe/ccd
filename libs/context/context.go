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

package context

import (
	"log"

	"github.com/adobe/ccd/internal/chttp"
	"github.com/adobe/ccd/libs/config"
)

// Context is used to collect common state and not have to hand it down to
// function not having to reference them one by one.
type Context struct {
	L      *log.Logger
	Cfg    *config.Config
	Client chttp.HttpClient
}

func New(l *log.Logger, cfg *config.Config, c chttp.HttpClient) *Context {
	return &Context{L: l, Cfg: cfg, Client: c}
}
