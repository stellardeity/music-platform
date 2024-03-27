package discover

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"rockwall/proto"
	"strings"
	"time"
)

func StartDiscover(node *proto.Node) {
	go startFuck("224.0.0.1:35035", node)
	go listenFuck("224.0.0.1:35035", node)
}

func listenFuck(address string, node *proto.Node) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		buffer := make([]byte, 1024)
		conn.ReadFromUDP(buffer)

		trim := bytes.Trim(buffer, "\x00")

		parts := strings.Split(string(trim), "::")

		_, found := node.Connections[address]
		if found || node.Address.Port == ":"+parts[1] {
			continue
		}
		// log.Printf("%s", trim)

		node.Connections[":"+string(parts[1])] = true

	}
}

func startFuck(address string, node *proto.Node) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Printf(err.Error())
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf(err.Error())
	}

	for {
		_, err := conn.Write([]byte(fmt.Sprintf("fuck:%v", node.Address.Port)))
		if err != nil {
			log.Printf(err.Error())
		}
		time.Sleep(1 * time.Second)
	}
}
