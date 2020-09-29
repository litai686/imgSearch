package controller

import (
	"encoding/json"
	"imgSearch/framework/db"
	"net/http"
	"strconv"
)

//注册路由
func init() {
	println("注册路由,优先级高:controller/test")
	show()
}

func show() {
	http.HandleFunc("/img/list", func(response http.ResponseWriter, request *http.Request) {
		//检测路由
		query := request.URL.Query()
		page := query["p"][0]
		//println("page:", page)
		resMap := make(map[string]interface{})
		resMap["code"] = "ok"
		//检测路由
		pageInt, _ := strconv.Atoi(page)
		//println("pageInt:", pageInt)
		startPage := (pageInt - 1) * 10
		sql := "select id,name,createTime from dfs where pid <= (select pid from dfs order by pid desc limit " + strconv.Itoa(startPage) + " ,1) order by pid desc limit 10"
		rows, err := db.Query_sql(sql)
		if err == nil {
			var id string
			var name string
			var createTime string
			var ls_group []interface{}
			{
			}
			i := 0
			for rows.Next() {
				i++
				_ = rows.Scan(&id, &name, &createTime)
				//println(id, name, createTime)
				ls := map[string]interface{}{
					"id":         id,
					"name":       name,
					"createTime": createTime,
				}
				ls_group = append(ls_group, ls)
			}
			//println("i----:", i)
			if i == 0 {
				resMap["code"] = "0"
			}
			resMap["data"] = ls_group
			defer rows.Close()
		} else {
			resMap["code"] = "系统错误"
		}
		resStr, _ := json.Marshal(resMap)
		response.Write(resStr)
	})
}
