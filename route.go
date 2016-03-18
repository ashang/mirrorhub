package main

import (
	"encoding/json"
	"net"
)

type RoutingTable struct {
	Routes  []Route `json:"routes"`
	Default string  `json:"default"`
}

func (tab *RoutingTable) Route(ip net.IP) string {
	// Quick path for invalid IP
	if ip == nil {
		return tab.Default
	}
	for _, r := range tab.Routes {
		if r.IPNet.Contains(ip) {
			return r.URL
		}
	}
	return tab.Default
}

type RoutingTableWithURLMap struct {
	RoutingTable
	URLMap map[string]string `json:"url-map"`
}

func lookup(m map[string]string, k string) string {
	if v, ok := m[k]; ok {
		return v
	}
	return k
}

func (tabm *RoutingTableWithURLMap) ResolveURLMap() *RoutingTable {
	tab := tabm.RoutingTable
	for i := range tab.Routes {
		tab.Routes[i].URL = lookup(tabm.URLMap, tab.Routes[i].URL)
	}
	tab.Default = lookup(tabm.URLMap, tab.Default)
	return &tab
}

type Route struct {
	IPNet IPNet  `json:"ipnet"`
	URL   string `json:"url"`
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
