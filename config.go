package main

import (
	"encoding/json"
	"log"
	"net"
)

type Config struct {
	Sites   map[string]Site `json:"sites"`
	Routes  []Route         `json:"routes"`
	Default []string        `json:"default"`
}

type Site struct {
	URL     string    `json:"url"`
	Distros StringSet `json:"distros"`
}

func (conf *Config) FindURL(ip net.IP, distro string) string {
	sites := conf.FindSites(ip)
	for _, sitename := range sites {
		site, ok := conf.Sites[sitename]
		if !ok {
			log.Println("bad site in conf:", sitename)
			return ""
		}
		if _, ok = site.Distros[distro]; ok {
			return site.URL
		}
	}
	return ""
}

func (conf *Config) FindSites(ip net.IP) []string {
	// Quick path for invalid IP
	if ip == nil {
		return conf.Default
	}
	for _, r := range conf.Routes {
		if r.IPNet.Contains(ip) {
			return r.Sites
		}
	}
	return conf.Default
}

type Route struct {
	IPNet IPNet    `json:"ipnet"`
	Sites []string `json:"sites"`
}

// StringSet is a set of strings, represented as a list in JSON.
type StringSet map[string]struct{}

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
	*ss = make(map[string]struct{}, len(a))
	for _, s := range a {
		(*ss)[s] = struct{}{}
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
