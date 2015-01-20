// Currentbrowsers package attempts to find the most recent
// versions of popular browsers.  This data is then easily
// consumable as an API.
package currentbrowsers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"

	"appengine"
	"appengine/urlfetch"
)

const (
	url = "http://codepen.io/"
)

// Browser contains the necessary information for browser
// type and release version.
type Browser struct {
	//Chrome Desktop, Chrome Android, Chrome iOS
	Type    string `json:"type"`
	Version string `json:"version"`
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/{user}/details/{pen}", CheckHandler)
	http.Handle("/", withCORS(r))
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

// CheckHandler is responsible for refreshing the list of most
// recent browsers.
func CheckHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	pt := url + r.URL.Path
	resp, err := client.Get(pt)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	re := regexp.MustCompile(`<li class="module`)
	li := regexp.MustCompile(`</li>`)
	num := regexp.MustCompile("\\d+")
	mats := re.FindAllIndex(bs, -1)
	log.Println("request for stuff", mux.Vars(r))
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
			"url":     pt,
			"src":     string(bs),
		})
		return
	}

	WriteJSON(w, map[string]interface{}{
		"views":    cts[0],
		"comments": cts[1],
		"hearts":   cts[2],
		"referrer": pt,
	})

}
