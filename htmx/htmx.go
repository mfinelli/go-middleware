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

// The htmx package checks for the presence of the HX-Request header and its
// expected value "true".
package htmx

import (
	"context"
	"net/http"
)

type key int

const ctxKey key = iota
const header = "HX-Request"

// Middleware to inject the result of checking if it's an HTMX reuest or not.
func CheckHTMX(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxKey,
			r.Header.Get(header) == "true")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Given the context returns the status of the HTMX that we injected with the
// middleware.
func IsHTMX(ctx context.Context) bool {
	if ctx == nil {
		return false
	}

	if isHtmx, ok := ctx.Value(ctxKey).(bool); ok {
		return isHtmx
	}

	return false
}
