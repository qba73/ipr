package main

import (
	"fmt"
	"os"
)

const awsurl = "https://ip-ranges.amazonaws.com/ip-ranges.json"


func main()  {
	b, err := fetchData(awsurl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", awsurl, err)
		os.Exit(1)
	}

	fmt.Printf("%s", b)
}
