package main

import (
	"net/http"

	"appengine"
)

type MyContext struct {
	appengine.Context
	*http.Client
	W http.ResponseWriter
	R *http.Request
}

func Wrap(wrapee func(*MyContext), mc *MyContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mc.W = w
		mc.R = r
		wrapee(mc)
	}
}
