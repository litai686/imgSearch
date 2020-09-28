package controller

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"imgSearch/framework/db"
	"imgSearch/framework/public"
	"net/http"
)

//需要暴露出去的路由
func init() {
	login()
}

func login() {
	resMap := make(map[string]interface{})
	http.HandleFunc("/login", func(response http.ResponseWriter, request *http.Request) {
		//检测路由
		jsonStr, err := public.PostData(request.Body)
		if err == nil {
			println(jsonStr)
			//JSON对象
			jsonData := gjson.Parse(jsonStr)
			username := jsonData.Get("name").String()
			pwd := jsonData.Get("pwd").String()
			rows, err := db.Query_sql("select id from ipfs_user where name='" + username + "' and pwd='" + pwd + "'")
			if err == nil {
				if rows.Next() {
					resMap["code"] = "ok"
					var id int
					_ = rows.Scan(&id)
					resMap["id"] = id
				} else {
					resMap["code"] = "用户名密码错误"
				}
			} else {
				resMap["code"] = "系统错误"
			}
			resStr, _ := json.Marshal(resMap)
			response.Write(resStr)
		} else {
			response.Write([]byte(err.Error()))
			println(err)
		}
	})
}
