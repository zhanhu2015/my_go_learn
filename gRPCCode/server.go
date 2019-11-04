package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"my_go_learn/gRPCCode/message"
	"net"
	"time"
)

type OrderServiceImpl struct {
}

//type OrderServiceServer interface {
//	GetOrderInfo(context.Context, *OrderRequest) (*OrderInfo, error)
//}

// 这个函数实现了message.pd.go中的OrderServiceServer接口里的GetOrderInfo方法
func (os *OrderServiceImpl) GetOrderInfo(ctx context.Context, request *message.OrderRequest) (*message.OrderInfo, error) {
	orderMap := map[string]message.OrderInfo{
		"201907300001": message.OrderInfo{OrderId: "201907300001", OrderName: "衣服", OrdreStatus: "已付款"},
		"201907300002": message.OrderInfo{OrderId: "201907300002", OrderName: "零食", OrdreStatus: "已付款"},
		"201907300003": message.OrderInfo{OrderId: "201907300003", OrderName: "食品", OrdreStatus: "未付款"},
	}
	var response *message.OrderInfo
	current := time.Now().Unix()
	if (request.TimeStamp > current) {
		*response = message.OrderInfo{OrderId: "0", OrderName: "", OrdreStatus: "订单信息异常"}
	} else {
		result := orderMap[request.OrderId]
		if result.OrderId != "" {
			fmt.Println(result)
			return &result, nil
		} else {
			return nil, errors.New("server error")
		}
	}
	return response, nil
}

func main() {
	server := grpc.NewServer()
	message.RegisterOrderServiceServer(server, new(OrderServiceImpl))

	lis, err := net.Listen("tcp", ":9090")

	if err != nil {
		panic(err.Error())
	}
	server.Serve(lis)
}
