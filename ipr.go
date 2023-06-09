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

// Response represents data received from AWS after
// successfull call to the whitelisted ip ranges url:
// "https://ip-ranges.amazonaws.com/ip-ranges.json"
type Response struct {
	SyncToken  string `json:"syncToken"`
	CreateDate string `json:"createDate"`
	Prefixes   []struct {
		IPPrefix           string `json:"ip_prefix"`
		Ipv6Prefix         string `json:"ipv6_prefix"`
		Region             string `json:"region"`
		Service            string `json:"service"`
		NetworkBorderGroup string `json:"network_border_group"`
	} `json:"prefixes"`
}

// IPRange represents information about whitelisted IP range.
type IPRange struct {
	Type               string `json:"ipv"`
	IPprefix           string `json:"ip_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

// IPRanges represents IPv4 and IPv6 AWS whitelisted IP ranges.
type IPRanges struct {
	SyncToken  int       `json:"sync_token"`
	CreateDate time.Time `json:"create_date"`
	IPv4Ranges []IPRange `json:"ipv4_ranges"`
	IPv6Ranges []IPRange `json:"ipv6_ranges"`
}

// ToCSVRecords creates csv records ready to pass to a csv writer.
func (ip IPRanges) CSVRecords() [][]string {
	rx := [][]string{{"ip_type", "ip_prefix", "region", "service", "network_border_group"}}
	for _, i := range ip.IPv4Ranges {
		row := []string{i.Type, i.IPprefix, i.Region, i.Service, i.NetworkBorderGroup}
		rx = append(rx, row)
	}
	for _, i := range ip.IPv6Ranges {
		row := []string{i.Type, i.IPprefix, i.Region, i.Service, i.NetworkBorderGroup}
		rx = append(rx, row)
	}
	return rx
}

// Client holds data for making calls
// to the AWS whitelisted endpoint.
type Client struct {
	URL        string
	HTTPClient *http.Client
}

// NewClient creates a default API client ready to talk to AWS endpoint.
//
// Default client uses the official AWS endpoint to fetch ip ranges.
// You can specify different enpoint by exporting the `AWS_IP_URL` env var.
func NewClient() *Client {
	return &Client{
		URL: getEnv("AWS_IP_URL", "https://ip-ranges.amazonaws.com/ip-ranges.json"),
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Ranges returns whitelisted AWS IPv4 and IPv6 ranges.
func (c *Client) GetRanges() (IPRanges, error) {
	res, err := c.HTTPClient.Get(c.URL)
	if err != nil {
		return IPRanges{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return IPRanges{}, fmt.Errorf("ipr: got status code %v", res.StatusCode)
	}
	var resp Response
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return IPRanges{}, fmt.Errorf("ipr: decoding response body %w", err)
	}
	return ParseResponse(resp)
}

// ParseResponse takes response struct and returns IPRanges.
//
// It errors if the token or timestamp is malformed, or when
// either IPv4 or IPv6 cannot be parsed.
func ParseResponse(resp Response) (IPRanges, error) {
	token, err := strconv.Atoi(resp.SyncToken)
	if err != nil {
		return IPRanges{}, fmt.Errorf("ipr: malformed sync token: %v, %w", resp.SyncToken, err)
	}
	createDate, err := time.Parse("2006-01-02-15-04-05", resp.CreateDate)
	if err != nil {
		return IPRanges{}, err
	}

	var (
		ip4ranges, ip6ranges []IPRange
		prefix, iptype       string
	)

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

// GetIPRanges pulls AWS Whitelisted IP ranges.
//
// It uses default IP Range client. If you want
// to use a different AWS Whitelisted IP URL
// export the env variable "AWS_IP_URL".
func GetIPRanges() (IPRanges, error) {
	return NewClient().GetRanges()
}

func runCLI(r io.Reader, w, er io.Writer) int {
	c := NewClient()
	rx, err := c.GetRanges()
	if err != nil {
		fmt.Fprintln(er, err)
		os.Exit(1)
	}
	fmt.Fprintln(w, rx)
	return 0
}

func Main() int {
	return runCLI(os.Stdin, os.Stdout, os.Stderr)
}
