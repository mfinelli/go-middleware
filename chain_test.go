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

package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testHandler struct{}

func (*testHandler) ServeHTTP(http.ResponseWriter, *http.Request) {}

type testMiddleware struct{}

func (*testMiddleware) ServeHTTP(http.ResponseWriter, *http.Request) {}

func TestChain(t *testing.T) {
	endpoint := &testHandler{}

	t.Run("no middleware", func(t *testing.T) {
		r := Chain(endpoint)
		assert.Equal(t, endpoint, r)
	})

	t.Run("one middleware", func(t *testing.T) {
		m := &testMiddleware{}
		r := Chain(endpoint, func(h http.Handler) http.Handler {
			return m
		})
		assert.Equal(t, m, r)
	})

	t.Run("ensure order", func(t *testing.T) {
		var (
			calls []int

			callMiddleware = func(n int) middleware {
				return func(next http.Handler) http.Handler {
					return http.HandlerFunc(
						func(w http.ResponseWriter,
							r *http.Request) {
							calls = append(calls,
								n)
							next.ServeHTTP(w, r)
						})
				}
			}
		)

		r := Chain(endpoint,
			callMiddleware(1),
			callMiddleware(2),
			callMiddleware(3),
		)

		r.ServeHTTP(nil, nil)
		assert.Equal(t, []int{1, 2, 3}, calls)
	})
}
