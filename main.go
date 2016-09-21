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
	flagAddr = flag.String("addr", ":8080", "address to listen")
	flagConf = flag.String("conf", "", "Configuration file")
)

var config Config

func getIPOverride(h http.Header) string {
	ff := h.Get("X-Mirrorhub-IP-Override")
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
	path := r.URL.Path
	if path == "" || path == "/" {
		w.Write(config.Homepage)
		return
	}

	if path[0] != '/' {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Find distro part (the first level in the path)
	i := strings.IndexByte(path[1:], '/')
	var distro string
	if i == -1 {
		distro = path[1:]
	} else {
		distro = path[1 : i+1]
	}

	// r.RemoteAddr is IP:port; stripe the port to get IP of the requester
	remote := r.RemoteAddr
	reqip := remote[:strings.LastIndexByte(remote, ':')]

	override := getIPOverride(r.Header)

	var ip net.IP
	if override == "" {
		ip = net.ParseIP(reqip)
	} else {
		ip = net.ParseIP(reqip)
	}
	mirror := config.FindMirrorURL(ip, distro)

	if override == "" {
		log.Printf("%s %s %q -> %s", reqip, r.Method, r.URL, mirror)
	} else {
		log.Printf("%s (%s) %s %q -> %s", override, reqip, r.Method, r.URL, mirror)
	}

	if mirror != "" {
		// NOTE: We drop the query string of the original URL.
		w.Header().Add("Location", mirror+r.URL.EscapedPath())
		w.WriteHeader(http.StatusFound)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func loadConfig(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	return decoder.Decode(&config)
}

func main() {
	flag.Parse()
	if *flagConf == "" {
		log.Fatal("-conf required")
	}
	err := loadConfig(*flagConf)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("loaded config from", *flagConf)
	}

	http.HandleFunc("/", handler)
	log.Println("going to listen", *flagAddr)
	log.Fatal(http.ListenAndServe(*flagAddr, nil))
}
