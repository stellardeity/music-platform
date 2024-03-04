package discover

import (
	"bufio"
	"os"
	"stellard/proto"
	"strings"
)

func StartDiscover(node *proto.Node) {
	for {
		msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		message := strings.Replace(msg, "\n", "", 1)

		splited := strings.Split(message, " ")

		switch splited[0] {
		case "/exit":
			os.Exit(0)
		case "/connect":
			node.ConnectTo(splited[1:])
		case "/network":
			node.PrintNetwork()
		default:
			node.SendMessageToAll(message)
		}
	}
}
