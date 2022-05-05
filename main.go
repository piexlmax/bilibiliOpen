package main

import (
	"bilibiliOpen/client"
	"bilibiliOpen/module"
	"bilibiliOpen/server"
	"sync"
)

var b = client.BiliBiliClient{
	RoomId: "22310900",
	Chan:   make(chan module.Chan, 10000),
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go b.StartWsLinkInfo(&wg)
	go server.BiliBiliWSServer(&wg, b.Chan)
	wg.Wait()
}
