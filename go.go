// +build !appengine

package main

import "net/http"

func main() {
	http.ListenAndServe("localhost:8080", Register())
}

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	path := CodePenURL + r.URL.Path
	resp, err := http.Get(path)

	if err != nil {
		w.Write([]byte(err.Error()))
	}
	check(w, resp, path)
}

func CollectionHandler(w http.ResponseWriter, r *http.Request) {

}
