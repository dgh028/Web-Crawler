package saver

import (
	"net/url"
	"os"
	"path/filepath"
)

// 网页存储时每个网页单独存为一个文件，以URL为文件名。注意对URL中的特殊字符，需要做转义
func SaveContent(rawURL string, decodedContent []byte, outPutDir string) (err error) {
	// 转义URL中的特殊字符
	fileName := filepath.Join(outPutDir, url.QueryEscape(rawURL))
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	// 将网页内容写入文件
	_, err = file.Write(decodedContent)
	if err != nil {
		return
	}
	return nil
}
