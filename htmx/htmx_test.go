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
	tests := []struct {
		name        string
		setHeader   bool
		headerValue string
		expected    bool
	}{
		{
			name:        "No HX-Request header set",
			setHeader:   false,
			headerValue: "",
			expected:    false,
		},
		{
			name:        "HX-Request set to empty value",
			setHeader:   true,
			headerValue: "",
			expected:    false,
		},
		{
			name:        "HX-Request header set but not true",
			setHeader:   true,
			headerValue: "nottrue",
			expected:    false,
		},
		{
			name:        "HX-Request header set",
			setHeader:   true,
			headerValue: "true",
			expected:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, e := http.NewRequest("GET", "/", nil)
			assert.Nil(t, e)

			if test.setHeader {
				r.Header.Set("HX-Request", test.headerValue)
			}

			w := httptest.NewRecorder()
			h := CheckHTMX(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expected, IsHTMX(r.Context()))
				w.WriteHeader(http.StatusOK)
			}))
			h.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
