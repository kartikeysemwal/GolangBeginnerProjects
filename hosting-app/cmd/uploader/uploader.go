package uploader

import (
	"errors"
	"fmt"
	"hosting-app/cmd/constants"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
)

// using this file we have to download the repo from github and upload it to some cloud

func DownloadRepo(repoURL string) error {
	splittedURL := strings.Split(repoURL, "/")

	if !strings.Contains(splittedURL[len(splittedURL)-1], ".git") {
		return errors.New("URL does not seems to be formatted")
	}

	savePath := constants.GIT_REPO_DOWNLOAD_PATH + "/" + splittedURL[len(splittedURL)-1]

	err := removeDirIfExists(savePath)

	if err != nil {
		fmt.Println("Error in removing existing dir present at path", savePath, "Error: ", err)
		return err
	}

	_, err = git.PlainClone(savePath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})

	if err != nil {
		return err
	}

	return nil
}

func removeDirIfExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// we are good
		return nil
	}

	// remove the directory
	fmt.Printf("Path: [%s] already exists. Removing\n", path)

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}
