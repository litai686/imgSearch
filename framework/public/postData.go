package public

import (
	"bytes"
	"io"
	"io/ioutil"
)

//获取POST数据
func PostData(data io.ReadCloser) (string, error) {
	postResult, err := ioutil.ReadAll(data)
	if err == nil {
		//JSON字符串
		jsonStr := bytes.NewBuffer(postResult).String()
		return jsonStr, nil
	} else {
		return "", err
	}
}
