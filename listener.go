package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func startListener(port int) {
	if port <= 0 || port > 65535 {
		port = 35035
	}

	service := fmt.Sprintf("0.0.0.0:%v", port)

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		log.Printf("ResolveTCPAddr: %s", err.Error())
		os.Exit(1)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Printf("ListenTCP: %s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("\n\tService start on %s\n\n", tcpAddr.String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go onConnection(conn)
	}
}

func onConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("New connection from: %v", conn.RemoteAddr().String())

	NewPeer(conn)
}
