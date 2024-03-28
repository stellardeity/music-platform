package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"os/user"
	"rockwall/discover"
	"rockwall/listener"
	"rockwall/proto"
	"sync"
	"syscall"
)

type InitParams struct {
	Name *string
	Port *int
}

var initParams InitParams

func init() {
	currentUser, _ := user.Current()
	hostName, _ := os.Hostname()

	initParams = InitParams{
		Name: flag.String("name", currentUser.Username+"@"+hostName, "you name"),
		Port: flag.Int("port", 35035, "port that have to listen"),
	}

	flag.Parse()
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

}

func main() {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		log.Printf("Exit by signal: %s", sig)
		os.Exit(1)
	}()

	p := proto.NewProto(*initParams.Name, *initParams.Port)

	var wg sync.WaitGroup
	wg.Add(2)
	go discover.StartDiscover(p)
	go listener.StartListener(p, *initParams.Port)
	wg.Wait()
}
