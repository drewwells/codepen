// +build appengine

package main

import (
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

func init() {
	r := Register()
	http.Handle("/",
		withCORS(r),
	)
}

// CollectionHandler parses collection URLs on Codepen
func CollectionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	collection(client, w, r)
}

// CheckHandler is responsible for refreshing the list of most
// recent browsers.
func CheckHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	pt := "http://" + CodePenURL + r.URL.Path
	resp, err := client.Get(pt)

	if err != nil {
		w.Write([]byte(err.Error()))
	}
	check(w, resp, pt)
}
