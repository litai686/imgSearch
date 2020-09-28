package controller

import (
	"encoding/json"
	"imgSearch/framework/db"
	"net/http"
)

func init() {
	device_list()
}

func device_list() {

	http.HandleFunc("/device_list", func(response http.ResponseWriter, request *http.Request) {
		resMap := make(map[string]interface{})
		//检测路由
		query := request.URL.Query()
		if len(query) > 0 {
			pageindex := query["page"][0]
			limit := query["limit"][0]
			println("pageindex:", pageindex, " limit:", limit)
		}

		//jsonStr, err := public.PostData(request.Body)
		//if err == nil {
		//	println(jsonStr)
		//	////JSON对象
		//	//jsonData := gjson.Parse(jsonStr)
		//	//username := jsonData.Get("name").String()
		//	//pwd := jsonData.Get("pwd").String()
		//
		//}
		resMap["code"] = "ok"
		resMap["code"] = 0
		resMap["msg"] = ""

		ls_group := []interface{}{}
		rows, err := db.Query_sql("select id,device_name,device_local_ip,device_state,device_soft_version,device_start_time,device_remark,device_type from ipfs_device where device_type=0")
		if err == nil {
			i := 0
			for rows.Next() {
				var id int
				var device_name string
				var device_local_ip string
				var device_state string
				var device_soft_version string
				var device_start_time string
				var device_remark string
				var device_type string
				_ = rows.Scan(&id, &device_name, &device_local_ip, &device_state, &device_soft_version, &device_start_time, &device_remark, &device_type)
				ls := map[string]interface{}{
					"id":                  id,
					"device_name":         device_name,
					"device_local_ip":     device_local_ip,
					"device_state":        device_state,
					"device_soft_version": device_soft_version,
					"device_start_time":   device_start_time,
					"device_remark":       device_remark,
					"device_type":         device_type,
				}
				ls_group = append(ls_group, ls)
				i++
			}
			resMap["count"] = i
			resMap["data"] = ls_group

		} else {
			resMap["code"] = "接口出错"
		}
		resStr, _ := json.Marshal(resMap)
		println((string(resStr)))
		response.Write(resStr)
	})
}
