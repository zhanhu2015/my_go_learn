package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"my_go_learn/gRPCCode/message"
	"time"
)

func main() {
	// 1, Dial 连接
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		panic(err.Error())
	}

	defer conn.Close()
	orderServiceClient := message.NewOrderServiceClient(conn)

	OrderRequest := &message.OrderRequest{OrderId: "201907300002", TimeStamp: time.Now().Unix()}

	orderInfo, err := orderServiceClient.GetOrderInfo(context.Background(), OrderRequest)

	if orderInfo != nil {
		fmt.Println(orderInfo)
	}
}
