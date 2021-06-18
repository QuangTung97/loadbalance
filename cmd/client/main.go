package main

import (
	"balance/balancepb"
	"context"
	"fmt"
	"github.com/QuangTung97/goblin"
	"google.golang.org/grpc"
	"time"
)

func main() {
	pool := goblin.NewPoolClient(goblin.ClientConfig{
		Addresses: []string{
			"localhost:5001",
			"localhost:5002",
			"localhost:5003",
		},
		Options: []grpc.DialOption{
			grpc.WithInsecure(),
		},
	})

	time.Sleep(1 * time.Second)

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
