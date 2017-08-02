# graphite-web-proxy
Small proxy service to allow running graphite-web locally with a GrafanaCloud hostedMetrics service as a backend


## Usage
```
$ graphite-web-proxy -h
Usage of graphite-web-proxy:
  -addr string
    	host:port to listen on. (default "0.0.0.0:8181")
  -alsologtostderr
    	log to standard error as well as files
  -api-key string
    	grafana.com api key (default "xxxxxx")
  -log_backtrace_at value
    	when logging hits line file:N, emit a stack trace
  -log_dir string
    	If non-empty, write log files in this directory
  -logtostderr
    	log to standard error instead of files
  -stderrthreshold value
    	logs at or above this threshold go to stderr
  -tsdb-url string
    	gateway address of hosted-metrics service. (default "https://tsdb-x-foo.hosted-metrics.grafana.net")
  -v value
    	log level for V logs
  -vmodule value
    	comma-separated list of pattern=N settings for file-filtered logging
```

## Getting started

- download the proxy. 
```
go install github.com/raintank/graphite-web-proxy
```
- run the proxy
```
graphite-web-proxy -logtostderr -tsdb-url https://<your HostedMetrics url> -api-key <grafana.com API KEY>
```

- update your graphite-web installation to use the proxy as a *CLUSTER_SERVERS* (this is defined in graphite local_settings.py file)
If you dont have graphite-web running, you can use docker
```
docker run -p 8080:80 -e GRAPHITE_CLUSTER_SERVERS=<localIP>:8181 raintank/graphite-mt
```
