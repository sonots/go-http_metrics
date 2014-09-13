package http_metrics

import (
	"fmt"
	metrics "github.com/yvasiyarov/go-metrics" // max,mean,min,stddev,percentile
	"net/http"
	"time"
)

// print infomation on each request
var Verbose = false

// Set Enable = false to turn off instrumentation
var Enable = true

// a set of proxies
var proxyRegistry = make(map[string](*proxyHandler))

type tHttpHandlerFunc func(http.ResponseWriter, *http.Request)
type proxyHandler struct {
	name                string
	originalHandler     http.Handler
	originalHandlerFunc tHttpHandlerFunc
	isFunc              bool
	timer               metrics.Timer
}

func newProxyHandlerFunc(name string, h tHttpHandlerFunc) *proxyHandler {
	return &proxyHandler{
		name:                name,
		originalHandlerFunc: h,
		isFunc:              true,
		timer:               metrics.NewTimer(),
	}
}
func newProxyHandler(name string, h http.Handler) *proxyHandler {
	return &proxyHandler{
		name:            name,
		originalHandler: h,
		isFunc:          false,
		timer:           metrics.NewTimer(),
	}
}

//print the elapsed time on each request if Verbose flag is true
func (proxy *proxyHandler) printVerbose(r *http.Request, elapsedTime time.Duration) {
	fmt.Printf("time:%v\thandler:%s\tmethod:%s\tpath:%s\tparams:%s\telapsed:%f\n",
		time.Now(),
		proxy.name,
		r.Method,
		r.URL.Path,
		r.URL.Query().Encode(),
		elapsedTime.Seconds(),
	)
}

// measure elapsed time
func (proxy *proxyHandler) measure(startTime time.Time, r *http.Request) {
	elapsedTime := time.Now().Sub(startTime)
	proxy.timer.Update(elapsedTime)
	if Enable && Verbose {
		proxy.printVerbose(r, elapsedTime)
	}
}

///// instrument functions

// instrument handler
func (proxy *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var startTime time.Time
	if Enable {
		startTime = time.Now()
	}
	if proxy.isFunc {
		proxy.originalHandlerFunc(w, r)
	} else {
		proxy.originalHandler.ServeHTTP(w, r)
	}
	if Enable {
		defer proxy.measure(startTime, r)
	}
}

///// package functions

//WrapHandlerFunc  instrument HTTP handler functions to collect HTTP metrics
func WrapHandlerFunc(name string, h tHttpHandlerFunc) tHttpHandlerFunc {
	proxy := newProxyHandlerFunc(name, h)
	proxyRegistry[name] = proxy
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

//WrapHandler  instrument HTTP handler object to collect HTTP metrics
func WrapHandler(name string, h http.Handler) http.Handler {
	proxy := newProxyHandler(name, h)
	proxyRegistry[name] = proxy
	return proxy
}

//Print  print the metrics in each specified second
func Print(duration int) {
	timeDuration := time.Duration(duration)
	go func() {
		time.Sleep(timeDuration * time.Second)
		for {
			startTime := time.Now()
			for name, proxy := range proxyRegistry {
				timer := proxy.timer
				count := timer.Count()
				if count > 0 {
					fmt.Printf("time:%v\thandler:%s\tcount:%d\tmax:%f\tmean:%f\tmin:%f\tpercentile95:%f\tduration:%d\n",
						time.Now(),
						name,
						timer.Count(),
						float64(timer.Max())/float64(time.Second),
						timer.Mean()/float64(time.Second),
						float64(timer.Min())/float64(time.Second),
						timer.Percentile(0.95)/float64(time.Second),
						duration,
					)
					proxy.timer = metrics.NewTimer()
				}
			}
			elapsedTime := time.Now().Sub(startTime)
			time.Sleep(timeDuration*time.Second - elapsedTime)
		}
	}()
}
