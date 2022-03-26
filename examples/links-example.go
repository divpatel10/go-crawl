package main

import (
	"fmt"

	lightcrawl "github.com/divpatel10/go-crawl"
)

func main() {

	var url [1]string
	url[0] = "https://google.com"

	links := lightcrawl.FromLink(url[:], "a")

	// Print values to check
	for url, _ := range links {
		fmt.Println("\t" + url)
	}

}
