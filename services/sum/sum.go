package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"net"
	"os"
	etcd "service-collection/components"
	proto "service-collection/services/sum/proto"
)

var (
	//The return value is the address of a string variable that stores the value of the flag.
	ser  = flag.String("service", "sum_service", "service name")
	host = flag.String("host", "localhost", "listening host")
	port = flag.String("port", "50001", "listening port")
	reg  = flag.String("reg", "http://localhost:2379", "register etcd address")
)

func (s *server) GetSum(ctx context.Context, in *proto.SumRequest) (*proto.SumResponse, error) {
	return &proto.SumResponse{Output: in.Input}, nil
}

func main() {
	flag.Parse()

	listen, err := net.Listen("tcp", net.JoinHostPort(*host, *port))
	if err != nil {
		panic(err)
	}

	err = etcd.Register(*reg, *ser, *host, *port, 15)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)

	go func() {
		<-ch
		etcd.UnRegister()
		os.Exit(1)
	}()

	s := grpc.NewServer()
	proto.RegisterSumServiceServer(s, &server{})
	s.Serve(listen)
}

type server struct {
}
