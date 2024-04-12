package parser

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func GetNewUrlFromWebPage(content []byte, preUrl string) (urls []string, err error) {
	//html解析请使用 https://go.googlesource.com/net/
	rootNode, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("%s: html.Parse(): %s", preUrl, err.Error())
	}
	// parse url
	baseURL, err := url.ParseRequestURI(preUrl)
	if err != nil {
		return nil, fmt.Errorf("%s: url.ParseRequestURL(): %s", preUrl, err.Error())
	}
	urls = extractLinks(rootNode, baseURL)
	return
}

// 提取链接
func extractLinks(node *html.Node, baseURL *url.URL) []string {
	var links []string

	// 递归遍历节点
	var visitNode func(*html.Node)
	visitNode = func(n *html.Node) {
		/*
			<div>
				<a href="https://www.creatorblue.com/">创蓝科技</a>
				<br/>
				<a href="https://www.budaos.com/">布道师</a>
			</div>
		*/
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" && attr.Val != "javascript:;" && attr.Val != "javascript:void(0)" {
					linkURL, err := baseURL.Parse(attr.Val)
					if err != nil {
						continue
					}
					links = append(links, linkURL.String())
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			visitNode(child)
		}
	}
	visitNode(node)

	return links
}

// Parse hostname from raw url.
func ParseHostName(rawUrl string) (string, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	if u.Host == "" {
		return "", fmt.Errorf("empty host")
	}

	// 可能出现如xxx.baidu.com:8080这样带端口号的情况
	hostName := strings.Split(u.Host, ":")
	if len(hostName) == 0 {
		return "", fmt.Errorf("invalid hostname")
	}

	return hostName[0], nil
}
