package main

import (
	"fmt"
	"hosting-app/cmd/builder"
	"hosting-app/cmd/uploader"
)

func main() {
	err := uploader.DownloadRepo("https://github.com/hkirat/react-boilerplate.git")

	if err != nil {
		fmt.Println("Download repo failed")
	}

	err = builder.BuildRepo("C:\\Important Files\\GolangSmallProjects\\hosting-app\\react-boilerplate\\my-app")
	if err != nil {
		fmt.Println("Build of repo failed")
	}
}
