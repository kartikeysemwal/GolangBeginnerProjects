package main

import "hosting-app/cmd/uploader"

func main() {
	uploader.DownloadRepo("https://github.com/hkirat/react-boilerplate.git")
}
