package main

import (
	"net"
	"math"
	"net/http"
	"net/rpc"
)

type MathUtil struct {
}

func (mu *MathUtil) CalculateCircleArea(req float32, resp *float32) error {
	*resp = math.Pi * req * req
	return nil
}

func main()  {
	mathUtil := new(MathUtil)

	err := rpc.Register(mathUtil)
	if err != nil {
		panic(err.Error())
	}

	rpc.HandleHTTP()

	listen, err := net.Listen("tcp", ":9090")

	if err != nil{
		panic(err.Error())
	}

	http.Serve(listen, nil)
}
