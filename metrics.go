package http_metrics

import (
	"fmt"
	metrics "github.com/sonots/go-metrics" // max,mean,min,stddev,percentile
	"net/http"
	"time"
)

type Metrics struct {
	name  string
	timer metrics.Timer
}

func newMetrics(name string) *Metrics {
	return &Metrics{
		name:  name,
		timer: metrics.NewTimer(),
	}
}

// measure elapsed time
func (proxy *Metrics) measure(startTime time.Time, r *http.Request) {
	elapsedTime := time.Now().Sub(startTime)
	proxy.timer.Update(elapsedTime)
	if Enable && Verbose {
		proxy.printVerbose(r, elapsedTime)
	}
}

//print the elapsed time on each request if Verbose flag is true
func (proxy *Metrics) printVerbose(r *http.Request, elapsedTime time.Duration) {
	fmt.Printf(
		"time:%v\thandler:%s\tmethod:%s\tpath:%s\tparams:%s\telapsed:%f\n",
		time.Now(),
		proxy.name,
		r.Method,
		r.URL.Path,
		r.URL.Query().Encode(),
		elapsedTime.Seconds(),
	)
}

func (proxy *Metrics) printMetrics(duration int) {
	timer := proxy.timer
	count := timer.Count()
	if count > 0 {
		fmt.Printf(
			"time:%v\thandler:%s\tcount:%d\tmax:%f\tmean:%f\tmin:%f\tpercentile95:%f\tsum:%f\tduration:%d\n",
			time.Now(),
			proxy.name,
			timer.Count(),
			float64(timer.Max())/float64(time.Second),
			timer.Mean()/float64(time.Second),
			float64(timer.Min())/float64(time.Second),
			timer.Percentile(0.95)/float64(time.Second),
			float64(timer.Sum())/float64(time.Second),
			duration,
		)
		proxy.timer = metrics.NewTimer()
	}
}
