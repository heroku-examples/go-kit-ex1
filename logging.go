package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

const (
	requestIDKey = iota
	pathKey
)

func beforeIDExtractor(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestIDKey, r.Header.Get("X-Request-Id"))
}

func beforePATHExtractor(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, pathKey, r.URL.EscapedPath())
}

func loggingMiddlware(l log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (result interface{}, err error) {
			var req, resp string

			defer func(b time.Time) {
				l.Log(
					"path", ctx.Value(pathKey),
					"request", req,
					"result", resp,
					"err", err,
					"request_id", ctx.Value(requestIDKey),
					"elapsed", time.Since(b),
				)
			}(time.Now())
			if r, ok := request.(fmt.Stringer); ok {
				req = r.String()
			}
			result, err = next(ctx, request)
			if r, ok := result.(fmt.Stringer); ok {
				resp = r.String()
			}
			return
		}
	}
}
