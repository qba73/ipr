package ipr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type response struct {
	SyncToken  string `json:"syncToken"`
	CreateDate string `json:"createDate"`
	Prefixes   []struct {
		IPPrefix           string `json:"ip_prefix"`
		Region             string `json:"region"`
		Service            string `json:"service"`
		NetworkBorderGroup string `json:"network_border_group"`
	} `json:"prefixes"`
}

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

var myHttpClient = &http.Client{Timeout: 10 * time.Second}

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
	rx, err := getIPRanges(getEnv("AWS_IP_URL", "https://ip-ranges.amazonaws.com/ip-ranges.json"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rx)
}
