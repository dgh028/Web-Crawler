package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func LoadSeed(path string) ([]string, error) {
	// 读取种子文件
	seedData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取种子文件失败：%v", err)
	}

	// 解析 JSON
	var seeds []string
	err = json.Unmarshal(seedData, &seeds)
	if err != nil {
		return nil, fmt.Errorf("解析 JSON 失败：%v", err)
	}
	return seeds, nil
}
