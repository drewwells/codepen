// Currentbrowsers package attempts to find the most recent
// versions of popular browsers.  This data is then easily
// consumable as an API.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/net/html"
)

const (
	CodePenURL = "codepen.io"
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
		// Apply JSON headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Requested-With", "XMLHttpRequest")

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

	FetchCollection(client, w, id, p)

}

// FetchCollection attempts acquiring a collection, returning error
// if empty
func FetchCollection(client *http.Client, w http.ResponseWriter, id string, page int) {

	// Iterate through all collections
	if page == 0 {
		FetchCollection(client, w, id, 1)
		return
	}
	//scheme: //[userinfo@]host/path[?query][#fragment]

	path := url.URL{
		Scheme:   "http",
		Host:     CodePenURL,
		Path:     "collection/next_grid_for_collection/" + id,
		RawQuery: "page=%d" + strconv.Itoa(page),
	}

	resp, err := client.Get(path.String())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	bs, err := parseCollection(resp.Body)
	// There's no more stuffs
	if err == io.EOF {
		w.Write(bs)
	} else {
		w.Write([]byte("Theres more!"))
	}
}

// Collection defines JSON return from /collection/{id}
type Collection struct {
	HTML      string `json:"html"`
	Success   bool   `json:"success"`
	PageType  string `json:"pageType"`
	UpdateURL bool   `json:"updateURL"`
}

// Response format
type Response map[string][]int

func NewResponse() Response {
	return make(map[string][]int)
}

func (r Response) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("[")
	virgin := true
	for u := range r {
		if !virgin {
			buf.WriteString(",")
		}
		virgin = false
		buf.WriteString("{")
		frags := strings.Split(u, "/")
		id := frags[len(frags)-1]
		s := fmt.Sprintf(
			`"pen":"%s", "url": "%s", "comments":%d, `+
				`"views":%d, "loves": %d`,
			u, id, r[u][0], r[u][1], r[u][2],
		)
		buf.WriteString(s)
		buf.WriteString("}")
	}
	buf.WriteString("]")
	return buf.Bytes(), nil
}

func parseCollection(r io.Reader) ([]byte, error) {

	dec := json.NewDecoder(r)
	col := Collection{}
	if err := dec.Decode(&col); err != nil {
		return nil, err
	}
	if len(col.HTML) == 0 {
		return nil, errors.New("No HTML found")
	}

	buf := bytes.NewBufferString(col.HTML)

	node, err := html.Parse(buf)
	if err != nil {
		return nil, err
	}

	m := NewResponse()
	walker(node, m)
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return bs, io.EOF
}

func parseAttributes(node *html.Node) (string, bool) {
	var hit bool
	var id string
	for _, attr := range node.Attr {
		switch attr.Key {
		case "class":
			if strings.Contains(attr.Val, "single-stat") {
				hit = true
			}
		case "href":
			id = attr.Val
		case "data-hashid":
			// Loves have another id, find the pen id
			return parseAttributes(node.PrevSibling.PrevSibling)
		}
	}
	return id, hit
}

func walker(node *html.Node, m map[string][]int) {
	dig := true

	id, hit := parseAttributes(node)

	if hit {
		if _, ok := m[id]; !ok {
			m[id] = []int{}
		}

		spanText := node.FirstChild
		if spanText != nil {
			count := strings.TrimSpace(spanText.Data)
			if len(count) == 0 {
				count = strings.TrimSpace(spanText.NextSibling.FirstChild.Data)
			}

			ct, err := strconv.Atoi(count)
			if err != nil {
				log.Fatal(err)
			}
			m[id] = append(m[id], ct)

		}
		dig = false
	}

	// Breadth first
	if next := node.NextSibling; next != nil {
		walker(next, m)
	}
	if child := node.FirstChild; dig && child != nil {
		walker(child, m)
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
