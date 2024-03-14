package listener

import (
	"encoding/json"
	"log"
	"net/http"

	"rockwall/proto"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWs(w http.ResponseWriter, r *http.Request, node *proto.Node) {
	c, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	log.Printf("Ws started")

	br := make(chan bool)

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("ws read error: %v", err)
			break
		}
		log.Printf("ws read: [%v] %s", mt, message)

		decodedMessage := &proto.WsMessage{}
		err = json.Unmarshal(message, decodedMessage)

		if err != nil {
			log.Printf("error on unmarshal message: %v", err)
			continue
		}

		writeToWs(c, mt, message)
		var new_pack = &proto.Package{
			From: node.Address.IPv4 + node.Address.Port,
			Date: decodedMessage.Content,
		}
		node.Send(new_pack)
	}

	br <- true
}

func writeToWs(c *websocket.Conn, mt int, message []byte) {
	err := c.WriteMessage(mt, message)
	if err != nil {
		log.Printf("ws write error: %s", err)
	}
}
