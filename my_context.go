package main

import (
	"net/http"

	gae "appengine"
	"appengine/urlfetch"
	"golang.org/x/net/context"
	stdgae "google.golang.org/appengine"
)

// MyContext holds two contexts, seems like eventually the net/context will
// be the standard one??
// https://groups.google.com/forum/#!searchin/google-appengine-go/golang.org$2Fx$2Fnet$2Fcontext
type MyContext struct {
	Env string
	gae.Context
	StdContext context.Context
	*http.Client
	W http.ResponseWriter
	R *http.Request
}

func Wrap(wrapee func(*MyContext), mc *MyContext) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if mc.Env != "test" {
			mc.Context = gae.NewContext(r)
			mc.StdContext = stdgae.NewContext(r)
			mc.Client = urlfetch.Client(mc.Context) // TODO: remove this, only create contexts here
		}

		mc.W = w
		mc.R = r
		wrapee(mc)
	}
}
