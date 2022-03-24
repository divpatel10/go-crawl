package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Function that takes an html token and gets the href value from the tag
func getHref(t html.Token) (ok bool, href string) {
	// iterate through the token attributes
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

// This function crawls through a given url
func crawl(url string, ch chan string, chFin chan bool) {
	// make a Get request for the URL and store the response
	res, err := http.Get(url)

	// Run this at the end of the funciton
	defer func() {
		chFin <- true
	}()

	if err != nil {
		fmt.Print("Error: ", err, "\n For URL\t", url)
		return
	}

	// Get the body of the response
	body := res.Body

	// Close the body at the end of the function
	defer body.Close()

	// Divide the html body into tokens
	z := html.NewTokenizer(body)

	// Iterate through all the tokens
	for {
		curToken := z.Next()

		switch {
		// Stop processing if there is an error tokem
		case curToken == html.ErrorToken:
			return

		// StartTag -> eg <html> <a> <body> etc
		case curToken == html.StartTagToken:
			ancTag := z.Token()
			// if its an anchor tag:
			isAnchor := ancTag.Data == "a"

			// if its not an anchor, just continue
			if !isAnchor {
				continue
			}

			// get url from Href from the <a> tag
			ok, url := getHref(ancTag)

			if !ok {
				continue
			}

			//store if the href starts with http
			hasProto := strings.Index(url, "http") == 0

			// publish the url to the channel
			if hasProto {
				ch <- url
			}
		}
	}
}

func main() {

	// Map of passed URL and whether URLs were found for the given URL
	foundUrls := make(map[string]bool)

	// this variable stores the urls passed by the user
	seedUrls := os.Args[1:]

	// This channel is used to output all the found urls
	chUrls := make(chan string)

	// This channel lets us know that that we have found all the URLs
	chFin := make(chan bool)

	// Go over all the URLs in the Seed URLs
	for _, url := range seedUrls {

		// For each URL, start a go routine to scrape a website
		go crawl(url, chUrls, chFin)

	}

	// Go over all Urls, and subscribe to the channels
	for c := 0; c < len(seedUrls); {
		select {

		// if its a url channel, change the foundUrls map value to true
		case url := <-chUrls:
			foundUrls[url] = true

		// if a channel is finished outputting, move on to the next channel
		case <-chFin:
			c++
		}
	}

	fmt.Println("\n\nFound ", len(foundUrls), " urls: ")

	// Print all urls from the found URLs map
	for url, _ := range foundUrls {
		fmt.Println("\t" + url)
	}
	fmt.Print("\n\n")

	// Close the channels
	close(chUrls)
	close(chFin)
}
