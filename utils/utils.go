package utils

import (
	"fmt"
	"path"
	"time"
)

func CreateFullName(prefix string, now time.Time, fileName string, index int) string {
	return path.Join(CreatePath(prefix, now), CreateName(now, fileName, index))
}

func CreatePath(prefix string, now time.Time) string {
	result := ""
	if len(prefix) >= 1 {
		result += prefix + "/"
	}
	result += now.Format("2006") + "/" + now.Format("01")
	return result
}

func CreateName(now time.Time, fileName string, index int) string {
	return now.Format("20060102_150405") + "_" + fmt.Sprintf("%0.3d", now.UnixMilli()%1000) + "_" + fmt.Sprintf("%d", index) + path.Ext(fileName)
}
