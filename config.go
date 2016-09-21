package main

import (
	"encoding/json"
	"log"
	"net"
)

type Config struct {
	Sites           map[string]Site `json:"sites"`
	Routes          []Route         `json:"routes"`
	DefaultOrdering []string        `json:"default-ordering"`
	Distros         StringSet       `json:"-"`
	Homepage        []byte          `json:"-"`
}

type Site struct {
	URL     string    `json:"url"`
	Distros StringSet `json:"distros"`
}

func (conf *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	err := json.Unmarshal(data, (*Alias)(conf))
	if err != nil {
		return err
	}
	// Populate conf.distros and conf.homepage.
	conf.Distros = make(StringSet)
	for _, site := range conf.Sites {
		for distro := range site.Distros {
			conf.Distros.Add(distro)
		}
	}
	conf.Homepage = conf.makeHomepage()
	return nil
}

// FindMirrorURL finds the most suitable mirror based on the combination of
// client IP and the distro required. If no mirror can be found for the distro,
// an empty string is returned.
func (conf *Config) FindMirrorURL(ip net.IP, distro string) string {
	if !conf.Distros.Has(distro) {
		return ""
	}
	sites := conf.FindOrdering(ip)
	for _, sitename := range sites {
		site, ok := conf.Sites[sitename]
		if !ok {
			log.Println("bad site in conf:", sitename)
			return ""
		}
		if site.Distros.Has(distro) {
			return site.URL
		}
	}
	return ""
}

func (conf *Config) FindOrdering(ip net.IP) []string {
	// Quick path for invalid IP
	if ip == nil {
		return conf.DefaultOrdering
	}
	for _, r := range conf.Routes {
		if r.IPNet.Contains(ip) {
			return r.Ordering
		}
	}
	return conf.DefaultOrdering
}

type Route struct {
	IPNet    IPNet    `json:"ipnet"`
	Ordering []string `json:"ordering"`
}

// StringSet is a set of strings, represented as a list in JSON.
type StringSet map[string]struct{}

func (ss StringSet) Add(s string) {
	ss[s] = struct{}{}
}

func (ss StringSet) Has(s string) bool {
	_, ok := ss[s]
	return ok
}

func (ss StringSet) MarshalJSON() ([]byte, error) {
	a := make([]string, 0, len(ss))
	for s := range ss {
		a = append(a, s)
	}
	return json.Marshal(a)
}

func (ss *StringSet) UnmarshalJSON(data []byte) error {
	var a []string
	err := json.Unmarshal(data, &a)
	if err != nil {
		return err
	}
	*ss = make(StringSet, len(a))
	for _, s := range a {
		ss.Add(s)
	}
	return nil
}

// IPNet adds MarshalJSON and UnmarshalJSON methods to net.IPNet.
type IPNet struct {
	net.IPNet
}

func (n *IPNet) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.IPNet.String())
}

func (n *IPNet) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	_, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return err
	}
	n.IPNet = *ipnet
	return nil
}
