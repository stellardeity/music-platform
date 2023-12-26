package main

import (
	"log"
	"net"
)

func main() {
	startDiscover()
}

func startMeow(address string) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		log.Printf(err.Error())
	}

	_, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Printf(err.Error())
	}

	for {
		// _, err := conn.Write([]byte("meow"))
		// if err != nil {
		// 	log.Printf(err.Error())
		// }
		// time.Sleep(1 * time.Second)
	}
}

func startDiscover( /*peersFile string*/ ) {
	startMeow("127.0.0.1:35035")
}

func handShake(conn net.Conn) {
	log.Printf("DISCOVERY: try handshake with %s", conn.RemoteAddr())
}
