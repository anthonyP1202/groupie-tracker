package Websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message est une structure pour les messages WebSocket
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		// read message from client
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println(err)
			break
		}

		// Traitement du message JSON
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("Error decoding JSON message:", err)
			continue
		}

		switch msg.Type {
		case "start_game":
			// Traiter le message pour démarrer une nouvelle partie de Petit Bac
			// Vous pouvez appeler la fonction correspondante de votre package Game ici
			// Game.StartGameHandler(w, r)
		case "submit_response":
			// Traiter le message pour soumettre une réponse au Petit Bac
			// Vous pouvez appeler la fonction correspondante de votre package Game ici
			// Game.SubmitResponseHandler(w, r)
		}

		//send message to client
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			break
		}
	}
}
