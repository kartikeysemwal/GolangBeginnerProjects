package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"hosting-app/cmd/constants"
	"hosting-app/cmd/data"

	"github.com/go-git/go-git/v5"
)

func (app *Config) InitUploaderService() error {
	// go app.UploadRepoToCloud()

	return nil
}

// using this file we have to download the repo from github and upload it to some cloud

func (app *Config) UploadRepoToCloud() {
	subscriber := app.UploaderRedisQueue.Subscribe(context.Background(), constants.UPLOAD_TO_CLOUD_CHANNEL)
	defer subscriber.Close()

	for {
		msg, err := subscriber.ReceiveMessage(context.Background())

		if err != nil {
			fmt.Println("Received some bad request from redis queue")
			continue
		}

		uploadReq := data.UploaderRedisStruct{}

		if err = json.Unmarshal([]byte(msg.Payload), &uploadReq); err != nil {
			fmt.Println("Error while parsing the data")
			continue
		}

		DownloadRepo(uploadReq.RepoURL)

		app.PublishToDownloadQueue()
	}
}

func DownloadRepo(repoURL string) error {
	if len(repoURL) == 0 {
		return errors.New("Got empty repo URL for download")
	}

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

func (app *Config) PublishToDownloadQueue() {
}
