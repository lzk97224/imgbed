package qiniu

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"imgbed/conf"
	"imgbed/utils"
	"time"
)

var mac = qbox.NewMac(conf.CF.Qiniu.AccessKey, conf.CF.Qiniu.SecretKey)
var putPolicy = storage.PutPolicy{
	Scope: conf.CF.Qiniu.Bucket,
}
var upToken = putPolicy.UploadToken(mac)

func Upload(now time.Time, fileArgs []string) ([]string, error) {

	var cfg = storage.Config{}
	// 空间对应的机房
	cfg.Region = &storage.ZoneHuabei
	// 是否使用https域名
	cfg.UseHTTPS = true

	bucketManager := storage.NewBucketManager(mac, &cfg)
	domainInfoList, err := bucketManager.ListBucketDomains(conf.CF.Qiniu.Bucket)
	if err != nil {
		return nil, err
	}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	urls := []string{}
	for i, fileName := range fileArgs {
		key := utils.CreateFullName("", now, fileName, i)

		err := formUploader.PutFile(context.Background(), &ret, upToken, key, fileName, nil)
		if err != nil {
			return nil, err
		}

		url := "http://" + domainInfoList[0].Domain + "/" + ret.Key
		urls = append(urls, url)
	}
	return urls, nil
}
