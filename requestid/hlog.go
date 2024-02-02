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

package requestid

import (
	"net/http"

	zl "github.com/rs/zerolog"
)

// LogHandler adds a zerolog/hlog middlewar handler to inject the request id
// into the log context with the desired key.
func LogHandler(key string) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()

				if rid := Get(ctx); rid != "" {
					log := zl.Ctx(ctx)
					log.UpdateContext(
						func(c zl.Context) zl.Context {
							return c.Str(key, rid)
						},
					)
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}
