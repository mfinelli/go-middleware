/*!
 * Copyright 2024 Mario Finelli
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// The requestid package creates a middleware the generates and injects a
// unique request id into the request context and response header.
package requestid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type (
	generator  func() string
	key        int // used to avoid context collisions
	middleware func(http.Handler) http.Handler
	option     func(*config)
)

type config struct {
	generator generator
	header    string
}

const ctxKey key = iota
var defaultHeaderKey = string("X-Request-ID")

// Returns a new requestid middleware with the requested options.
func New(opts ...option) middleware {
	cfg := &config{
		generator: func() string {
			return uuid.New().String()
		},
		header: defaultHeaderKey,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				rid := r.Header.Get(cfg.header)
				if rid == "" {
					rid = cfg.generator()
				}
				ctx = context.WithValue(ctx, ctxKey, rid)
				w.Header().Set(cfg.header, rid)
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}

// Tries to get the request id from the request context.
func Get(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if rid, ok := ctx.Value(ctxKey).(string); ok {
		return rid
	}

	return ""
}

// Override the request id header from the default "X-Request-ID".
func WithHeader(h string) option {
	return func(cfg *config) {
		cfg.header = h
	}
}

// Override the ID generation function which defaults to a v4 UUID.
func WithGenerator(g generator) option {
	return func(cfg *config) {
		cfg.generator = g
	}
}
