package lightcrawl

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

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
func crawl(url string, ch chan string, chFin chan bool, etype string) {
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
	tag := map[string]bool{
		"li":    false,
		"a":     false,
		"title": false,
		"td":    false,
		"h1":    false,
		"h2":    false,
		"h3":    false,
		"p":     false,
		// "td":    false,
	}
	// Iterate through all the tokens
	for {
		curToken := z.Next()

		// Map key to true if the specific tag is asked for

		switch {
		// Stop processing if there is an error tokem
		case curToken == html.ErrorToken:
			return

		// StartTag -> eg <html> <a> <body> etc
		case curToken == html.StartTagToken:
			tt := z.Token()

			_, ok := tag[etype]

			if ok {
				tag[etype] = tt.Data == etype
			} else {
				chFin <- true
			}

			// if its not an anchor, just continue
			if tag["a"] {

				// get url from Href from the <a> tag
				ok, a_url := getHref(tt)

				if !ok {
					continue
				}

				//store if the href starts with http
				hasHttp := strings.Index(url, "http") == 0

				// publish the url to the channel
				if hasHttp {
					ch <- a_url
				}
			}

		case curToken == html.TextToken:
			tt := z.Token()

			if tag["li"] || tag["h1"] || tag["h2"] || tag["p"] {
				ch <- tt.Data
			}
		}
	}
}

func Scrape(element string, seedUrls []string) map[string]bool {
	// Map of passed URL and whether URLs were found for the given URL
	foundUrls := make(map[string]bool)

	// This channel is used to output all the found urls
	chUrls := make(chan string)

	// This channel lets us know that that we have found all the URLs
	chFin := make(chan bool)

	// Go over all the URLs in the Seed URLs
	for _, url := range seedUrls {
		// For each URL, start a go routine to scrape a website
		go crawl(url, chUrls, chFin, element)
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

	// Close the channels
	close(chUrls)
	close(chFin)

	return foundUrls
}
