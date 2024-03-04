package listener

import (
	"net"
	"stellard/proto"
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
	listen, err := net.Listen("tcp", "0.0.0.0"+node.Address.Port)
	if err != nil {
		panic("listen error")
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			break
		}
		go handleConnection(node, conn)
	}
}

func handleConnection(node *proto.Node, conn net.Conn) {
	defer conn.Close()

	// reader := bufio.NewReader(conn)
	// writer := bufio.NewWriter(conn)

	// readWriter := bufio.NewReadWriter(reader, writer)

	// buf, err := readWriter.Peek(4)
	// if err != nil {
	// 	if err != io.EOF {
	// 		log.Printf("Read peak ERROR: %s", err)
	// 	}
	// 	return
	// }

	// if ItIsHttp(buf) {
	// 	handleHttp(readWriter)
	// 	return
	// } else {
	node.HandleNode(conn)
	// }
}

// func handleHttp(rw *bufio.ReadWriter) {
// 	request, _ := http.ReadRequest(rw.Reader)
// 	log.Printf("%s", request.URL)
// }
