package crawler

import (
	"testing"
)

func TestGetHTTPResp(t *testing.T) {
	_, contentType, err := getHTTPResp("https://www.baidu.com", 1)
	if err != nil {
		t.Error(err)
	} else if contentType == "text/html; charset=utf-8" {
		t.Log("success")
	}
}

func TestGetWebPageContent(t *testing.T) {}
