package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}

	}
	return
}

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

	fmt.Println("\n\nFound ", len(foundUrls), " urls: ")

	for url, _ := range foundUrls {
		fmt.Println("\t" + url)
	}
	fmt.Print("\n\n")
}

func crawl(url string, ch chan string, chFin chan bool) {
	res, err := http.Get(url)

	defer func() {
		chFin <- true
	}()

	if err != nil {
		fmt.Print("Error: ", err, "\n For URL\t", url)
		return
	}

	body := res.Body

	defer body.Close()

	z := html.NewTokenizer(body)

	for {
		curToken := z.Next()

		switch {

		case curToken == html.ErrorToken:
			return

		case curToken == html.StartTagToken:
			ancTag := z.Token()

			isAnchor := ancTag.Data == "a"

			if !isAnchor {
				continue
			}

			ok, url := getHref(ancTag)

			if !ok {
				continue
			}

			hasProto := strings.Index(url, "http") == 0

			if hasProto {
				ch <- url
			}
		}

	}

}
