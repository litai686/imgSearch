package webSocketClient

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
)

var isAlive bool = false
var conn *websocket.Conn

var ipHost = ""
var pathUrl = ""

var jsonMap = make(map[string]interface{})

//启动websocket
func Start(ip string, path string) {
	connect(ip, path) //连接
}

//连接服务器
func connect(ip string, path string) {
	ipHost = ip
	pathUrl = path
	u := url.URL{Scheme: "ws", Host: ip, Path: path}
	con, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		println(err.Error())
		println("断线，5秒后重连")
		time.Sleep(5 * time.Second)
		connect(ipHost, pathUrl)
	} else {
		conn = con
		isAlive = true
		println("websocket连接成功")
		jsonMap["id"] = conn.LocalAddr().String()
		jsonMap["route"] = "test"
		jsonMap["data"] = "123"
		bytedata, err := json.Marshal(jsonMap)
		println("要发送的数据:", string(bytedata))
		if err != nil {
			println(err.Error())
		} else {
			SendMsg(string(bytedata)) //发送绑定设备消息
		}
		readMsgThread(conn) //循环读取消息
	}
}

// 读取消息
func readMsgThread(con *websocket.Conn) {
	for {
		_, message, err := con.ReadMessage()
		if err != nil {
			println(err.Error())
			isAlive = false
			break
		}
		println("接收消息：", string(message))
	}
	//重新连接
	println("断线，5秒后重连")
	time.Sleep(5 * time.Second)
	_ = con.Close()
	connect(ipHost, pathUrl)
}

//发送消息
func SendMsg(msg string) bool {
	if isAlive {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err == nil {
			return true
		} else {
			println(err)
			return false
		}
	} else {
		println("掉线")
		return false
	}

}
