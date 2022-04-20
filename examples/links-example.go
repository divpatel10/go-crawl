package main

import (
	"fmt"

	lightcrawl "github.com/divpatel10/lightcrawl"
)

func main() {

	var url [2]string
	url[0] = "https://google.com"
	url[1] = "https://netflix.com"
	links := lightcrawl.Scrape("a", url[:])

	// Print values to check
	for url, _ := range links {
		fmt.Println("\t" + url)
	}

}
