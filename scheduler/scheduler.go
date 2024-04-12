package scheduler

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/baidu/go-lib/log"
	"github.com/baidu/go-lib/queue"
	"github.com/web-crawler/loader"
)

type Scheduler struct {
	TaskQue          queue.Queue //线程安全队列
	Visited          sync.Map
	WorkingCrawerNum chan struct{} //channel若为空表示没有爬虫在工作
	Cfg              loader.Config
	// 需要存储的目标网页正则表达式
	TargetUrlPattern *regexp.Regexp
	// 站点爬取间隔timer表
	TimerTable sync.Map
}

func NewScheduler(cfg *loader.Config) (*Scheduler, error) {
	// 初始化爬虫调度器
	sch := &Scheduler{
		WorkingCrawerNum: make(chan struct{}, cfg.Spider.ThreadCount),
		Cfg:              *cfg,
	}
	targetUrlRegex, err := regexp.Compile(cfg.Spider.TargetUrl)
	if err != nil {
		return nil, fmt.Errorf("regexp.Compile: %s", err.Error())
	}
	sch.TargetUrlPattern = targetUrlRegex
	sch.TaskQue.Init()
	return sch, nil
}

func (s *Scheduler) Run() {
	// 初始化种子队列
	seeds, err := loader.LoadSeed(s.Cfg.Spider.UrlListFile)
	if err != nil {
		log.Logger.Error("loader.LoadSeed: %s", err.Error())
		return
	}
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
			defer func() {
				log.Logger.Info("depth %d ,task %s done", task.Depth, task.Url)
				<-s.WorkingCrawerNum
			}() // 爬虫工作完成后，将channel中空闲的爬虫数+1
			if err := s.RunTask(task); err != nil {
				log.Logger.Error("%s", err.Error())
			}
		}(task.(Task))
	}
	close(s.WorkingCrawerNum)
	log.Logger.Info("所有爬虫任务均完成")
}
