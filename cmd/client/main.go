package main

import (
	"balance/balancepb"
	"context"
	"fmt"
	"github.com/QuangTung97/goblin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	pool := goblin.NewPoolClient(goblin.ClientConfig{
		Addresses: []string{
			"goblin-server:5001",
		},
		Options: []grpc.DialOption{
			grpc.WithInsecure(),
		},
	}, goblin.WithClientLogger(logger))

	for !pool.Ready() {
		time.Sleep(100 * time.Millisecond)
	}

	for {
		err := pool.GetConn(func(conn *grpc.ClientConn) error {
			client := balancepb.NewBalanceServiceClient(conn)
			resp, err := client.Hello(context.Background(), &balancepb.HelloRequest{
				Name: "Ta Quang Tung",
			})
			if err != nil {
				panic(err)
			}
			fmt.Println(resp)
			return nil
		})
		if err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)
	}
}
