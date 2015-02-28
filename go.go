// +build !appengine

package main

import "net/http"

func main() {
	http.ListenAndServe("localhost:8080", Register())
}

func CheckHandler(w http.ResponseWriter, r *http.Request) {

}

func CollectionHandler(w http.ResponseWriter, r *http.Request) {

}
