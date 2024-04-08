package crawler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

const (
	HeaderKeyUserAgent = "User-Agent"
)

const (
	Mozilla = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.36"
)

// 网络请求获取网页内容
func GetWebPageContent(url string) (content []byte, err error) {
	//单个网页抓取或解析失败，不能导致整个程序退出。需要在日志中记录下错误原因并继续。
	content, contentType, err := getHTTPResp(url)
	if err != nil {
		return nil, err
	}
	if !strings.Contains(contentType, "text") {
		return nil, fmt.Errorf("%s: Content-Type: %s", url, contentType)
	}
	//需要能够处理不同字符编码的网页(例如utf-8或gbk)
	//使用正确的字符编码进行解码
	charsetReader, err := charset.NewReader(bytes.NewReader(content), contentType)
	if err != nil {
		return nil, err
	}
	// 读取已解析的HTML页面内容
	decodedBody, err := ioutil.ReadAll(charsetReader)
	if err != nil {
		return nil, err
	}
	return decodedBody, nil
}

// A successful call returns data, content type of url and err == nil.
func getHTTPResp(url string) (content []byte, contentType string, err error) {
	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置超时时间为 10 秒
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("%s: http.NewRequest(): %s", url, err.Error())
	}
	req.Header.Add(HeaderKeyUserAgent, Mozilla)
	// 发送 GET 请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("http status code is %d", resp.StatusCode)
	}

	// 读取响应内容
	content, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// 获取响应类型
	contentType = resp.Header.Get("Content-Type")

	return content, contentType, nil
}
