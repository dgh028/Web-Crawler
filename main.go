package main

import (
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/web-crawler/crawler"
	"github.com/web-crawler/loader"
	"github.com/web-crawler/parser"
	"github.com/web-crawler/saver"

	"github.com/baidu/go-lib/log"
	"github.com/baidu/go-lib/log/log4go"
	"github.com/baidu/go-lib/queue"
	"github.com/go-gcfg/gcfg"
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

type Scheduler struct {
	TaskQue          queue.Queue //线程安全队列
	Visited          sync.Map
	WorkingCrawerNum chan struct{} //channel若为空表示没有爬虫在工作
	Cfg              Config
	// 需要存储的目标网页正则表达式
	TargetUrlPattern *regexp.Regexp
}

type Task struct {
	Url   string
	Depth int
}

func initLog(logSwitch string, logPath string, stdOut bool) error {
	log4go.SetLogBufferLength(10000)
	log4go.SetLogWithBlocking(false)

	err := log.Init("web_crawler", logSwitch, logPath, stdOut, "midnight", 5)
	if err != nil {
		return fmt.Errorf("err in log.Init(): %s", err.Error())
	}

	return nil
}
func main() {
	defer func() {
		log.Logger.Close()
		time.Sleep(100 * time.Millisecond)
	}()
	//日志库请使用http://icode.baidu.com/repos/baidu/go-lib/log
	if err := initLog("INFO", "./log", true); err != nil {
		fmt.Printf("initLog(): %s\n", err.Error())
		return
	}
	//主配置文件读取使用https://github.com/go-gcfg/gcfg
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, "./conf/spider.conf")
	if err != nil {
		log.Logger.Error("无法读取配置文件：%v", err)
	}
	//创建目录
	err = os.MkdirAll(cfg.Spider.OutputDirectory, 0755)
	if err != nil {
		fmt.Println("创建目录失败:", err)
		return
	}
	processPingceTest(cfg)
}

func processPingceTest(cfg Config) {
	sch := &Scheduler{
		WorkingCrawerNum: make(chan struct{}, cfg.Spider.ThreadCount),
		Cfg:              cfg,
	}
	targetUrlRegex, err := regexp.Compile(cfg.Spider.TargetUrl)
	if err != nil {
		log.Logger.Error("regexp.Compile: %s", err.Error())
		return
	}
	sch.TargetUrlPattern = targetUrlRegex
	sch.TaskQue.Init()
	sch.Run()
}

func (s *Scheduler) Run() {
	seeds, err := loader.LoadSeed(s.Cfg.Spider.UrlListFile)
	if err != nil {
		log.Logger.Error("loader.LoadSeed: %s", err.Error())
		return
	}
	// 初始化种子队列
	for _, seed := range seeds {
		seedTask := Task{
			Url:   seed,
			Depth: 0,
		}
		s.TaskQue.Append(seedTask)
	}

	for {
		//当程序完成所有抓取任务后，必须优雅退出  即s.TaskQue为空，且爬虫协程均以完成工作退出
		if s.TaskQue.Len() == 0 && len(s.WorkingCrawerNum) == 0 {
			break
		}
		if s.TaskQue.Len() == 0 {
			continue // 队列为空，继续等待爬虫工作完成
		}
		// 从队列中拿到任务
		task := s.TaskQue.Remove()
		// 判断是否已经访问过
		if _, ok := s.Visited.LoadOrStore(task.(Task).Url, struct{}{}); ok {
			continue
		}
		s.WorkingCrawerNum <- struct{}{}
		go func(task Task) {
			defer func() { <-s.WorkingCrawerNum }() // 爬虫工作完成后，将channel中空闲的爬虫数+1
			if err := s.RunTask(task); err != nil {
				log.Logger.Error("%s", err.Error())
			}
		}(task.(Task))
	}
	close(s.WorkingCrawerNum)
	fmt.Println("爬虫完成")
}

func (s *Scheduler) RunTask(task Task) error {
	//读取内容。如果读取失败，则重新读取。
	content, err := crawler.GetWebPageContent(task.Url, s.Cfg.Spider.CrawlTimeout)
	if err != nil {
		return fmt.Errorf("crawler.GetWebPageContent: %s", err.Error())
	}
	//# 需要存储的目标网页URL pattern(正则表达式)  .*.(htm|html)$
	if s.TargetUrlPattern.MatchString(task.Url) {
		//存储数据
		if err := saver.SaveContent(task.Url, content, s.Cfg.Spider.OutputDirectory); err != nil {
			return fmt.Errorf("saver.SaveContent: %s", err.Error())
		}
	}

	//检查深度，如果深度不够，就放到队列里以便继续爬取
	if task.Depth > s.Cfg.Spider.MaxDepth {
		return nil
	}
	newUrls, err := parser.GetNewUrlFromWebPage(content, task.Url)
	if err != nil {
		return fmt.Errorf("parser.GetNewUrlFromWebPage: %s", err.Error())
	}
	for _, url := range newUrls {
		s.TaskQue.Append(Task{url, task.Depth + 1}) // 将解析出来的链接加入队列中
	}
	return nil
}
