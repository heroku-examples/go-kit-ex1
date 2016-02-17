package main

import (
	_ "expvar"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/expvar"
	kitratelimit "github.com/go-kit/kit/ratelimit"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/juju/ratelimit"
	"golang.org/x/net/context"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stdout)

	ctx := context.Background()

	c := &countService{}

	var svc endpoint.Endpoint
	svc = makeAddEndpoint(c)

	limit := ratelimit.NewBucket(2*time.Second, 1)
	svc = kitratelimit.NewTokenBucketLimiter(limit)(svc)

	requestCount := expvar.NewCounter("request.count")
	svc = metricsMiddleware(requestCount)(svc)
	svc = loggingMiddlware(logger)(svc)

	addHandler := httptransport.NewServer(
		ctx,
		svc,
		decodeAddRequest,
		encodeResponse,
		httptransport.ServerBefore(beforeIDExtractor, beforePATHExtractor),
	)

	http.Handle("/add", addHandler)

	port := os.Getenv("PORT")
	logger.Log("listening on", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Log("listen.error", err)
	}
}
