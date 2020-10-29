package controller

import (
	"context"
	"encoding/json"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/api/apistruct"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/tidwall/gjson"
	"imgSearch/framework/public"
	"log"
	"net/http"
	"strings"
)

func init() {
	walletBalance()
	createWalletAddress()
	sendWallet()
}

//var authToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.I2lMfbqqBzB0_-3lpUskCFkFw0A1j5Z9Rx8IQSMJ7tg"  //本机
var authToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.02xaP2N7s388CjchnQp6zlYnOT402nd0CCVrL5m74Po" //公司
var ipaddr = "127.0.0.1:1234"
var headers = http.Header{"Authorization": []string{"Bearer " + authToken}}
//查询钱包余额
func walletBalance() {
	http.HandleFunc("/walletBalance", func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
		resMap := make(map[string]interface{})
		//检测路由
		//query := request.URL.Query()
		//if len(query) > 0 {
		//	pageindex := query["page"][0]
		//	limit := query["limit"][0]
		//	println("pageindex:", pageindex, " limit:", limit)
		//}

		var walletAddress string
		resMap["code"] = "ok"
		jsonStr, err := public.PostData(request.Body)
		if err == nil {
			println(jsonStr)
			//JSON对象
			jsonData := gjson.Parse(jsonStr)
			token := jsonData.Get("token").String()
			if public.CheckToken(token) == false {
				resMap["code"] = "token错误"
				resStr, _ := json.Marshal(resMap)
				response.Write(resStr)
				return
			}
			walletAddress = jsonData.Get("data").Get("walletAddress").String()
			println("walletAddress:", walletAddress)
		} else {
			resMap["code"] = "参数错误"
		}

		wallet, _ := address.NewFromString(walletAddress) //钱包地址
		var api apistruct.FullNodeStruct
		closer, err := jsonrpc.NewMergeClient(context.Background(), "ws://"+ipaddr+"/rpc/v0", "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
		if err != nil {
			log.Fatalf("connecting with lotus failed: %s", err)
		}
		defer closer()
		filNum, err := api.WalletBalance(context.Background(), wallet)
		if err == nil {
			wallet_res := map[string]string{
				"wallte": big.Div(filNum, big.NewInt(1000000000000000000)).String() + "." + big.Mod(filNum, big.NewInt(1000000000000000000)).String(),
			}
			resMap["data"] = wallet_res

		} else {
			resMap["msg"] = err.Error()
		}
		resStr, _ := json.Marshal(resMap)
		response.Write(resStr)
	})

}

//创建钱包地址
func createWalletAddress() {
	http.HandleFunc("/createWalletAddress", func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
		resMap := make(map[string]interface{})
		resMap["code"] = "ok"
		jsonStr, err := public.PostData(request.Body)
		if err == nil {
			println(jsonStr)
			//JSON对象
			jsonData := gjson.Parse(jsonStr)
			token := jsonData.Get("token").String()
			if public.CheckToken(token) == false {
				resMap["code"] = "token错误"
				resStr, _ := json.Marshal(resMap)
				response.Write(resStr)
				return
			}
		} else {
			resMap["code"] = "参数错误"
		}

		var api apistruct.FullNodeStruct
		closer, err := jsonrpc.NewMergeClient(context.Background(), "ws://"+ipaddr+"/rpc/v0", "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
		if err != nil {
			log.Fatalf("connecting with lotus failed: %s", err)
		}
		defer closer()
		addr, err := api.WalletNew(context.Background(), types.KTBLS)
		if err == nil {
			wallet_res := map[string]interface{}{
				"wallteAddress": addr,
			}
			resMap["data"] = wallet_res

		} else {
			resMap["msg"] = err.Error()
		}
		resStr, _ := json.Marshal(resMap)
		response.Write(resStr)
	})
}

//转账
func sendWallet() {
	http.HandleFunc("/sendWallet", func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
		resMap := make(map[string]interface{})
		var fromAddress = ""
		var sendAddress = ""
		var fil = ""
		resMap["code"] = "ok"
		jsonStr, err := public.PostData(request.Body)
		if err == nil {
			println(jsonStr)
			//JSON对象
			jsonData := gjson.Parse(jsonStr)
			token := jsonData.Get("token").String()
			if public.CheckToken(token) == false {
				resMap["code"] = "token错误"
				resStr, _ := json.Marshal(resMap)
				response.Write(resStr)
				return
			}
			fromAddress = jsonData.Get("data").Get("fromAddress").String()
			sendAddress = jsonData.Get("data").Get("sendAddress").String()
			fil = jsonData.Get("data").Get("fil").String()
			println("fromAddress:", fromAddress)
			println("sendAddress:", sendAddress)

		} else {
			resMap["code"] = "参数错误"
		}

		//var api apistruct.FullNodeStruct
		//closer, err := jsonrpc.NewMergeClient(context.Background(), "ws://"+ipaddr+"/rpc/v0", "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
		//defer closer()
		//if err != nil {
		//	resMap["msg"] = err.Error()
		//	resStr, _ := json.Marshal(resMap)
		//	response.Write(resStr)
		//	return
		//}

		if fromAddress != "" && sendAddress != "" {
			res := public.ExecShell("lotus send --from " + fromAddress + " " + sendAddress + " " + fil)
			res = strings.Replace(res, "\n", "", -1) //去掉尾部换行符
			println("msgID:" + res)
			println("len:", len(res))
			if len(res) == 62 {
				resMap["data"] = map[string]string{
					"msgID": res,
				}
			} else {
				resMap["code"] = "地址参数错误"
			}

		} else {
			resMap["code"] = "地址错误"
		}
		resStr, _ := json.Marshal(resMap)
		response.Write(resStr)
		//addr, err := api.Wallets
		//if err == nil {
		//	wallet_res := map[string]interface{}{
		//		"wallteAddress": addr,
		//	}
		//	resMap["data"] = wallet_res
		//
		//} else {
		//	resMap["msg"] = err.Error()
		//}
		//resStr, _ := json.Marshal(resMap)
		//response.Write(resStr)
	})
}
