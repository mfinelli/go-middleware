# go-middleware

A collection of generic go middlewares.

## usage

### requestid

The `requestid` middleware is heavily inspired by the
[gin](https://gin-gonic.com)
[version](https://github.com/gin-contrib/requestid). It also allows the
possibility to override the generation function (something which is helpful
during testing for example).

```go
import "go.finelli.dev/middleware/requestid"

// ...

        &http.Server{
                Handler: requestid.New()(mux),
        }

        // or with customization

        &http.Server{
                Handler: requestid.New(requestid.WithHeader("X-Req-ID"))(mux),
        }

// ...
```

You can then retrieve the id in your routes:

```go
import "go.finelli.dev/middleware/requestid"

func Route(w http.ResponseWriter, r *http.Request) {
        rid := requestid.Get(r.Context())
        // ...
}
```

## license

```
Copyright 2024 Mario Finelli

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
