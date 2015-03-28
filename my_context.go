package main

import (
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

type MyContext struct {
	Env string
	appengine.Context
	*http.Client
	W http.ResponseWriter
	R *http.Request
}

func Wrap(wrapee func(*MyContext), mc *MyContext) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if mc.Env != "test" {
			mc.Context = appengine.NewContext(r)
			mc.Client = urlfetch.Client(mc.Context)
		}

		mc.W = w
		mc.R = r
		wrapee(mc)
	}
}
