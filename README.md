# Web-Crawler

## 1. a toy web crawler

[[Exercise: Web Crawler](https://go.dev/tour/concurrency/10) from a tour of Go]

In this exercise you'll use Go's concurrency features to parallelize a web crawler.且并发量不能超过 3

Modify the `Crawl` function to fetch URLs in parallel without fetching the same URL twice.

_Hint_: you can keep a cache of the URLs that have been fetched on a map, but maps alone are not safe for concurrent use!
