package main

import (
	"fmt"
	EndPoint "kitrpcclient/endpoint"
	"kitrpcclient/tool"
	"kitrpcclient/transport"
)

func main() {
	// data := map[string]string{"key": "lipper", "val": "30"}
	// jsonData, err := json.Marshal(data)
	// prep, err := http.NewRequest("POST", "http://127.0.0.1:8000/set", bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// prep.Header.Set("Content-Type", "application/json")
	// client := &http.Client{}
	// presp, err := client.Do(prep)
	// var res response

	// if err = json.NewDecoder(presp.Body).Decode(&res); err != nil {
	// 	log.Fatalf("json.NewDecoder: %v", err)
	// }
	// log.Println(res)
	// i, err := tool.Direct("Post", "http://127.0.0.1:8000", transport.EncodeRequestFunc, transport.DecodeResponseFunc, EndPoint.PostRequest{Key: "lipper", Val: "30"})

	i, err := tool.ServiceDiscovery("Post", "http://127.0.0.1:8500", transport.EncodeRequestFunc, transport.DecodeResponseFunc, EndPoint.PostRequest{Key: "lipper", Val: "30"}, "测试1", true, "test")
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	res, ok := i.(EndPoint.Response)
	if !ok {
		fmt.Println("no ok")
		return
	}
	fmt.Println(res)
}

type response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
