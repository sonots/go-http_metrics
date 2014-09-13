package http_metrics

import (
	"net/http"
	"time"
)

// print infomation on each request
var Verbose = false

// Set Enable = false to turn off instrumentation
var Enable = true

// a set of proxies
var proxyHandlerFuncRegistry = make(map[*tHttpHandlerFunc](*HandlerFunc))
var proxyHandlerRegistry = make(map[*http.Handler](*Handler))

// a set of metrics
var metricsRegistry = make(map[string](*Metrics))

//WrapHandlerFunc  instrument HTTP handler functions to collect HTTP metrics
func WrapHandlerFunc(name string, h tHttpHandlerFunc) tHttpHandlerFunc {
	metrics := metricsRegistry[name]
	if metrics == nil {
		metrics = newMetrics(name)
		metricsRegistry[name] = metrics
	}
	proxy := proxyHandlerFuncRegistry[&h]
	if proxy == nil {
		proxy = newHandlerFunc(h, metrics)
		proxyHandlerFuncRegistry[&h] = proxy
	}
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

//WrapHandler  instrument HTTP handler object to collect HTTP metrics
func WrapHandler(name string, h http.Handler) http.Handler {
	metrics := metricsRegistry[name]
	if metrics == nil {
		metrics = newMetrics(name)
		metricsRegistry[name] = metrics
	}
	proxy := proxyHandlerRegistry[&h]
	if proxy == nil {
		proxy = newHandler(h, metrics)
		proxyHandlerRegistry[&h] = proxy
	}
	return proxy
}

//Print  print the metrics in each specified second
func Print(duration int) {
	timeDuration := time.Duration(duration)
	go func() {
		time.Sleep(timeDuration * time.Second)
		for {
			startTime := time.Now()
			for _, metrics := range metricsRegistry {
				metrics.printMetrics(duration)
			}
			elapsedTime := time.Now().Sub(startTime)
			time.Sleep(timeDuration*time.Second - elapsedTime)
		}
	}()
}
