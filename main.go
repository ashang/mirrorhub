package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	addr  = flag.String("addr", ":8080", "address to listen")
	route = flag.String("route", "", "JSON file to load routing table")
)

var routing RoutingTable

func forwardedForIP(h http.Header) string {
	ff := h.Get("X-Forwarded-For")
	if ff == "" {
		return ""
	}
	comma := strings.IndexByte(ff, ',')
	if comma == -1 {
		return ff
	}
	return ff[:comma]
}

func handler(w http.ResponseWriter, r *http.Request) {
	remote := r.RemoteAddr
	realip := remote[:strings.LastIndexByte(remote, ':')]

	ipstr := forwardedForIP(r.Header)
	if ipstr == "" {
		ipstr = realip
	}
	mirror := routing.Route(net.ParseIP(ipstr))

	if ipstr != realip {
		log.Printf("%s (%s) %s %q -> %s", ipstr, realip, r.Method, r.URL, mirror)
	} else {
		log.Printf("%s %s %q -> %s", ipstr, r.Method, r.URL, mirror)
	}

	// NOTE: We drop the query string of the original URL.
	w.Header().Add("Location", mirror+r.URL.EscapedPath())
	w.WriteHeader(http.StatusFound)
}

func loadRoutingTable(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	var tab RoutingTableWithURLMap
	err = decoder.Decode(&tab)
	if err != nil {
		return err
	}
	routing = *tab.ResolveURLMap()
	return nil
}

func main() {
	flag.Parse()
	if *route == "" {
		log.Fatal("-route required")
	}
	err := loadRoutingTable(*route)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
