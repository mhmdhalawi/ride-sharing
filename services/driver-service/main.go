package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

var GrpcAddr = ":9092"

func main() {
	svc := NewService()
	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpcserver.NewServer()
	NewGrpcHandler(grpcServer, svc)

	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting gRPC server Trip service on port %s", lis.Addr().String())
		serverErrors <- grpcServer.Serve(lis)
	}()

	select {
	case err := <-serverErrors:
		log.Printf("server error: %v", err)

	case sig := <-shutdown:
		log.Printf("shutting down due to %v", sig)
		grpcServer.GracefulStop()
	}

}
