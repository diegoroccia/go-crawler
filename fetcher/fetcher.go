package fetcher

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type fetcherResult struct {
	body string
	urls []string
}

type httpFetcher map[string]*fetcherResult

func (f httpFetcher) Fetch(url string) (string, []string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", nil, nil
	}
	defer resp.Body.Close()
	tokenizer := html.NewTokenizer(resp.Body)

	urls := make([]string, 0)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				// End of the document, exit the loop
				fmt.Println("Done")
				return "", urls, nil
			} else {
				log.Fatal("HTML parsing error:", err)
				return "", nil, err
			}
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" { // Check if it's an anchor (link) tag
				for _, attr := range token.Attr {
					if attr.Key == "href" {

						if strings.HasPrefix(attr.Val, "/") {
							urls = append(urls, url+attr.Val)
						} else if strings.HasPrefix(attr.Val, url) {
							fmt.Println(attr.Val)
							urls = append(urls, attr.Val)
						}

					}
				}
			}
		}
	}
}

func NewFetcher() Fetcher {
	return &httpFetcher{}
}
