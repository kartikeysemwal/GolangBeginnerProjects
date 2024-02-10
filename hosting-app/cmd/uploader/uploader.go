package uploader

import (
	"fmt"
	"hosting-app/cmd/constants"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
)

// using this file we have to download the repo from github and upload it to some cloud

func DownloadRepo(repoURL string) {
	splittedURL := strings.Split(repoURL, "/")

	if !strings.Contains(splittedURL[len(splittedURL)-1], ".git") {
		fmt.Println("URL does not seems to be formatted")
		return
	}

	_, err := git.PlainClone(constants.GIT_REPO_DOWNLOAD_PATH+"/"+splittedURL[len(splittedURL)-1], false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	if err != nil {
		fmt.Println("Something went wrong")
	}
}
