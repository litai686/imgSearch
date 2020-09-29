package main

import (
	"fmt"
	"github.com/aWildProgrammer/goticker"
	"gopkg.in/ini.v1"
	"imgSearch/framework/db"
	"imgSearch/framework/httpServer"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	task := goticker.New(100, false)
	_ = task.AddTaskCallBackFunc(fileCoinUpFile, 10, "cli") // 每间隔3秒执行一次 test 函数

	httpServer.Start()

}

var state = false // 任务执行状态
//filecoin文件上传
func fileCoinUpFile(param interface{}) {
	//println(param.(string))
	if state == false {
		state = true
		//读取配置文件
		cfg, err := ini.Load("./config.ini")
		if err != nil {
			fmt.Println("文件读取错误", err)
			state = false
			return
		}
		minerIDs := cfg.Section("").Key("minerID")
		rows, err := db.Query_sql("select name,id,info,cid from dfs where info='0' limit 1")
		if err == nil {
			var name string
			var id string
			var info string
			var cid string
			if rows.Next() {
				_ = rows.Scan(&name, &id, &info, &cid)
				rows.Close()
				if info == "1" {
					println("文件已经上传过：", name)
				} else {
					println("上传文件：", name)
					dir, err := filepath.Abs(name)
					if err != nil {
						println("文件不存在：", err)
						state = false
						return
					}
					println("上传文件到filecoin：", dir)
					res := execShell("lotus client import " + dir)
					res = strings.Replace(res, "\n", "", -1) //去掉尾部换行符
					println("shell res :", res)
					//println("--:", strings.Index(res, "Import"))
					if strings.Index(res, "Import") == -1 {
						println("上传出错：" + dir)
						state = false
						return
					}
					resSplit := strings.Split(res, " ")
					if len(resSplit) > 2 {
						println("cid: ", resSplit[3])
						cid = resSplit[3]
					}
				}

				if cid == "0" {
					println("上传出错，获取不到CID")
					state = false
					return
				}

				//更新上传状态
				_, err := db.Update_sql("update dfs set info='1' where id='" + id + "'")
				//更新CID
				RowsAffected, err := db.Update_sql("update dfs set cid='" + cid + "' where id='" + id + "'")
				if err == nil {
					if RowsAffected > 0 {
						println("上传到filcoin成功: ", cid)
						println("开始交易: ", cid)
						println("minerIDs:", minerIDs.Value())
						minerIDSplit := strings.Split(minerIDs.Value(), ",")
						for _, val := range minerIDSplit {
							if strings.Replace(val, " ", "", -1) != "" {
								println("执行脚本：", "lotus client deal "+cid+" "+val+" 0.000000005 518400")
								res_deal := execShell("lotus client deal " + cid + " " + val + " 0.000000005 518400")
								res_deal = strings.Replace(res_deal, "\n", "", -1) //去掉尾部换行符
								if len(res_deal) == 59 {
									println("交易成功")
								} else {
									println("交易失败:", res_deal)
								}
							}
						}

					}
				} else {
					println("上传到filcoin失败: ", cid)
				}

			} else {
				//println("暂无新数据")
			}
		} else {
			println(err.Error())
		}
	}

	state = false
}

func execShell(strCommand string) string {
	cmd := exec.Command("/bin/bash", "-c", strCommand)

	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		fmt.Println("Execute failed when Start:" + err.Error())
		return ""
	}

	out_bytes, _ := ioutil.ReadAll(stdout)
	defer stdout.Close()

	if err := cmd.Wait(); err != nil {
		fmt.Println("Execute failed when Wait:" + err.Error())
		return ""
	}
	return string(out_bytes)
}
