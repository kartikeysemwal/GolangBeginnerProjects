package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var AllRooms RoomMap

func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomID := AllRooms.CreateRoom()

	type resp struct {
		RoomID string `json:"room_id"`
	}

	log.Println(AllRooms.Map)
	json.NewEncoder(w).Encode(resp{RoomID: roomID})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg, 100)

func broadcaster() {
	for {
		msg := <-broadcast

		log.Println("broadcaster msg: ", msg.Message)

		for _, client := range AllRooms.Map[msg.RoomID] {
			if client.Conn != msg.Client {
				err := client.Conn.WriteJSON(msg.Message)

				if err != nil {
					log.Println("this is an error message")
					log.Println(err)
					client.Conn.Close()
				}
			}
		}
	}
}

func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomID"]

	if !ok {
		log.Println("roomID is missing from the URL parameters")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Web socket  upgrade error", err)
	}

	AllRooms.InsertIntoRoom(roomID[0], false, ws)

	// go broadcaster()

	for {
		var msg broadcastMsg

		err := ws.ReadJSON(&msg.Message)

		// if err != nil {
		// 	log.Println("Read error: ", err)
		// 	continue
		// }

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Println("WebSocket connection is going away")
			} else {
				log.Println("Read error:", err)
			}
			break // Exit the loop and function on read error
		}

		msg.Client = ws
		msg.RoomID = roomID[0]

		log.Println("JoinRoomRequestHandler ", msg.Message)

		broadcast <- msg
	}
}

func InitBroadcaster() {
	go broadcaster()
}
