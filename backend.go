package main

import (
	"net/http/httputil"
	"net/url"
)

type backend struct {
	url     *url.URL
	proxy   *httputil.ReverseProxy
	isAlive bool
}
