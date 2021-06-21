package main

import (
	"balance/balancepb"
	"context"
	"flag"
	"fmt"
	"github.com/QuangTung97/goblin"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var port = flag.Int("port", 5001, "port of grpc")

var nodes = flag.String("nodes", "", "the nodes in cluster")

func getAddresses() []string {
	if len(*nodes) == 0 {
		return nil
	}
	result := strings.Split(*nodes, ",")
	for i := range result {
		result[i] = strings.TrimSpace(result[i])
	}
	return result
}

func main() {
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	pool, err := goblin.NewPoolServer(goblin.ServerConfig{
		GRPCPort:     uint16(*port),
		IsDynamicIPs: true,
		ServiceAddr:  "goblin-server:5001",
		DialOptions:  []grpc.DialOption{grpc.WithInsecure()},
	}, goblin.WithServerLogger(logger))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger),
		),
	)
	pool.Register(grpcServer)
	balancepb.RegisterBalanceServiceServer(grpcServer, &balanceServer{
		pool: pool,
	})

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			panic(err)
		}

		fmt.Println("Start gRPC")
		err = grpcServer.Serve(listener)
		if err != nil {
			panic(err)
		}
		fmt.Println("Shutdown successfully")
	}()

	<-exit
	err = pool.Shutdown()
	if err != nil {
		panic(err)
	}
	grpcServer.GracefulStop()
	wg.Wait()
}

type balanceServer struct {
	balancepb.UnimplementedBalanceServiceServer
	pool *goblin.PoolServer
}

func (s balanceServer) Hello(ctx context.Context, req *balancepb.HelloRequest) (*balancepb.HelloResponse, error) {
	return &balancepb.HelloResponse{
		Msg: fmt.Sprintf("%s %s", s.pool.GetName(), s.pool.GetMemberlistAddress()),
	}, nil
}
