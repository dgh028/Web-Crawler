package loader

import (
	"fmt"
	"regexp"

	"gopkg.in/gcfg.v1"
)

// Config 结构体用于保存配置信息
type Config struct {
	Spider struct {
		// 种子文件路径
		UrlListFile string `gcfg:"urlListFile"`
		// 抓取结果存储目录
		OutputDirectory string `gcfg:"outputDirectory"`
		// 最大抓取深度(种子为0级)
		MaxDepth int `gcfg:"maxDepth"`
		// 抓取间隔. 单位: 秒
		CrawlInterval int `gcfg:"crawlInterval"`
		// 抓取超时. 单位: 秒
		CrawlTimeout int `gcfg:"crawlTimeout"`
		// 需要存储的目标网页URL pattern(正则表达式)
		TargetUrl string `gcfg:"targetUrl"`
		// 抓取routine数
		ThreadCount int `gcfg:"threadCount"`
	}
}

func LoadConfig(cfg *Config) error {
	err := gcfg.ReadFileInto(cfg, "./conf/spider.conf")
	if err != nil {
		return fmt.Errorf("无法读取配置文件：%s", err.Error())
	}
	if err := cfg.Check(); err != nil {
		return fmt.Errorf("配置文件错误：%s", err.Error())
	}
	return nil
}

func (cfg *Config) Check() error {
	if cfg.Spider.UrlListFile == "" {
		return fmt.Errorf("cfg.Spider.UrlListFile is empty")
	}
	if cfg.Spider.OutputDirectory == "" {
		return fmt.Errorf("cfg.Spider.OutputDirectory is empty")
	}
	if cfg.Spider.MaxDepth < 0 {
		return fmt.Errorf("cfg.Spider.MaxDepth is invalid")
	}
	if cfg.Spider.CrawlTimeout < 0 {
		return fmt.Errorf("cfg.Spider.CrawlTimeout is invalid")
	}
	if cfg.Spider.CrawlInterval < 0 {
		return fmt.Errorf("cfg.Spider.CrawlInterval is invalid")
	}
	_, err := regexp.Compile(cfg.Spider.TargetUrl)
	if err != nil {
		return fmt.Errorf("%s: regexp.Compile(): %s", cfg.Spider.TargetUrl, err.Error())
	}
	if cfg.Spider.ThreadCount < 1 {
		return fmt.Errorf("cfg.Spider.ThreadCount is invalid")
	}
	return nil
}
