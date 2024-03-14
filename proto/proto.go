package proto

import (
	"encoding/json"
	"log"
	"net"
	"strings"
)

type Node struct {
	Connections map[string]bool
	Address     Address
}

type Address struct {
	IPv4 string
	Port string
}

type Package struct {
	To   string
	From string
	Date string
}

func NewNode(address string) *Node {
	splited := strings.Split(address, ":")
	if len(splited) != 2 {
		return nil
	}

	return &Node{
		Connections: make(map[string]bool),
		Address: Address{
			IPv4: splited[0],
			Port: ":" + splited[1],
		},
	}
}

func (node *Node) PrintNetwork() {
	for addr := range node.Connections {
		log.Println("|", addr)
	}
}

func (node *Node) HandleNode(conn net.Conn) {
	var (
		buffer  = make([]byte, 512)
		message string
		pack    *Package
	)

	for {
		length, err := conn.Read(buffer)
		if err != nil {
			break
		}
		message += string(buffer[:length])
	}
	err := json.Unmarshal([]byte(message), &pack)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	node.ConnectTo([]string{pack.From})
	log.Println(pack.Date)
}

func (node *Node) SendMessageToAll(message string) {
	var new_pack = &Package{
		From: node.Address.IPv4 + node.Address.Port,
		Date: message,
	}
	for addr := range node.Connections {
		new_pack.To = addr
		node.Send(new_pack)
	}

}

func (node *Node) Send(pack *Package) {
	conn, err := net.Dial("tcp", pack.To)
	if err != nil {
		delete(node.Connections, pack.To)
		return
	}

	defer conn.Close()

	json_pack, _ := json.Marshal(*pack)
	conn.Write(json_pack)
}

func (node *Node) ConnectTo(addresses []string) {
	for _, addr := range addresses {
		node.Connections[addr] = true
	}
}
