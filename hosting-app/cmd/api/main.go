package main

import (
	"context"
	"encoding/json"
	"hosting-app/cmd/constants"
	"hosting-app/cmd/data"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	UploaderRedisQueue *redis.Client
}

func main() {

	appConfig := Config{
		UploaderRedisQueue: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}

	appConfig.InitUploaderService()

	testReq := data.UploaderRedisStruct{
		RepoURL: "https://github.com/hkirat/react-boilerplate.git",
	}

	payload, err := json.Marshal(testReq)
	if err != nil {
		panic(err)
	}

	if err = appConfig.UploaderRedisQueue.Publish(context.Background(), constants.UPLOAD_TO_CLOUD_CHANNEL, payload).Err(); err != nil {
		// fmt.Println("Unable to publish to", constants.UPLOAD_TO_CLOUD_CHANNEL, "channel")
		panic(err)
	}

	appConfig.UploadRepoToCloud()

	// err := appConfig.DownloadRepo("https://github.com/hkirat/react-boilerplate.git")

	// if err != nil {
	// 	fmt.Println("Download repo failed")
	// }

	// err := BuildRepo("C:\\Important Files\\GolangSmallProjects\\hosting-app\\react-boilerplate\\my-app")
	// if err != nil {
	// 	fmt.Println("Build of repo failed")
	// }
}
