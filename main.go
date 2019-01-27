package main

import "fmt"

const awsurl = "https://ip-ranges.amazonaws.com/ip-ranges.json"

func main() {
	rx, err := getIPRanges(awsurl)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(rx)
}
