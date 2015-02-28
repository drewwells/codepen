// Currentbrowsers package attempts to find the most recent
// versions of popular browsers.  This data is then easily
// consumable as an API.
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	CodePenURL = "http://codepen.io/"
)

// Browser contains the necessary information for browser
// type and release version.
type Browser struct {
	//Chrome Desktop, Chrome Android, Chrome iOS
	Type    string `json:"type"`
	Version string `json:"version"`
}

func Register() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/{user}/details/{pen}", CheckHandler)
	r.HandleFunc("/collection/{id}", CollectionHandler)
	// Allow specifying individual pages
	//r.HandleFunc("/collection/{id}/{page}", CollectionHandler)
	return r
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods",
			"POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		h.ServeHTTP(w, r)
	})
}

// IndexHandler is responsible for listing the most
// recent browsers.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Access URL /{user}/details/{pen}"))
}

func WriteJSON(w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func collection(client *http.Client, w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	page := mux.Vars(r)["page"]

	var p int
	if page != "" {
		p, _ = strconv.Atoi(page)
	}
	w.Write([]byte("fetch"))
	FetchCollection(w, client, id, p)

}

// FetchCollection attempts acquiring a collection, returning error
// if empty
func FetchCollection(w http.ResponseWriter, c *http.Client, id string, page int) {
	urlTmpl := "http://codepen.io/collection/next_grid_for_collection/%s?page=1"
	_ = urlTmpl
	// Iterate through all collections
	if page == 0 {
		FetchCollection(w, c, id, 1)
		return
	}

}

func check(w http.ResponseWriter, resp *http.Response, path string) {
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	re := regexp.MustCompile(`<li class="module`)
	li := regexp.MustCompile(`</li>`)
	num := regexp.MustCompile("\\d+")
	mats := re.FindAllIndex(bs, -1)

	cts := make([]int64, len(mats))
	var hits string
	_ = hits
	for i, m := range mats {
		pos := li.FindIndex(bs[m[0]:])
		hits += string(bs[m[0] : m[0]+pos[0]])

		n := string(num.Find(bs[m[0] : m[0]+pos[0]]))
		cts[i], _ = strconv.ParseInt(n, 10, 32)
	}

	if len(mats) == 0 {
		WriteJSON(w, map[string]interface{}{
			"message": "No details found",
			"url":     path,
			"src":     string(bs),
		})
		return
	}

	WriteJSON(w, map[string]interface{}{
		"views":    cts[0],
		"comments": cts[1],
		"hearts":   cts[2],
		"referrer": path,
	})

}
