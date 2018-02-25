package main

import (
	_ "expvar"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/expvar"
	krl "github.com/go-kit/kit/ratelimit"
	kth "github.com/go-kit/kit/transport/http"
	"github.com/juju/ratelimit"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stdout)

	var c countService
	svc := makeAddEndpoint(&c)

	limit := ratelimit.NewBucket(2*time.Second, 1)
	svc = krl.NewTokenBucketLimiter(limit)(svc)

	requestCount := expvar.NewCounter("request.count")
	svc = metricsMiddleware(requestCount)(svc)
	svc = loggingMiddlware(logger)(svc)

	http.Handle("/add",
		kth.NewServer(
			svc,
			decodeAddRequest,
			encodeResponse,
			kth.ServerBefore(beforeIDExtractor, beforePATHExtractor),
		),
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Log("listening-on", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Log("listen.error", err)
	}
}
