package main

import (
	"fmt"
	"log"
	"net/http"

	"goproj.com/user"
)

const port = "80"

type Server struct {
	userManager user.UserManager
}

func main() {
	// app := user.InitApp()

	// m * testing.M

	// user.TestMain(m)

	userInMemory := user.InitInMemoryUserApp()

	// userDB, err := user.InitSQLiteUserApp("./")

	// if err != nil {
	// 	log.Fatal("Error in iniliatizing SQLite user app", err)
	// }

	server := &Server{
		userManager: userInMemory,
	}

	// _, err := server.userManager.CreateUser(user.User{Name: "John Smith", Email: "john@gmail.com"})

	// if err != nil {
	// 	fmt.Println(err)
	// }

	log.Printf("Starting go-task service on port %s\n", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: server.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
