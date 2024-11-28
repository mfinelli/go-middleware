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

package htmx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckHTMX(t *testing.T) {
	t.Run("No HX-Request header set", func(t *testing.T) {
		r, e := http.NewRequest("GET", "/", nil)
		assert.Nil(t, e)

		w := httptest.NewRecorder()
		h := CheckHTMX(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isHtmx := IsHTMX(r.Context())
			assert.Equal(t, false, isHtmx)
			w.WriteHeader(http.StatusOK)
		}))
		h.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("HX-Request header set but not true", func(t *testing.T) {
		r, e := http.NewRequest("GET", "/", nil)
		assert.Nil(t, e)
		r.Header.Set("HX-Request", "")

		w := httptest.NewRecorder()
		h := CheckHTMX(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isHtmx := IsHTMX(r.Context())
			assert.Equal(t, false, isHtmx)
			w.WriteHeader(http.StatusOK)
		}))
		h.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("HX-Request header set", func(t *testing.T) {
		r, e := http.NewRequest("GET", "/", nil)
		assert.Nil(t, e)
		r.Header.Set("HX-Request", "true")

		w := httptest.NewRecorder()
		h := CheckHTMX(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isHtmx := IsHTMX(r.Context())
			assert.Equal(t, true, isHtmx)
			w.WriteHeader(http.StatusOK)
		}))
		h.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
