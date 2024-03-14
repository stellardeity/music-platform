package listener

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"rockwall/proto"
)

var itHttp = map[string]bool{
	"GET ": true,
	"POST": true,
	"PUT ": true,
	"DELE": true,
	"OPTI": true,
	"PATC": true,
}

func ItIsHttp(ba []byte) bool {
	return itHttp[string(ba)]
}

func StartListener(node *proto.Node) {
	service := fmt.Sprintf("0.0.0.0%v", node.Address.Port)

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
		go handleConnection(node, conn)
	}
}

func handleConnection(node *proto.Node, conn net.Conn) {
	defer conn.Close()

	log.Printf("New connection: %s", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	readWriter := bufio.NewReadWriter(reader, writer)

	buf, err := readWriter.Peek(4)
	if err != nil {
		if err != io.EOF {
			log.Printf("Read peak ERROR: %s", err)
		}
	}

	if ItIsHttp(buf) {
		handleHttp(readWriter, conn, node)
	} else {
		node.HandleNode(conn)
	}
}

func handleHttp(rw *bufio.ReadWriter, conn net.Conn, node *proto.Node) {
	request, err := http.ReadRequest(rw.Reader)

	if err != nil {
		log.Printf("Read request ERROR: %s", err)
		return
	}

	response := http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	log.Printf("ROCK %s", path.Clean(request.URL.Path))
	if path.Clean(request.URL.Path) == "/ws" {
		handleWs(NewMyWriter(conn), request, node)
	} else {
		processRequest(request, &response)
	}

	err = response.Write(rw)
	if err != nil {
		log.Printf("Write response ERROR: %s", err)
		return
	}

	err = rw.Writer.Flush()
	if err != nil {
		log.Printf("Flush response ERROR: %s", err)
		return
	}
}
