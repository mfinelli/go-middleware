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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func emptyResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func TestRequestId(t *testing.T) {
	r, e := http.NewRequest("GET", "/", nil)
	assert.Nil(t, e)

	w := httptest.NewRecorder()
	h := New()(http.HandlerFunc(emptyResponse))
	h.ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.NotEmpty(t, w.Header().Get(defaultHeaderKey))
}

func TestRequestIdPassThru(t *testing.T) {
	r, e := http.NewRequest("GET", "/", nil)
	assert.Nil(t, e)
	r.Header.Set(defaultHeaderKey, "ABC")

	w := httptest.NewRecorder()
	h := New()(http.HandlerFunc(emptyResponse))
	h.ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Header().Get(defaultHeaderKey), "ABC")
}

func TestRequestIdCustomHeader(t *testing.T) {
	r, e := http.NewRequest("GET", "/", nil)
	assert.Nil(t, e)

	w := httptest.NewRecorder()
	h := New(WithHeader("X-Custom-ID"))(http.HandlerFunc(emptyResponse))
	h.ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.NotEmpty(t, w.Header().Get("X-Custom-ID"))
	assert.Empty(t, w.Header().Get(defaultHeaderKey))
}

func TestRequestIdCustomGenerator(t *testing.T) {
	r, e := http.NewRequest("GET", "/", nil)
	assert.Nil(t, e)

	w := httptest.NewRecorder()
	h := New(WithGenerator(func() string {
		return "custom"
	}))(http.HandlerFunc(emptyResponse))
	h.ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Header().Get(defaultHeaderKey), "custom")
}

func TestRequestIdGet(t *testing.T) {
	r, e := http.NewRequest("GET", "/", nil)
	assert.Nil(t, e)

	w := httptest.NewRecorder()
	h := New(WithGenerator(func() string {
		return "test"
	}))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := Get(r.Context())
		assert.Equal(t, rid, "test")

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}))
	h.ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusOK)
}
