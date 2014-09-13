package http_metrics

import (
	"net/http"
	"time"
)

type Handler struct {
	Original http.Handler
	*Metrics
}

func newHandler(name string, h http.Handler) *Handler {
	return &Handler{
		h,
		newMetrics(name),
	}
}

func (proxy *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var startTime time.Time
	if Enable {
		startTime = time.Now()
	}
	proxy.Original.ServeHTTP(w, r)
	if Enable {
		defer proxy.measure(startTime, r)
	}
}
