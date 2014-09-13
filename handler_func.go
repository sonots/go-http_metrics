package http_metrics

import (
	"net/http"
	"time"
)

type tHttpHandlerFunc func(http.ResponseWriter, *http.Request)

type HandlerFunc struct {
	Original tHttpHandlerFunc
	metrics  *Metrics
}

func newHandlerFunc(h tHttpHandlerFunc, metrics *Metrics) *HandlerFunc {
	return &HandlerFunc{
		h,
		metrics,
	}
}

func (proxy *HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if Enable {
		startTime := time.Now()
		defer proxy.metrics.measure(startTime, r)
	}
	proxy.Original(w, r)
}
