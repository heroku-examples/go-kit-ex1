package main

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
)

func metricsMiddleware(requestCount metrics.Counter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			requestCount.Add(1)
			return next(ctx, request)
		}
	}
}
