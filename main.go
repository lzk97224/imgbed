package main

import (
	"flag"
	"fmt"
	"imgbed/gitlab"
	"imgbed/qiniu"
	"log"
	"sync"
	"time"
)

var upType string

func init() {
	flag.StringVar(&upType, "t", "u", "上传类型")
	flag.Parse()
}

func main() {
	//fmt.Println(upType)
	//fmt.Println(flag.Args())
	upload()
}

func upload() {
	now := time.Now()
	var urls []string
	var waitG sync.WaitGroup

	upFileList := flag.Args()
	if len(upFileList) <= 0 {
		log.Panic("上传文件不存在")
	}

	var errQiniu error
	waitG.Add(1)
	go func() {
		defer waitG.Done()
		_, errQiniu = qiniu.Upload(now, upFileList)
	}()

	var errGitlab error
	waitG.Add(1)
	go func() {
		defer waitG.Done()
		urls, errGitlab = gitlab.Upload(now, upFileList)
	}()

	waitG.Wait()

	if errQiniu != nil {
		log.Panic(errQiniu)
	}

	if errGitlab != nil {
		log.Panic(errGitlab)
	}

	for _, url := range urls {
		if upType == "m" {
			fmt.Println("![](" + url + ")")
		} else {
			fmt.Println(url)
		}
	}
}
