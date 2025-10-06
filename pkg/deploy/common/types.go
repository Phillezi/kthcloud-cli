package common

import "net/http"

type RequestOption func(req *http.Request)
