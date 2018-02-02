package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"

	"github.com/golang/glog"
)

var (
	TSDB_URL string
	APIKEY   string
	ADDR     string
	ORGID    int
)

func init() {
	flag.StringVar(&TSDB_URL, "tsdb-url", "https://tsdb-x-foo.hosted-metrics.grafana.net", "gateway address of hosted-metrics service.")
	flag.StringVar(&APIKEY, "api-key", "xxxxxx", "grafana.com api key")
	flag.StringVar(&ADDR, "addr", "0.0.0.0:8181", "host:port to listen on.")
	flag.IntVar(&ORGID, "org", 0, "Include X-Tsdb-Org header, required when using an admin api-key")
}

func main() {
	flag.Parse()
	glog.CopyStandardLogTo("INFO")
	glog.Infof("graphite-web proxy starting up on %s", ADDR)

	u, err := url.Parse(TSDB_URL)
	if err != nil {
		glog.Fatalf("failed to parse tsdburl. %s", err)
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			glog.V(5).Info("director rewriting request URL")
			req.URL.Scheme = u.Scheme
			req.URL.Host = u.Host
			req.Host = u.Host
			req.URL.Path = path.Join("/graphite", req.URL.Path)
			glog.V(5).Infof("path rewriten to %s", req.URL.Path)
			if ORGID != 0 {
				req.Header.Set("X-Tsdb-Org", fmt.Sprintf("%d", ORGID))
			}
			req.ParseMultipartForm(1024 * 1024 * 10)
			data := req.Form

			if data.Get("local") == "1" {
				glog.V(4).Info("setting local=0 in query params")
				data.Set("local", "0")
			}

			body := data.Encode()
			switch method := req.Method; method {
			case "POST":
				glog.V(5).Infof("rewriting request body to %s", body)
				req.Body = ioutil.NopCloser(bytes.NewReader([]byte(body)))
				req.Header.Set("content-type", "application/x-www-form-urlencoded")
			case "PUT":
				glog.V(5).Infof("rewriting request body to %s", body)
				req.Body = ioutil.NopCloser(bytes.NewReader([]byte(body)))
				req.Header.Set("content-type", "application/x-www-form-urlencoded")

			case "GET":
				glog.V(5).Infof("rewriting query params to %s", body)
				req.URL.RawQuery = body
			}

			req.ContentLength = int64(len(body))
			// add auth headers
			req.Header.Set("AUTHORIZATION", fmt.Sprintf("Bearer %s", APIKEY))

			glog.Infof("%s %+v", req.Method, req.URL)
		},
	}
	glog.Fatal(http.ListenAndServe(ADDR, proxy))
}
