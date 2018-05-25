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
	"time"

	"github.com/adobe/ccd/libs/context"
	"github.com/adobe/ccd/libs/events"
	"github.com/adobe/ccd/libs/services/ims"
)

func RunEventLoop(ctx *context.Context) (<-chan events.Event, error) {
	ctx.L.Printf("Starting access token provider")
	atp, err := ims.RunAccessTokenProvider(ctx)
	if err != nil {
		return nil, err
	}

	eC := make(chan events.Event)
	go run(ctx, atp, eC)
	return eC, nil
}

func run(ctx *context.Context, atp ims.AccessTokenProvider, ec chan<- events.Event) {
	for {
		err := func() error {
			// NewDiscovery call.
			discovery, err := NewDiscovery(ctx, atp)
			if err != nil {
				return err
			}

			notif, err := discovery.Notification(ctx, atp)
			if err != nil {
				return err
			}

			for { // Another loop to not call for token and discovery every time
				var journal *Journal
				for journal == nil {
					notif, err = notif.Next(ctx, atp)
					if err != nil {
						return err
					}

					journal, err = notif.Journal(ctx, atp)
					if err != nil {
						return err
					}
				}

				ctx.L.Printf("notification received; journal at %s", journal)
				// TODO: Handle Notification
			}
			return nil
		}()
		if err != nil {
			ctx.L.Printf("ERR: %s", err)
			time.Sleep(20 * time.Second)
		}
	}
}
