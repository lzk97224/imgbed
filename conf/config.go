package conf

import (
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var CF config

type config struct {
	Qiniu struct {
		AccessKey string `json:"accessKey"`
		SecretKey string `json:"secretKey"`
		Bucket    string `json:"bucket"`
	} `json:"qiniu"`
	Gitlab struct {
		PrivateToken string `json:"privateToken"`
		UserName     string `json:"userName"`
		ProjectName  string `json:"projectName"`
		Branch       string `json:"branch"`
	} `json:"gitlab"`
}

func init() {
	var confStr []byte
	var err error

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Panic(err)
	}

	cmdPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Panic(err)
	}
	cmdAbsPath, err := filepath.Abs(cmdPath)
	if err != nil {
		log.Panic(err)
	}
	cmdDir := filepath.Dir(cmdAbsPath)

	configFiles := []string{
		"./qnimg_config.json",
		cmdDir + "/qnimg_config.json",
		homeDir + "/qnimg_config.json",
		homeDir + "/go/bin/qnimg_config.json",
		userConfigDir + "/qnimg_config.json",
	}

	for _, fileName := range configFiles {
		confStr, err = os.ReadFile(fileName)
		if err == nil {
			break
		}
	}

	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(confStr, &CF)
	if err != nil {
		log.Panic(err)
	}
}
