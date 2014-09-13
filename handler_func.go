package http_metrics

import (
	"net/http"
	"time"
)

type tHttpHandlerFunc func(http.ResponseWriter, *http.Request)

type HandlerFunc struct {
	Original tHttpHandlerFunc
	*Metrics
}

func newHandlerFunc(name string, h tHttpHandlerFunc) *HandlerFunc {
	return &HandlerFunc{
		h,
		newMetrics(name),
	}
}

func (proxy *HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var startTime time.Time
	if Enable {
		startTime = time.Now()
	}
	proxy.Original(w, r)
	if Enable {
		defer proxy.measure(startTime, r)
	}
}
