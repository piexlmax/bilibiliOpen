package client

import (
	"bilibiliOpen/content"
	"bilibiliOpen/module"
	"bilibiliOpen/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type GetTokenRes struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
	Data      struct {
		AuthBody string   `json:"auth_body"`
		Host     []string `json:"host"`
		Ip       []string `json:"ip"`
		TcpPort  []int    `json:"tcp_port"`
		WsPort   []int    `json:"ws_port"`
		WssPort  []int    `json:"wss_port"`
	} `json:"data"`
}

type BiliBiliClient struct {
	*websocket.Conn
	RoomId string
	Chan   chan module.Chan
}

var dia websocket.Dialer

func (b *BiliBiliClient) StartWsLinkInfo(wg *sync.WaitGroup) {
	defer wg.Done()
	var client http.Client
	body := `
		{"room_id":` + b.RoomId + `}
		`
	req := utils.NewBiliBiReq("/v1/common/websocketInfo", body)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	err, tokenRes := utils.FmtBody[GetTokenRes](res.Body)
	if err != nil {
		return
	}
	if tokenRes.Code == 0 {
		var link url.URL
		for i := range tokenRes.Data.Host {
			port := strconv.Itoa(tokenRes.Data.WsPort[0])
			link.Host = tokenRes.Data.Host[i] + ":" + port
			link.Scheme = "ws"
			link.Path = "sub"
			b.Conn, _, err = dia.Dial(link.String(), nil)
			if err == nil {
				break
			}
			if err != nil {
				log.Println(err)
			}
		}
		reqProto := utils.MakeProto(tokenRes.Data.AuthBody, content.OP_AUTH)
		e := b.Conn.WriteMessage(websocket.BinaryMessage, reqProto)
		if e != nil {
			log.Println(e)
			fmt.Println("连接失败")
		}
		for {
			_, p, e := b.Conn.ReadMessage()
			if e != nil {
				log.Println(e)
				break
			}
			b.readMsg(p)
		}
	}
}

func (b *BiliBiliClient) pingpang() {
	for {
		time.Sleep(30 * time.Second)
		b.Conn.WriteMessage(websocket.BinaryMessage, utils.MakeProto("", content.OP_HEARTBEAT))
	}
}

func (b *BiliBiliClient) readMsg(byteMsg []byte) {
	switch byteMsg[7] {
	// 普通不需要解码的信息
	case content.COMMONMSG:
		if byteMsg[11] == content.OP_AUTH_REPLY {
			var authRes module.AuthRes
			json.Unmarshal(byteMsg[16:], &authRes)
			if authRes.Code == 0 {
				log.Println("服务注册成功")
				go b.pingpang()
			}
		}
		if byteMsg[11] == content.OP_SEND_SMS_REPLY {
			b.danmu(byteMsg[16:])
			b.Chan <- module.Chan{Msg: byteMsg[16:], Err: nil, MsgType: websocket.BinaryMessage}
		}
		if byteMsg[11] == content.OP_HEARTBEAT_REPLY {
			log.Println("心跳返回")
		}
		break
		// 需要ZLIB解码的消息
	case content.ZLIBMSG:
		err, msg := utils.ReadMsg(byteMsg[16:])
		fmt.Println(err, string(msg))
		break
	}
}

func (b *BiliBiliClient) danmu(jsonB []byte) {
	var dm module.DanmuRes
	err := json.Unmarshal(jsonB, &dm)
	if err != nil {
		fmt.Println(dm)
	}

	fmt.Printf(`
	房间号:%d,
	用户ID:%d,
	用户名:%s,
	大航海等级:%d,
	是不是真爱粉:%t,
	头像:%s,
	说:%s,
	`, dm.Data.RoomId, dm.Data.Uid, dm.Data.Uname, dm.Data.GuardLevel, dm.Data.FansMedalWearingStatus, dm.Data.Uface, dm.Data.Msg)
}
