package gitlab

import (
	"encoding/base64"
	"errors"
	"fmt"
	gitlab "github.com/xanzy/go-gitlab"
	"log"
	"os"
	"qnimg/conf"
	"qnimg/utils"
	"time"
)

var encoding = "base64"
var commitMessage = "上传图片"

func Upload(now time.Time, fileArgs []string) ([]string, error) {
	git, err := gitlab.NewClient(conf.CF.Gitlab.PrivateToken)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	users, _, err := git.Users.ListUsers(&gitlab.ListUsersOptions{Username: &conf.CF.Gitlab.UserName})
	if len(users) != 1 {
		return nil, errors.New(fmt.Sprintf("users 查询错误,%v", users))
	}
	userId := users[0].ID

	projects, _, err := git.Projects.ListUserProjects(userId, &gitlab.ListProjectsOptions{
		Search: &conf.CF.Gitlab.ProjectName,
	})
	if len(projects) != 1 {
		return nil, errors.New(fmt.Sprintf("projects 查询错误,%v", projects))
	}
	pID := projects[0].ID

	urls := []string{}
	for i, fileName := range fileArgs {
		key := utils.CreateFullName("", now, fileName, i)

		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
		encodeToString := base64.StdEncoding.EncodeToString(fileContent)

		file, _, err := git.RepositoryFiles.CreateFile(pID, key, &gitlab.CreateFileOptions{
			Branch:        &conf.CF.Gitlab.Branch,
			Encoding:      &encoding,
			Content:       &encodeToString,
			CommitMessage: &commitMessage,
		})
		if err != nil {
			return nil, err
		}
		urls = append(urls, fmt.Sprintf("https://gitlab.com/"+conf.CF.Gitlab.UserName+"/"+conf.CF.Gitlab.ProjectName+"/raw/main/"+file.FilePath))
	}

	return urls, nil
}
