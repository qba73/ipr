package main

type ipv4 struct {
	IPv4prefix string `json:"ipv4_prefix"`
	Region     string `json:"region"`
	Service    string `json:"service"`
}

type ipv6 struct {
	IPv6prefix string `json:"ipv6_prefix"`
	Region     string `json:"region"`
	Service    string `json:"service"`
}

type ipranges struct {
	SyncToken    string `json:"syncToken"`
	CreateDate   string `json:"createDate"`
	IPv4prefixes []ipv4 `json:"prefixes"`
	IPv6prefixes []ipv6 `json:"ipv6_prefixes"`
}
