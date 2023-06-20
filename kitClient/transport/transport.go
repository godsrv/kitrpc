package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	EndPoint "kitrpcclient/endpoint"
	"net/http"
	"strconv"
)

func EncodeRequestFunc(c context.Context, request *http.Request, r interface{}) error {
	req, ok := r.(EndPoint.PostRequest)
	if !ok {
		return errors.New("断言失败")
	}
	request.Header.Set("Content-Type", "application/json")
	jsonData, _ := json.Marshal(req)
	// 拿到自定义的请求对象对url做业务处理
	request.URL.Path += "/set"
	request.Body = ioutil.NopCloser(bytes.NewBuffer(jsonData))
	return nil
}

func DecodeResponseFunc(c context.Context, res *http.Response) (response interface{}, err error) {
	// 判断响应
	if res.StatusCode != 200 {
		return nil, errors.New("异常的响应码" + strconv.Itoa(res.StatusCode))
	}
	// body中的内容需要我们解析成我们通用定义好的内容
	var r EndPoint.Response
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
