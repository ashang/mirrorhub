package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "address to listen")

func findMirror(remote string) string {
	// TODO Implement something slightly more useful.
	return "https://mirrors.tuna.tsinghua.edu.cn"
}

func handler(w http.ResponseWriter, r *http.Request) {
	mirror := findMirror(r.RemoteAddr)

	log.Printf("%s %s %q -> %s", r.RemoteAddr, r.Method, r.URL, mirror)

	// NOTE: We drop the query string of the original URL.
	w.Header().Add("Location", mirror+r.URL.EscapedPath())
	w.WriteHeader(http.StatusFound)
}

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)
	http.ListenAndServe(*addr, nil)
}
