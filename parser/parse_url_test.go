package parser

import (
	"testing"

	"github.com/web-crawler/crawler"
)

func TestGetNewUrlFromWebPage(t *testing.T) {
	url := "https://www.baidu.com"
	content, err := crawler.GetWebPageContent(url, 1)
	if err != nil {
		t.Errorf("crawler.GetWebPageContent: %s", err.Error())
	}
	urls, err := GetNewUrlFromWebPage(content, url)
	if err != nil {
		t.Errorf("parser.GetNewUrlFromWebPage: %s", err.Error())
	}
	if len(urls) == 0 {
		t.Errorf("no sublink in %s", url)
		return
	}
}
