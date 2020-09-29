package controller

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"imgSearch/framework/db"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

//需要暴露出去的路由
func init() {
	println("注册路由,优先级高:controller/upload")
	upload()
	img()
}

//图片上传
func upload() {
	resMap := make(map[string]interface{})
	http.HandleFunc("/dfs/upload", func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域

		//文件上传目录
		uploadPath := "./upload/" + time.Now().Format("20060102") + "/"
		_ = request.ParseMultipartForm(32 << 20)
		file, handler, err := request.FormFile("file")
		defer file.Close()
		if handler.Size > 1024*1024*20 {
			resMap["code"] = "上传失败，图片尺寸太大"
			resStr, _ := json.Marshal(resMap)
			response.Write(resStr)
			return
		}
		if err != nil {
			resMap["code"] = "上传失败 10001"
			resStr, _ := json.Marshal(resMap)
			response.Write(resStr)
			println("错误1", err.Error())
		} else {
			hash := sha256.New()
			if _, err := io.Copy(hash, file); err != nil {
				println("计算sha256错误")
				resMap["code"] = "上传失败 10002"
				resStr, _ := json.Marshal(resMap)
				response.Write(resStr)
			}
			sha256 := fmt.Sprintf("%x", hash.Sum(nil))
			println("文件sha256:" + sha256)

			isMakeDir := makeDir(uploadPath) //创建文件夹
			if isMakeDir {
				//获取文件后缀
				fileType := path.Ext(handler.Filename)
				//重新命名文件
				curTime := strconv.Itoa(int(time.Now().UnixNano()))
				println("当前时间戳：", curTime)
				filename := curTime + fileType //单位秒 time.Now().UnixNano() //单位纳秒
				println("文件名称:", filename)
				f, err := os.OpenFile(uploadPath+filename, os.O_WRONLY|os.O_CREATE, 0666)
				defer f.Close()
				if err != nil {
					println(err.Error())
					resMap["code"] = "上传失败"
				} else {
					file2, _, err := request.FormFile("file")
					defer file2.Close()
					fileSize, err := io.Copy(f, file2)
					if err == nil {

						//存数据库
						insertId, err := db.Insert_sql("insert into dfs (id,name,createTime) values ('" + sha256 + "','" + uploadPath + filename + "','" + time.Now().Format("2006-01-02 15:04:05") + "')")
						if err == nil && insertId > 0 {
							resMap["code"] = "ok"
							resMap["data"] = fileSize
						} else {
							resMap["code"] = "插入数据失败"
						}
					} else {
						resMap["code"] = "上传失败"
					}
				}
				resStr, _ := json.Marshal(resMap)
				response.Write(resStr)
			}

		}

	})
}

func img() {
	http.HandleFunc("/dfs/img", func(response http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		n := query["n"][0]
		//println("n:", n)
		//response.Write([]byte(jsonStr))
		//查找数据库
		rows, err := db.Query_sql("select name from dfs where id='" + n + "'")
		if err == nil {
			if rows.Next() {
				defer rows.Close()
				var name string
				_ = rows.Scan(&name)
				//println("name:", name)
				file, err := os.Open(name)
				defer file.Close()
				if err == nil {
					buff, err := ioutil.ReadAll(file)
					if err == nil {
						response.Write(buff)
					} else {
						response.Write([]byte("error1"))
					}
				} else {
					response.Write([]byte("error2"))
				}
			} else {
				response.Write([]byte("404"))
			}
		} else {
			response.Write([]byte("error db"))
		}
	})

}

//检查目录是否存在
func makeDir(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		//创建文件夹
		err := os.Mkdir(filename, os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹失败", err)
			exist = false
		}
	}
	return exist
}

func GetMd5FromFile(path string) (string, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil

}

//文件的 md5 和 SAH-256 的校验
func getSHA256FromFile(path string) (string, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return "", err
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum, nil
}

//字符串 md5 和 SHA-256 的校验
func getMd5FromString(data string) string {
	h := md5.New()
	io.WriteString(h, data)
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum
}

func getSHA256FromString(data string) string {
	h := sha256.New()
	io.WriteString(h, data)
	sum := fmt.Sprintf("%x", data)
	return sum
}
