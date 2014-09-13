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
var proxyHandlerFuncRegistry = make(map[string](*HandlerFunc))
var proxyHandlerRegistry = make(map[string](*Handler))

//WrapHandlerFunc  instrument HTTP handler functions to collect HTTP metrics
func WrapHandlerFunc(name string, h tHttpHandlerFunc) tHttpHandlerFunc {
	proxy := newHandlerFunc(name, h)
	proxyHandlerFuncRegistry[name] = proxy
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

//WrapHandler  instrument HTTP handler object to collect HTTP metrics
func WrapHandler(name string, h http.Handler) http.Handler {
	proxy := newHandler(name, h)
	proxyHandlerRegistry[name] = proxy
	return proxy
}

//Print  print the metrics in each specified second
func Print(duration int) {
	timeDuration := time.Duration(duration)
	go func() {
		time.Sleep(timeDuration * time.Second)
		for {
			startTime := time.Now()
			for _, proxy := range proxyHandlerFuncRegistry {
				proxy.printMetrics(duration)
			}
			for _, proxy := range proxyHandlerRegistry {
				proxy.printMetrics(duration)
			}
			elapsedTime := time.Now().Sub(startTime)
			time.Sleep(timeDuration*time.Second - elapsedTime)
		}
	}()
}
