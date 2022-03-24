package main

import (
	"fmt"
	"os"
)

func main() {

	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]

	chUrls := make(chan string)
	chFin := make(chan bool)

	for _, url := range seedUrls {
		go crawl(url, chUrls, chFin)
	}

	for c := 0; c < len(seedUrls); {
		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFin:
			c++

		}
	}

	fmt.Println("\n\nFound ", len(foundUrls), " urls: \n\n")

	for url, _ := range foundUrls {
		fmt.Println("\t" + url)
	}
}
