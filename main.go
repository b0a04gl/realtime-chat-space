package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   = make(map[*websocket.Conn]string)
	broadcast = make(chan Message)
	mu        sync.Mutex
)

type Message struct {
	Text      string `json:"text"`
	Sender    string `json:"sender"`
	Broadcast bool   `json:"broadcast"`
}

func main() {
	fs := http.FileServer(http.Dir("static"))
http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)
		go handleMessages()
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = "$"
	mu.Unlock()

	fmt.Println("new client... " + conn.LocalAddr().String())
	for {
			var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Print(err)
				delete(clients, conn)
			return
		}

		if clients[conn] == "$" {
			clients[conn] = msg.Sender
			msg.Sender = "System"
		}

	broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		fmt.Println("sender... " + msg.Sender)

		for client := range clients {

				clientUID := clients[client]

			fmt.Println("client.... id " + clientUID)

				if strings.Contains(msg.Text, clientUID) || msg.Sender == clientUID {
				continue
			}

			err := client.WriteJSON(msg)
			if err != nil {
					log.Printf("Error: %v", err)
			client.Close()
				delete(clients, client)
			}
		}

	}
}
