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
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" && attr.Val != "javascript:;" && attr.Val != "javascript:void(0)" {
					//从HTML提取链接时需要处理相对路径和绝对路径
					// 处理相对路径和绝对路径
					if strings.HasPrefix(attr.Val, "http://") || strings.HasPrefix(attr.Val, "https://") {
						// 绝对路径
						links = append(links, attr.Val)
					} else {
						// 相对路径
						linkURL, err := baseURL.Parse(attr.Val)
						if err != nil {
							fmt.Println("baseURL.Parse failed", err)
							continue
						}
						links = append(links, linkURL.String())
					}
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
