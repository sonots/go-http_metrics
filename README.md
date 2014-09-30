# go-http\_metrics

Instrument http request.

go-http\_metrics wraps http.Handle or HandleFunc to instrument each http request. 

# Usage

```go
import (
  "net/http"
  "github.com/sonots/go-http_metrics"
)

func main() {
  // http.HandleFunc("/", rootHandleFunc)
  http.HandleFunc("/", http_metrics.WrapHandlerFunc("rootHandleFunc", rootHandleFunc))
  // http.Handle("/static/", staticHandle)
  http.Handle("/static/", http_metrics.WrapHandler("staticHandle", staticHandle))

  http_metrics.Print(1) // print metrics on each 1 second
  // http_metrics.Verbose = true // if you want to print on each request
  // http_metrics.Enable = false // if you want to turn off instrumentation
  http.ListenAndServe("0.0.0.0:5050", nil)
}
```

Output Example (LTSV format):

```
time:2014-09-08 03:27:50.346193673 +0900 JST  handler:rootHandleFunc count:1 max:0.001626    mean:0.001626   min:0.001626    percentile95:0.001626    duration:1
```

Verbose Output Example (LTSV format):

```
time:2014-09-08 03:27:50.346193673 +0900 JST  handler:rootHandleFunc method:GET      path:/    params:foo=bar   elapsed:0.001626
```

# API

## Print

Print summarized metrics on each specified second:

```go
http_metrics.Print(60) # 60 seconds
```

## Verbose

Print metrics on each request:

```go
http_metrics.Verbose = true
```

## Enable

Diable instrumentation:

```go
http_metrics.Enable = false
```

## OnResponse

Set an arbitrary callback called on each response as:

```go
http_metrics.OnResponse = func() {
  sql_metrics.Flush()
  template_metrics.Flush()
}
```

This example flushes [sql_metrics](https://github.com/sonots/go-sql_metrics) and [template_metrics](https://github.com/sonots/go-template_metrics) on each request.

# ToDo

* write tests

# Contribution

* Fork (https://github.com/sonots/go-http_metrics/fork)
* Create a feature branch
* Commit your changes
* Rebase your local changes against the master branch
* Create new Pull Request

# Copyright

* See [LICENSE](./LICENSE)
