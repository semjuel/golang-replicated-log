package model

import "net/http"

type Routes []Route

type Route struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
}
