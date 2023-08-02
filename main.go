package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const musicDir = "music"
	const port = 8080

	http.Handle("/", fileServerHandler(musicDir))
	log.Printf("Serving %s on HTTP port: %v\n", musicDir, port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func fileServerHandler(dir string) http.HandlerFunc {
	musicDir := http.Dir(dir)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.FileServer(musicDir).ServeHTTP(w, r)
	}
}
