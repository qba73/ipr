package ipr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const awsurl = "https://ip-ranges.amazonaws.com/ip-ranges.json"

var myHttpClient = &http.Client{Timeout: 10 * time.Second}

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

// getData fetches the given url and decodes into target (interface)
func getData(url string, target interface{}) error {
	res, err := myHttpClient.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(target)
}

func getIPRanges(url string) (ipranges, error) {
	var rx ipranges
	if err := getData(url, &rx); err != nil {
		return ipranges{}, err
	}
	return rx, nil
}

func RunCLI() {
	fmt.Println("fetching whitelisted AWS IP Ranges")
	rx, err := getIPRanges(awsurl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rx)
}
