# OpenTelemetry handler

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/go-slog/otelslog/ci.yaml?style=flat-square)](https://github.com/go-slog/otelslog/actions/workflows/ci.yaml)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/github.com/go-slog/otelslog)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.21-61CFDD.svg?style=flat-square)
[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

**[log/slog](https://pkg.go.dev/log/slog) handler attaching OpenTelemetry trace details to logs.**


## Installation

```shell
go get github.com/go-slog/otelslog
```

## Usage

Wrap an existing handler:

```go
import(
    "log/slog"

    "github.com/go-slog/otelslog"
)

var handler slog.Handler

// Set up your handler
// handler = ...

// Wrap Handler
handler = otelslog.NewHandler(handler)

logger := slog.New(handler)

// Call logger with a context

logger.InfoContext(ctx, "hello world")

// Output: level=INFO msg="hello world" trace_id=74726163655f69645f74657374313233 span_id=7370616e5f696431
```

Use it as a middleware in [slogmulti.Pipe](https://pkg.go.dev/github.com/samber/slog-multi#Pipe):

```go
import (
    "github.com/go-slog/otelslog"
    "github.com/samber/slog-multi"
)

handler = slogmulti.Pipe(otelslog.Middleware()).Handler(handler)

// Set p logger
// ...
```

## Development

Run tests:

```shell
go test -race -v ./...
```

Run linter:

```shell
golangci-lint run
```

## License

The project is licensed under the [MIT License](LICENSE).
