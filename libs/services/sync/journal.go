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
	"github.com/adobe/ccd/internal/hal"
	"github.com/adobe/ccd/libs/context"
	"github.com/adobe/ccd/libs/services/ims"
)

type Journal struct {
	*hal.Object
}

func NewJournal(ctx *context.Context, atp ims.AccessTokenProvider, link string) (*Journal, error) {
	// TODO receive journal
	return &Journal{hal.New("self", link)}, nil
}

func (j *Journal) String() string {
	l, _ := j.Link("self")
	return l
}
