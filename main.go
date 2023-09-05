package main

import (
	"flag"
	"fmt"
	"imgbed/conf"
	"imgbed/gitlab"
	"imgbed/qiniu"
	"imgbed/utils"
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

	err := utils.AllCompression(upFileList)
	if err != nil {
		log.Panic("文件压缩是吧")
	}

	var errQiniu error
	var errGitlab error

	if len(conf.CF.Qiniu.AccessKey) > 1 {
		waitG.Add(1)
		go func() {
			defer waitG.Done()
			_, errQiniu = qiniu.Upload(now, upFileList)
		}()
	}

	if len(conf.CF.Gitlab.ProjectName) > 1 {
		waitG.Add(1)
		go func() {
			defer waitG.Done()
			urls, errGitlab = gitlab.Upload(now, upFileList)
		}()
	}

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
