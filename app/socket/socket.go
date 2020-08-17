package socket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[int]*websocket.Conn) // connected clients
var broadcast = make(chan Message)          // broadcast channel

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
	UserID   int    `json:"userID"`
	SendTo   int    `json:"sendTo"`
}

func WriteSocketMessage() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		//send private message
		if msg.SendTo != 0 {
			err := clients[msg.SendTo].WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				clients[msg.SendTo].Close()
				delete(clients, msg.SendTo)
			}
			// send public message
		} else {
			for index, client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, index)
				}
			}
		}
	}
}

func readSocketMessage(websocket *websocket.Conn) {
	for {
		var msg Message
		//read new message and it to the msg object
		err := websocket.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		clients[msg.UserID] = websocket
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}

}

func SocketEndPoint(response http.ResponseWriter, request *http.Request) {
	upgrader.CheckOrigin = func(request *http.Request) bool {
		return true
	}
	websocket, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Println(err)
	}
	defer websocket.Close()
	readSocketMessage(websocket)

}
