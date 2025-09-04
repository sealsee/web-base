package httpUtil

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendHttpRequest[T any](method, url string, header map[string]string, body map[string]interface{}) (*T, error) {
	// 创建要发送的数据
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		// 处理错误
		return nil, err
	}

	// 创建请求
	req, err := http.NewRequest(method, url, bytes.NewReader(jsonBytes))
	if err != nil {
		// 处理错误
		return nil, err
	}

	// 设置请求头
	// Set时候，如果原来这一项已存在，后面的就修改已有的
	// Add时候，如果原本不存在，则添加，如果已存在，就不做任何修改
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("tkn", "yx001")
	for k, v := range header {
		req.Header.Set(k, v)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// 处理错误
		return nil, err
	}
	// fmt.Println(resp.Body)
	defer resp.Body.Close()
	var result T
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
