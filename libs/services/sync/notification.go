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

package sync

import (
	"github.com/adobe/ccd/internal/hal"
	"github.com/adobe/ccd/libs/context"
	"github.com/adobe/ccd/libs/services/ims"
)

type Notification struct {
	*hal.Object
}

func NewNotification(ctx *context.Context, atp ims.AccessTokenProvider, link string) (*Notification, error) {
	next := new(Notification)
	return next, sendGET(ctx, atp, link, &next)
}

func (n *Notification) Next(ctx *context.Context, atp ims.AccessTokenProvider) (*Notification, error) {
	nextLink, err := n.Link("next")
	if err != nil {
		return nil, err
	}

	return NewNotification(ctx, atp, nextLink)
}

func (n *Notification) Journal(ctx *context.Context, atp ims.AccessTokenProvider) (*Journal, error) {
	journalL, err := n.Link("journal")
	if err != nil || journalL == "" {
		return nil, nil
	}
	return NewJournal(ctx, atp, journalL)
}
