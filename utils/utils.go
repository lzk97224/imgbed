package utils

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"imgbed/utils/pngquant"
	"log"
	"os"
	"os/exec"
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

func AllCompression(ss []string) error {
	for _, item := range ss {
		err := ImgCompression(item, item)
		if err != nil {
			return err
		}
	}
	return nil
}

var pngquantPath = ""

func init() {

	gopath := os.Getenv("GOPATH")
	if len(gopath) <= 0 {
		log.Fatal("环境变量不存在")
	}

	pngquantPath = gopath + "/bin/" + "pngquant"

	_, err := os.Stat(pngquantPath)
	if err == nil {
		//fmt.Println("文件已经存在")
		return
	}

	if err != nil {
		if os.IsNotExist(err) {
			//fmt.Println("文件不存在")
		} else {
			log.Fatal("文件判断失败", err)
		}
	}

	pngquantExe, err := os.OpenFile(pngquantPath, os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("压缩工具复制失败", err)
	}

	defer pngquantExe.Close()
	_, err = pngquantExe.Write(pngquant.PngquantBin)
	if err != nil {
		log.Fatal("压缩工具复制失败", err)
	}
}

func ImgCompression(s string, d string) error {
	imgFile, err := os.Open(s)
	if err != nil {
		return err
	}
	_, format, err := image.Decode(imgFile)
	if err != nil {
		return err
	}
	err = imgFile.Close()
	if err != nil {
		return err
	}

	if format == "png" {
		cmd := exec.Command(pngquantPath, "--force", "--skip-if-larger", s, "-o", d)
		err := cmd.Start()
		if err != nil {
			return err
		}
		err = cmd.Wait()
		if err != nil && err.Error() != "exit status 98" {
			return err
		}
	} else if format == "jpeg" {
		src, err := imaging.Open(s)
		if err != nil {
			return err
		}
		err = imaging.Save(src, d,
			imaging.PNGCompressionLevel(png.BestCompression),
			imaging.JPEGQuality(50),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
