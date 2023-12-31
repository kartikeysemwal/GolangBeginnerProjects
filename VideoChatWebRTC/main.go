package main

import (
	"VideoChatWebRTC/server"
	"log"
	"net/http"
)

func main() {
	server.AllRooms.Init()
	server.InitBroadcaster()

	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)

	log.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Println(err)
	}
}
