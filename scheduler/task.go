package scheduler

import (
	"fmt"
	"time"

	"github.com/baidu/go-lib/log"
	"github.com/web-crawler/crawler"
	"github.com/web-crawler/parser"
	"github.com/web-crawler/saver"
)

type Task struct {
	Url   string
	Depth int
}

func (s *Scheduler) RunTask(task Task) error {
	//控制抓取间隔，避免对方网站封禁IP    设计上可以针对单个routine控制，也可以针对站点进行控制
	hostName, err := parser.ParseHostName(task.Url)
	if err != nil {
		return fmt.Errorf("parser.ParseHostName: %s", err.Error())
	}
	lastAccessTime, ok := s.TimerTable.LoadOrStore(hostName, time.Now())
	if ok {
		<-time.After(time.Duration(s.Cfg.Spider.CrawlInterval)*time.Second - time.Since(lastAccessTime.(time.Time)))
	}
	log.Logger.Info("depth %d, start to crawl %s", task.Depth, task.Url)
	//读取内容。如果读取失败，则重新读取。
	content, err := crawler.GetWebPageContent(task.Url, s.Cfg.Spider.CrawlTimeout)
	if err != nil {
		return fmt.Errorf("crawler.GetWebPageContent: %s", err.Error())
	}
	// 需要存储的目标网页URL pattern(正则表达式)  .*.(htm|html)$
	if s.TargetUrlPattern.MatchString(task.Url) {
		//存储数据
		if err := saver.SaveContent(task.Url, content, s.Cfg.Spider.OutputDirectory); err != nil {
			return fmt.Errorf("saver.SaveContent: %s", err.Error())
		}
	}

	//检查深度，如果深度不够，就放到队列里以便继续爬取
	if task.Depth >= s.Cfg.Spider.MaxDepth {
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
