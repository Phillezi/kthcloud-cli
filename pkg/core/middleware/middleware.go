package middleware

import (
	"net/http"
	"net/url"
)

type Middleware func(req *http.Request)

func WithDefaultURL(url *url.URL) Middleware {
	return func(req *http.Request) {
		req.URL = url
	}
}
