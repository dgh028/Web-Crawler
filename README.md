# Web-Crawler

## 1. a toy web crawler

[[Exercise: Web Crawler](https://go.dev/tour/concurrency/10) from a tour of Go]

In this exercise you'll use Go's concurrency features to parallelize a web crawler.且并发量不能超过 3

Modify the `Crawl` function to fetch URLs in parallel without fetching the same URL twice.

_Hint_: you can keep a cache of the URLs that have been fetched on a map, but maps alone are not safe for concurrent use!

## 2. a web crawler

### 2.1 调度主要逻辑

项目要求程序完成所有抓取任务后，必须优雅退出 （退出的条件是任务队列为空且爬虫协程均以完成工作退出）

每次循环，从任务队列取出任务，然后拿一个未工作的爬虫协程去执行任务

### 2.2 任务执行主逻辑

第一步，从队列中拿到任务。

第二步，读取内容。如果读取失败，则重新读取。如果读取成功，则执行第三步。

第三步，存储数据。

第四步，检查数据深度。

第五步，如果数据深度不足，就进一步解析，并且放到队列中。

第六步，结束任务。
