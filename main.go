package main

import (
	"fmt"
	"sync"
)

var visited sync.Map
var wt sync.WaitGroup //使用sync.WaitGroup来同步
var maxConcurrency = 3
var sem = make(chan struct{}, maxConcurrency) //使用channel来控制并发数

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// Decrement the counter when the goroutine completes.
	defer wt.Done()
	// TODO: Don't fetch the same URL twice.
	_, ok := visited.Load(url)
	if ok {
		return
	}
	visited.Store(url, true)
	if depth <= 0 {
		return
	}
	sem <- struct{}{}
	defer func() {
		<-sem
	}()
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		// Increment the WaitGroup counter.
		wt.Add(1)
		go Crawl(u, depth-1, fetcher)
	}
	return
}

func main() {
	wt.Add(1)
	Crawl("https://golang.org/", 4, fetcher)
	//time.Sleep(5 * time.Second)
	// Wait for all goroutine to complete.
	wt.Wait()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
