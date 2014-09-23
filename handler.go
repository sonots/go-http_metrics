package http_metrics

import (
	"net/http"
	"time"
)

type Handler struct {
	Original http.Handler
	metrics  *Metrics
}

func newHandler(h http.Handler, metrics *Metrics) *Handler {
	return &Handler{
		h,
		metrics,
	}
}

func (proxy *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if Enable {
		startTime := time.Now()
		defer proxy.metrics.measure(startTime, r)
		if OnResponse != nil {
			defer OnResponse()
		}
	}
	proxy.Original.ServeHTTP(w, r)
}
