package main

import (
	"flag"
	"os"
	"os/user"
)

type InitParams struct {
	Name      *string
	Port      *int
	PeersFile *string
}

var initParams InitParams

func init() {
	currentUser, _ := user.Current()
	hostName, _ := os.Hostname()

	initParams = InitParams{
		Name:      flag.String("name", currentUser.Username+"@"+hostName, "you name"),
		Port:      flag.Int("port", 35035, "port that have to listen"),
		PeersFile: flag.String("peers", "peers.txt", "Path to file with peer addresses on each line"),
	}
}

func main() {
	startListener(*initParams.Port)
}
