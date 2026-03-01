package main

import (
	"net"
	"net/http"
	"time"

	"github.com/cjodo/go-cap"
)

func main() {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Second * 5,
			KeepAlive: time.Second * 30,
		}).DialContext,
		MaxIdleConns: 100,
		IdleConnTimeout: time.Second * 90,
		TLSHandshakeTimeout: time.Second * 10,
	}
	http := &http.Client{
		Transport: transport,
		Timeout: time.Second * 10,
	}
	opts := []redcap.Option{
		redcap.WithHTTPClient(http),
		redcap.WithMaxRetries(30),
	}

	redcap.NewClient("baseUrl", "token", opts...)
}
