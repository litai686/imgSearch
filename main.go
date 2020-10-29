package main

import (
	"crypto/md5"
	"fmt"
	"github.com/skip2/go-qrcode"
	"gopkg.in/ini.v1"
	"imgSearch/framework/db"
	"imgSearch/framework/httpServer"
	"imgSearch/framework/public"
	"io"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	//task := goticker.New(100, false)
	//_ = task.AddTaskCallBackFunc(fileCoinUpFile, 10, "cli") // 每间隔3秒执行一次 fileCoinUpFile 函数
	//_ = task.AddTaskCallBackFunc(qrRandom, 60*5, "cli") // 每间隔8秒执行一次 qrRandom
	//_ = task.AddTaskCallBackFunc(doImg, 60*1, "cli")
	httpServer.Start()
}

var mutex sync.Mutex
//filecoin文件上传
func fileCoinUpFile(param interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	//读取配置文件
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		fmt.Println("文件读取错误", err)
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
					return
				}
				println("上传文件到filecoin：", dir)
				res := public.ExecShell("lotus client import " + dir)
				res = strings.Replace(res, "\n", "", -1) //去掉尾部换行符
				println("shell res :", res)
				//println("--:", strings.Index(res, "Import"))
				if strings.Index(res, "Import") == -1 {
					println("上传出错：" + dir)
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
							println("执行脚本：", "lotus client deal "+cid+" "+val+" 1afil 1036800")
							res_deal := public.ExecShell("lotus client deal " + cid + " " + val + " 1afil 1036800")
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
	//matex.Unlock() //解锁
}

var mutex_doImg sync.Mutex

func doImg(param interface{}) {
	mutex_doImg.Lock()
	defer mutex_doImg.Unlock()
	//读取配置文件
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		fmt.Println("文件读取错误", err)
		return
	}
	minerIDs := cfg.Section("").Key("minerID") //minerID集合

	//读取数据库文件
	rows, err := db.Query_sql("select name,id from img_tab where state=0 limit 1")
	var name string
	var id int64
	var dir = ""
	if rows.Next() {
		_ = rows.Scan(&name, &id)
		rows.Close()
		dir = name
		//lsimg := strings.Replace(dir, ".jpg", ".png", 1)
		//_ = os.Remove(lsimg)
	}
	//var FileInfo []os.FileInfo
	//relativePath := "./upload/tmpimg/"
	//if FileInfo, err = ioutil.ReadDir(relativePath); err != nil {
	//	fmt.Println("读取 img 文件夹出错")
	//	return
	//}
	//
	//var dir = ""
	//
	//for _, fileInfo := range FileInfo {
	//	fmt.Println("fileInfoName:", fileInfo.Name())
	//	fmt.Println("fileInfoSize:", fileInfo.Size())
	//	if fileInfo.Size() > (1024 * 1024 * 15) {
	//		dir = relativePath + fileInfo.Name()
	//		lsimg := strings.Replace(dir, ".jpg", ".png", 1)
	//		println("删除：", lsimg)
	//		_ = os.Remove(lsimg)
	//		break
	//	}
	//}

	if dir == "" {
		//println("没有新文件")
		return
	}

	//println("上传文件到filecoin：", dir)
	res := public.ExecShell("lotus client import " + dir)
	res = strings.Replace(res, "\n", "", -1) //去掉尾部换行符
	println("shell res :", res)
	//println("--:", strings.Index(res, "Import"))
	if strings.Index(res, "Import") == -1 {
		println("上传出错：" + dir)
		return
	}
	var cid = "0"
	resSplit := strings.Split(res, " ")
	if len(resSplit) > 2 {
		//println("cid: ", resSplit[3])
		cid = resSplit[3]
	}

	if cid == "0" {
		println("上传出错，获取不到CID")
		return
	}

	println("上传到filcoin成功: ", cid)
	//更新数据库
	update_id, err := db.Update_sql("update img_tab set state=1,cid='" + cid + "' where id=" + strconv.FormatInt(id, 10) + "")
	if update_id > 0 {
		//println("更新数据库成功:" + dir)
	}
	println("开始交易: ", cid)
	//println("minerIDs:", minerIDs.Value())
	minerIDSplit := strings.Split(minerIDs.Value(), ",")
	for _, val := range minerIDSplit {
		if strings.Replace(val, " ", "", -1) != "" {
			//println("执行脚本：", "lotus client deal "+cid+" "+val+" 1afil 1050000")
			println("执行脚本：", "lotus client deal "+cid+" "+val+" 0.0000000001FIL 1040000")
			res_deal := public.ExecShell("lotus client deal " + cid + " " + val + " 0.0000000001 1040000")
			res_deal = strings.Replace(res_deal, "\n", "", -1) //去掉尾部换行符
			if len(res_deal) == 59 {
				println("交易成功")
			} else {
				println("交易失败:", res_deal)
			}
		}
	}
}

var mutex_qrRandom sync.Mutex
//生成随机图片
func qrRandom(param interface{}) {
	mutex_qrRandom.Lock()
	defer mutex_qrRandom.Unlock()
	strRandom := getRandomString(256)
	strRandomMd5 := md5str(strRandom)
	//println("成功的随机数:", strRandomMd5)
	imgName := "./upload/tmpimg/" + strRandomMd5 + ".png"
	//生成二维码
	qrcode.WriteFile(strRandom, qrcode.Medium, 100, imgName)
	time.Sleep(time.Second * 3)
	//合并文件
	println("开始合并文件")
	public.ExecShell("cat  ./upload/tmp/tmp.iso >> " + imgName + "")
	println("合并文件成功", imgName)

	////合成照片
	//tmpimgB, _ := os.Open("./upload/tmp/tmp.jpg")
	//tmpimg, _ := jpeg.Decode(tmpimgB)
	//defer tmpimgB.Close()
	//
	//hcImgB, _ := os.Open(imgName)
	//hcImg, _ := png.Decode(hcImgB)
	//defer hcImgB.Close()
	//
	//offset := image.Pt(0, 0)
	//b := tmpimg.Bounds()
	//m := image.NewRGBA(b)
	//draw.Draw(m, b, tmpimg, image.ZP, draw.Src)
	//draw.Draw(m, hcImg.Bounds().Add(offset), hcImg, image.ZP, draw.Over)
	//imgNewName := "./upload/tmpimg/" + strRandomMd5 + ".jpg"
	//imgw, _ := os.Create(imgNewName)
	//jpeg.Encode(imgw, m, &jpeg.Options{100})
	//defer imgw.Close()
	//把图片信息写入数据库
	insert_id, _ := db.Insert_sql("insert into img_tab (name) values ('" + imgName + "')")
	if insert_id > 0 {
		println("插入数据库成功：" + imgName)
	} else {
		println("插入数据库失败：" + imgName)
	}
}

//字符串 md5 和 SHA-256 的校验
func md5str(data string) string {
	h := md5.New()
	io.WriteString(h, data)
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum
}

//生成随机数
func getRandomString(l int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func ss() {

}
