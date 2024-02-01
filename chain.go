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

// The middleware package provides middleware helpers.
package middleware

import "net/http"

type middleware = func(http.Handler) http.Handler

// Wraps the given handler with the provided middlewares.
func Chain(h http.Handler, middlewares ...middleware) http.Handler {
	switch len(middlewares) {
	case 0:
		return h

	case 1:
		return middlewares[0](h)

	default:
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}

		return h
	}
}
