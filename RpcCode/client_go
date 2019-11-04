package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:9090")
	if err != nil {
		panic(err.Error())
	}

	var req float32
	req = 3

	// 同步调用方式
	//var resp *float32
	//err = client.Call("MathUtil.CalculateCircleArea", req, &resp)
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println(*resp)

	var respSync *float32
	// 异步调用方式
	syncCall := client.Go("MathUtil.CalculateCircleArea", req, &respSync, nil)
	replayDone := <-syncCall.Done
	fmt.Println(replayDone)
	fmt.Println(*respSync)
}
