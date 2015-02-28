// +build !appengine

package main

import "net/http"

func main() {
	http.ListenAndServe("localhost:8080", withCORS(Register()))
}

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	path := "http://" + CodePenURL + r.URL.Path
	resp, err := http.Get(path)

	if err != nil {
		w.Write([]byte(err.Error()))
	}
	check(w, resp, path)
}

// CollectionHandler captures /collection/{id}/{?page}
func CollectionHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	collection(client, w, r)

}
