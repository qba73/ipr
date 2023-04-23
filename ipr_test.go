package ipr_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/ipr"
)

func newTestTLSServer(resp io.Reader, t *testing.T) *httptest.Server {
	t.Helper()

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.Copy(w, resp)
		if err != nil {
			t.Fatal(err)
		}
	}))
	return ts
}

func TestGetRanges_SucceedsOnValidInput(t *testing.T) {
	t.Parallel()

	ts := newTestTLSServer(validResponse, t)
	defer ts.Close()

	c := ipr.NewClient()
	c.HTTPClient = ts.Client()
	c.URL = ts.URL

	got, err := c.GetRanges()
	if err != nil {
		t.Fatal(err)
	}
	want := ipr.IPRanges{
		SyncToken:  1676592786,
		CreateDate: time.Date(2023, 02, 17, 00, 13, 06, 00, time.UTC),
		IPv4Ranges: []ipr.IPRange{
			{
				Type:               "ipv4",
				IPprefix:           "13.34.65.64/27",
				Region:             "il-central-1",
				Service:            "AMAZON",
				NetworkBorderGroup: "il-central-1",
			},
		},
		IPv6Ranges: []ipr.IPRange{
			{
				Type:               "ipv6",
				IPprefix:           "2600:1ff8:e000::/40",
				Region:             "sa-east-1",
				Service:            "S3",
				NetworkBorderGroup: "sa-east-1",
			},
		},
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetRanges_ErrorsOnBogusToken(t *testing.T) {
	t.Parallel()

	ts := newTestTLSServer(invalidResponseWithBogusSyncToken, t)
	defer ts.Close()

	c := ipr.NewClient()
	c.HTTPClient = ts.Client()
	c.URL = ts.URL

	_, err := c.GetRanges()
	if err == nil {
		t.Error("want err, got nil")
	}
}

func TestGetRanges_ErrorsOnMalformedDate(t *testing.T) {
	t.Parallel()

	ts := newTestTLSServer(invalidResponseWithMalformedDate, t)
	defer ts.Close()

	c := ipr.NewClient()
	c.HTTPClient = ts.Client()
	c.URL = ts.URL

	_, err := c.GetRanges()
	if err == nil {
		t.Error("want err, got nil")
	}
}

func TestIPRanges_CreatesCSVRecordsOnValidInput(t *testing.T) {
	t.Parallel()

	got := validIPRanges.CSVRecords()
	want := [][]string{
		{"ip_type", "ip_prefix", "region", "service", "network_border_group"},
		{"ipv4", "13.34.65.64/27", "il-central-1", "AMAZON", "il-central-1"},
		{"ipv6", "2600:1ff8:e000::/40", "sa-east-1", "S3", "sa-east-1"},
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

var (
	validResponse = strings.NewReader(`{
		"syncToken": "1676592786",
		"createDate": "2023-02-17-00-13-06",
		"prefixes": [
		  {
			"ip_prefix": "13.34.65.64/27",
			"region": "il-central-1",
			"service": "AMAZON",
			"network_border_group": "il-central-1"
		  },
		  {
			"ipv6_prefix": "2600:1ff8:e000::/40",
			"region": "sa-east-1",
			"service": "S3",
			"network_border_group": "sa-east-1"
		  }
		  ]
		}`)

	invalidResponseWithBogusSyncToken = strings.NewReader(`{
		"syncToken": "",
		"createDate": "2023-02-17-00-13-06",
		"prefixes": [
		  {
			"ip_prefix": "13.34.65.64/27",
			"region": "il-central-1",
			"service": "AMAZON",
			"network_border_group": "il-central-1"
		  },
		  {
			"ipv6_prefix": "2600:1ff8:e000::/40",
			"region": "sa-east-1",
			"service": "S3",
			"network_border_group": "sa-east-1"
		  }
		  ]
		}`)

	invalidResponseWithMalformedDate = strings.NewReader(`{
			"syncToken": "1676592786",
			"createDate": "2023-02-17 00-13-06",
			"prefixes": [
			  {
				"ip_prefix": "13.34.65.64/27",
				"region": "il-central-1",
				"service": "AMAZON",
				"network_border_group": "il-central-1"
			  },
			  {
				"ipv6_prefix": "2600:1ff8:e000::/40",
				"region": "sa-east-1",
				"service": "S3",
				"network_border_group": "sa-east-1"
			  }
			  ]
			}`)

	validIPRanges = ipr.IPRanges{
		SyncToken:  1676592786,
		CreateDate: time.Date(2023, 02, 17, 00, 13, 06, 00, time.UTC),
		IPv4Ranges: []ipr.IPRange{
			{
				Type:               "ipv4",
				IPprefix:           "13.34.65.64/27",
				Region:             "il-central-1",
				Service:            "AMAZON",
				NetworkBorderGroup: "il-central-1",
			},
		},
		IPv6Ranges: []ipr.IPRange{
			{
				Type:               "ipv6",
				IPprefix:           "2600:1ff8:e000::/40",
				Region:             "sa-east-1",
				Service:            "S3",
				NetworkBorderGroup: "sa-east-1",
			},
		},
	}
)
