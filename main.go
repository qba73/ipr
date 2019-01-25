package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const awsurl = "https://ip-ranges.amazonaws.com/ip-ranges.json"


func main()  {
	res, err := http.Get(awsurl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", awsurl, err)
		os.Exit(1)
	}

	fmt.Printf("%s", b)
}
