package server

import (
	"bilibiliOpen/module"
	"bilibiliOpen/utils"
	"encoding/json"
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

type BConn struct {
	Conn *websocket.Conn
	Name string
	C    chan module.Chan
}

var Conns = make(map[string]BConn)

func ChanHandler(c <-chan module.Chan) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Conn, err := UP.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err, "web")
			return
		}
		var bconn BConn
		bconn.Conn = Conn
		bconn.C = make(chan module.Chan, 10)
		go ReadConn(bconn)
		// 客户端连接中断未断开
		for {
			for s := range Conns {
				ch := <-Conns[s].C
				err := Conns[s].Conn.WriteMessage(ch.MsgType, ch.Msg)
				if err != nil {
					log.Println(err)
					Conns[s].Conn.Close()
					delete(Conns, Conns[s].Name)
					fmt.Println(Conns[s].Name, "已经断开连接")
				}
			}

		}
	}
}

func BiliBiliWSServer(wg *sync.WaitGroup, c chan module.Chan) {
	defer wg.Done()
	http.HandleFunc("/", ChanHandler(c))
	go func() {
		for {
			ch := <-c
			for s := range Conns {
				Conns[s].C <- ch
			}
		}
	}()
	fmt.Println(http.ListenAndServe(":8888", nil))
}

func ReadConn(bconn BConn) {
	for {
		_, m, err := bconn.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			log.Println("读取消息失败")
			break
		}
		var c module.WebCMD
		json.Unmarshal(m, &c)
		HandleCMD(c, bconn)
	}
}

func HandleCMD(c module.WebCMD, bconn BConn) {
	switch c.CMD {
	case "register":
		err, registerRes := utils.FmtStrToStruct[module.RegisterData](c.Data)
		if err != nil {
			log.Println(err)
			return
		}
		bconn.Name = registerRes.Name
		Conns[registerRes.Name] = bconn
		log.Println(registerRes.Name, "注册成功")
	}
}
