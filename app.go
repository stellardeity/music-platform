package main

import (
	"os"
	"stellard/discover"
	"stellard/listener"
	"stellard/proto"
	"sync"
)

func init() {
	if len(os.Args) != 2 {
		panic("len args != 2")
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	node := proto.NewNode(os.Args[1])
	go listener.StartListener(node)
	go discover.StartDiscover(node)
	wg.Wait()
}
