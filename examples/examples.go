package main

import (
	"fmt"

	lightcrawl "github.com/divpatel10/lightcrawl"
)

// Example to crawl Links
func links() {

	var url [2]string
	url[0] = "https://www.google.com"
	url[1] = "https://netflix.com"
	links := lightcrawl.Scrape("a", url[:])

	// Print values to check
	fmt.Println("\n\nFound ", len(links), " urls: ")

	for url, _ := range links {
		fmt.Println("\t" + url)
	}

}

// Example to Crawl Lists
func lists() {
	var url [1]string
	url[0] = "https://web.ics.purdue.edu/~gchopra/class/public/pages/webdesign/05_simple.html"

	links := lightcrawl.Scrape("p", url[:])

	for url, _ := range links {
		fmt.Println("\t" + url)
	}
}

func main() {
	// lists()
	links()
}
