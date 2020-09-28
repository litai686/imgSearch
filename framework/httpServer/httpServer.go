package httpServer

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"imgSearch/controller"
	"net/http"
	"reflect"
)

func Start() {
	//注册路由
	controller.Init()
	//region 静态服务器
	http.Handle("/", http.FileServer(http.Dir("./static")))
	//endregion

	//region websocket
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	var clientFd = 0 //客户端描述符
	var clientMap = make(map[string]*websocket.Conn)
	http.HandleFunc("/ws", func(response http.ResponseWriter, request *http.Request) {
		var (
			conn *websocket.Conn
			err  error
			data []byte
		)
		if conn, err = upgrader.Upgrade(response, request, nil); err != nil {
			return
		}
		for {
			_, data, err = conn.ReadMessage()
			if err != nil {
				//客户端掉线处理
				for key := range clientMap {
					if clientMap[key] == conn {
						delete(clientMap, key)
					}
				}
				println("客户端剩余个数：", len(clientMap))
				println(err.Error())
				break
			}
			//接收信息
			fmt.Println(string(data))
			resJosn := gjson.Parse(string(data))
			id := resJosn.Get("id").String()
			if id == "" {
				println("非法客户端访问")
				clientMap[string(clientFd)] = conn
				clientFd++
				//break
			} else {
				//绑定客户端ID
				clientMap[id] = conn
			}
			println(conn)
			if err = conn.WriteMessage(websocket.TextMessage, []byte(resJosn.String())); err != nil {
				break
			}
		}
		_ = conn.Close()

	})
	//endregion
	_ = http.ListenAndServe("0.0.0.0:8001", nil)
}

func Call(m map[string]interface{}, name string, params ...interface{}) ([]reflect.Value, error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		return nil, errors.New("the number of input params not match!")
	}
	in := make([]reflect.Value, len(params))
	for k, v := range params {
		in[k] = reflect.ValueOf(v)
	}
	return f.Call(in), nil
}

//func Apply(f interface{}, args []interface{}) ([]reflect.Value) {
//	fun := reflect.ValueOf(f)
//	in := make([]reflect.Value, len(args))
//	for k, param := range args {
//		in[k] = reflect.ValueOf(param)
//	}
//	r := fun.Call(in)
//	return r
//}
