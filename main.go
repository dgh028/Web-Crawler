package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/web-crawler/crawler"
	"github.com/web-crawler/parser"
	"github.com/web-crawler/saver"

	"github.com/baidu/go-lib/log"
	"github.com/baidu/go-lib/log/log4go"
	"github.com/baidu/go-lib/queue"
)

var Depth int
var Seed string
var ParallelNum int //并发爬虫数量

type Scheduler struct {
	TaskQue          queue.Queue //线程安全队列
	Visited          sync.Map
	WorkingCrawerNum chan struct{} //channel若为空表示没有爬虫在工作
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
	if err := initLog("INFO", "./log", true); err != nil {
		fmt.Printf("initLog(): %s\n", err.Error())
		return
	}
	Depth = 2
	Seed = "https://www.baidu.com/"
	ParallelNum = 10
	processPingceTest()
}

func processPingceTest() {
	sch := &Scheduler{
		WorkingCrawerNum: make(chan struct{}, ParallelNum),
	}
	sch.TaskQue.Init()
	sch.Run()
}
func (s *Scheduler) Run() {
	// 初始化种子队列
	seedTask := Task{
		Url:   Seed,
		Depth: 0,
	}
	s.TaskQue.Append(seedTask)
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
			//读取内容。如果读取失败，则重新读取。
			content, err := crawler.GetWebPageContent(task.Url)
			if err != nil {
				return
			}
			//存储数据
			if err := saver.SaveContent(task.Url, content, ""); err != nil {
				return
			}
			//检查深度，如果深度不够，就放到队列里以便继续爬取
			if task.Depth > Depth {
				return
			}
			newUrls, err := parser.GetNewUrlFromWebPage(content)
			if err != nil {
				return
			}
			for _, url := range newUrls {
				s.TaskQue.Append(Task{url, task.Depth + 1}) // 将解析出来的链接加入队列中
			}
		}(task.(Task))
	}
	close(s.WorkingCrawerNum)
	fmt.Println("爬虫完成")
}
