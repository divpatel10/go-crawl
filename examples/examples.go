package main

import (
	"fmt"

	lightcrawl "github.com/divpatel10/lightcrawl"
)

// Example to crawl Links
func links() {

	var url [1]string
	url[0] = "https://netflix.com"
	links := lightcrawl.ScrapeElement("a", url[:])

	// Print values to check
	fmt.Println("\n\nFound ", len(links), " urls: ")

	for url, val := range links {

		for furl := range val {
			fmt.Printf("%s \t %s\n", url, links[url][furl])
		}
	}

}

func getElementById() {
	var url string = "https://pkg.go.dev/golang.org/x/net/html"

	lightcrawl.ScrapeId(url, "pkg-index")

}

// Example to Crawl Lists
func lists() {
	var url [1]string
	url[0] = "https://web.ics.purdue.edu/~gchopra/class/public/pages/webdesign/05_simple.html"

	links := lightcrawl.ScrapeElement("p", url[:])

	for url, val := range links {

		for furl := range val {
			fmt.Printf("%s \t %s\n", url, links[url][furl])
		}
	}
}

func main() {
	// lists()
	// links()
	getElementById()
}
