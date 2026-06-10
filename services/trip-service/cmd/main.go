package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

var GrpcAddr = ":9093"

func main() {
	inmemRepo := repository.NewInMemoryRepository()
	svc := service.NewTripService(inmemRepo)

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpcserver.NewServer()

	grpc.NewGRPCHandler(grpcServer, svc)

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
