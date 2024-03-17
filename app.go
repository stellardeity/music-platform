/*
Go training project

	-name string
	  	you name (default "eku@eku-HP-ProBook-450-G3")
	-peers string
	  	Path to file with peer addresses on each line (default "peers.txt")
	-port int
	  	port that have to listen (default 35035)
	-webview
	  	Start WebView ui (default false)
*/
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"os/user"
	"sync"
	"syscall"

	"rockwall/discover"
	"rockwall/listener"
	"rockwall/proto"
)

// InitParams initializing params
type InitParams struct {
	Name         *string
	Port         *int
	PeersFile    *string
	StartWebView *bool
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

	flag.Parse()

	//flag.PrintDefaults()

	// Настройки логирования
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

}

func main() {

	fff := retro()

	i := fff(1, 2)
	log.Printf("result = %v", i)

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		log.Printf("Exit by signal: %s", sig)
		os.Exit(1)
	}()

	p := proto.NewProto(*initParams.Name, *initParams.Port)

	startWithoutWebView(p)
}

func retro() func(a int, b int) int {
	return func(a int, b int) int {
		return a + b
	}
}

// startWithoutWebView Запуск приложения без запуска WebView
func startWithoutWebView(p *proto.Proto) {
	var wg sync.WaitGroup
	wg.Add(2)
	go discover.StartDiscover(p, *initParams.PeersFile)
	go listener.StartListener(p, *initParams.Port)
	wg.Wait()
}
