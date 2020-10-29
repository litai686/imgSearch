package public

import (
	"strconv"
	"time"
)

//是否启用token
var IsToken = false

func CheckToken(token string) bool {
	if IsToken {
		chk := false
		key := MD5("158485188")
		key1 := Substr(token, 0, 32)
		println("token_key_md5:", key)
		if key != key1 {
			return false
		}
		//还原时间戳
		splitStr := Substr(token, 32, 10)
		println("token_timespan:", splitStr)
		//字符串转int64时间戳
		timespan, err := strconv.ParseInt(splitStr, 10, 64)
		if err == nil {
			if (time.Now().Unix()-timespan) <= 120 && (time.Now().Unix()-timespan) >= 0 {
				chk = true
			}
		}
		return chk
	} else {
		return true
	}

}
