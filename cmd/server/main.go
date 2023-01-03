package main

import (
	"context"
	"fmt"
	"log"
	"net" // unix Interface
	"os"
	"os/signal"

	hellopb "my_grpc/pkg/my_grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type myServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func (s *myServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer() // like HTTP server

	// use hello method
	hellopb.RegisterGreetingServiceServer(s, NewMyServer())

	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// shutdown
	quit := make(chan os.Signal, 1) // channel, signal
	signal.Notify(quit, os.Interrupt)
	<-quit // only reception
	log.Panicln("stopping gRPC server...")
	s.GracefulStop()
}
