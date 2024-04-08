package saver

import (
	"testing"

	"github.com/web-crawler/crawler"
)

func TestSaveContent(t *testing.T) {
	// 示例URL
	rawURL := "https://zhuanlan.zhihu.com/p/643739624"
	content, err := crawler.GetWebPageContent(rawURL)
	if err != nil {
		t.Error(err)
	}
	//存储数据
	if err := SaveContent(rawURL, content, "./"); err != nil {
		t.Error(err)
	}
}
