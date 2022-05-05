package server

import (
	"bilibiliOpen/module"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var UP = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ChanHandler(c <-chan module.Chan) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Conn, err := UP.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err, "web")
			return
		}
		fmt.Println("web接入")
		// 客户端连接中断未断开
		for {
			ch := <-c
			fmt.Println(ch)
			err := Conn.WriteMessage(ch.MsgType, ch.Msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func BiliBiliWSServer(wg *sync.WaitGroup, c chan module.Chan) {
	defer wg.Done()
	http.HandleFunc("/", ChanHandler(c))
	fmt.Println(http.ListenAndServe(":8888", nil))
}
