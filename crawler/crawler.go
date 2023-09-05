package crawler

import (
	"crawler/fetcher"
	"fmt"
	"sync"
)

type Cache struct {
	mu  sync.Mutex
	url map[string]bool
}

type Crawler struct {
	Url     string
	Fetcher fetcher.Fetcher
}

func (crawler Crawler) CrawlerJob(url string, ch chan string, wg *sync.WaitGroup, c *Cache) {
	defer wg.Done() // Decrement the wait group when this goroutine finishes

	_, urls, err := crawler.Fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch <- url // Send the URL to the channel

	for _, u := range urls {
		if !c.url[u] {
			c.mu.Lock()
			c.url[u] = true
			c.mu.Unlock()

			wg.Add(1)
			go crawler.CrawlerJob(u, ch, wg, c)
		}
	}
}

func (crawler Crawler) Crawl() []string {
	ch := make(chan string) // Buffered channel to prevent deadlock
	wg := &sync.WaitGroup{}
	wg.Add(1) // Increment the wait group for the initial crawl

	var c Cache
	c.url = make(map[string]bool)

	go crawler.CrawlerJob(crawler.Url, ch, wg, &c)

	go func() {
		wg.Wait() // Wait for all goroutines to finish
		close(ch) // Close the channel after all goroutines are done
	}()

	urls := make([]string, 0)

	for url := range ch {
		fmt.Println("Received:", url)
		urls = append(urls, url)
	}

	return urls
}
