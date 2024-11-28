# go-middleware

[![Default](https://github.com/mfinelli/go-middleware/actions/workflows/default.yml/badge.svg)](https://github.com/mfinelli/go-middleware/actions/workflows/default.yml)
[![Go Reference](https://pkg.go.dev/badge/go.finelli.dev/middleware.svg)](https://pkg.go.dev/go.finelli.dev/middleware)

A collection of generic go middlewares.

## usage

### utilities

#### middleware chaining

This module provides a function to chain middlewares together to avoid needing
to wrap each middleware around the next which is harder to reason about. It's
unfortunately, not as simple as some web frameworks' `router.Use` functions.
It comes almost directly from the
[example](https://github.com/ogen-go/example/blob/main/internal/httpmiddleware/httpmiddleware.go#L104)
provided by the [ogen](https://ogen.dev) project.

```go
import (
        "net/http"
        mw "go.finelli.dev/middleware"
        "go.finelli.dev/middleware/requestid"
        chi "github.com/go-chi/chi/v5/middleware"
)

// ...

        mux := http.NewServeMux()

        // ...

        &http.Server{
                Handler: mw.Chain(mux,
                        chi.Logger,
                        chi.Recoverer,
                        requestid.New(),
                ),
        }

// ...
```

### middlewares

#### htmx

The `htmx` middleware simply checks if the `HX-Request`
[header](https://htmx.org/docs/#request-headers) has been set in the request
and if it's value is `"true"` (which is always the case when set).

```go
import "go.finelli.dev/middleware/htmx"

// ...

        &http.Server{
                Handler: htmx.CheckHTMX(mux),
        }

// ...
```

Then in a route you can check if the request was sent with HTMX:

```go
import "go.finelli.dev/middleware/htmx"

func Route(w http.ResponseWriter, r *http.Request) {
        isHtmx := htmx.IsHTMX(r.Context())
        // ...
}
```

#### requestid

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

There's also a middleware to automatically inject the request id into a
[zerolog](https://github.com/rs/zerolog) `hlog` context.

```go
import (
        "os"

        "github.com/rs/zerolog"
        "github.com/rs/zerolog/hlog"
        "go.finelli.dev/middleware/requestid"
)

// ...

        log := zerolog.New(os.Stdout).With().Timestamp().Logger()
        mux := http.NewServeMux()

        // ...

        &http.Server{
                Handler: requestid.New()(hlog.NewHandler(log)(
                        hlog.AccessHandler(/*...*/)(
                                requestid.LogHandler("request_id")(mux)))),
        }

// ...

func Route(w http.ResponseWriter, r *http.Request) {
        log := hlog.FromRequest(r)

        // these log messages have the "request_id" key populated now
        log.Info().Msg("test log message")

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
