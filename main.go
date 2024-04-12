package main

import (
	"fmt"
	"os"
	"time"

	"github.com/baidu/go-lib/log"
	"github.com/baidu/go-lib/log/log4go"

	"github.com/web-crawler/loader"
	"github.com/web-crawler/scheduler"
)

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
	cfg := &loader.Config{}
	if err := loader.LoadConfig(cfg); err != nil {
		fmt.Printf("loader.LoadConfig(): %s\n", err.Error())
		return
	}

	//创建目录
	if err := os.MkdirAll(cfg.Spider.OutputDirectory, 0755); err != nil {
		fmt.Println("创建目录失败:", err)
		return
	}

	//初始化调度器
	sch, err := scheduler.NewScheduler(cfg)
	if err != nil {
		fmt.Printf("scheduler.NewScheduler(): %s\n", err.Error())
		return
	}

	// 启动爬虫
	sch.Run()
}
