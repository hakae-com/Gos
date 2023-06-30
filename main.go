package main

import (
	"fmt"
	"gos/gos"
	"net/http"
)

func main() {
	r := gos.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})
	r.GET("/register", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})

	err := r.RUN(":9999")
	if err != nil {
		return
	}
}
