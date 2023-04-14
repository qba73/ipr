package ipr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"os"
	"strconv"
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
		IPPrefix           string `json:"ip_prefix,omitempty"`
		Ipv6Prefix         string `json:"ipv6_prefix,omitempty"`
		Region             string `json:"region"`
		Service            string `json:"service"`
		NetworkBorderGroup string `json:"network_border_group"`
	} `json:"prefixes"`
}

type IPRange struct {
	Type               string `json:"ipv"`
	IPprefix           string `json:"ip_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

type IPRanges struct {
	SyncToken  int       `json:"sync_token"`
	CreateDate time.Time `json:"create_date"`
	IPv4Ranges []IPRange `json:"ipv4_ranges"`
	IPv6Ranges []IPRange `json:"ipv6_ranges"`
}

type Client struct {
	URL        string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		URL: getEnv("AWS_IP_URL", "https://ip-ranges.amazonaws.com/ip-ranges.json"),
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Ranges() (IPRanges, error) {
	res, err := c.HTTPClient.Get(c.URL)
	if err != nil {
		return IPRanges{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return IPRanges{}, fmt.Errorf("ipr: got response code %v", res.StatusCode)
	}
	var resp response
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return IPRanges{}, fmt.Errorf("ipr: decoding response body %w", err)
	}
	return ProcessRanges(resp)
}

func ProcessRanges(resp response) (IPRanges, error) {
	token, err := strconv.Atoi(resp.SyncToken)
	if err != nil {
		return IPRanges{}, fmt.Errorf("ipr: malformed sync token: %v, %w", resp.SyncToken, err)
	}
	createDate, err := time.Parse("2006-01-02-15-04-05", resp.CreateDate)
	if err != nil {
		return IPRanges{}, err
	}

	var ip4ranges []IPRange
	var ip6ranges []IPRange
	var prefix string
	var iptype string

	for _, p := range resp.Prefixes {
		if p.IPPrefix != "" {
			prefix = p.IPPrefix
			iptype = "ipv4"
		}
		if p.Ipv6Prefix != "" {
			prefix = p.Ipv6Prefix
			iptype = "ipv6"
		}

		pr, err := netip.ParsePrefix(prefix)
		if err != nil {
			return IPRanges{}, err
		}
		if pr.Addr().Is4() {
			ipv := IPRange{
				Type:               iptype,
				IPprefix:           prefix,
				Region:             p.Region,
				Service:            p.Service,
				NetworkBorderGroup: p.NetworkBorderGroup,
			}
			ip4ranges = append(ip4ranges, ipv)
			continue
		}

		if pr.Addr().Is6() {
			ipv := IPRange{
				Type:               iptype,
				IPprefix:           prefix,
				Region:             p.Region,
				Service:            p.Service,
				NetworkBorderGroup: p.NetworkBorderGroup,
			}
			ip6ranges = append(ip6ranges, ipv)
			continue
		}
	}

	ipx := IPRanges{
		SyncToken:  token,
		CreateDate: createDate.UTC(),
		IPv4Ranges: ip4ranges,
		IPv6Ranges: ip6ranges,
	}
	return ipx, nil
}

func GetIPRanges() (IPRanges, error) {
	return NewClient().Ranges()
}

func runCLI(r io.Reader, w, er io.Writer) int {
	return 0
}

func Main() int {
	return runCLI(os.Stdin, os.Stdout, os.Stderr)
}
