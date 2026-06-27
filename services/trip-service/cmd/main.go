package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"ride-sharing/services/trip-service/internal/infrastructure/events"
	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

var GrpcAddr = ":9093"

func main() {
	rabbitMqURI := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")
	inmemRepo := repository.NewInMemoryRepository()
	svc := service.NewTripService(inmemRepo)

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// RabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQ(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	log.Println("Starting RabbitMQ connection")

	publisher := events.NewTripEventPublisher(rabbitmq)

	grpcServer := grpcserver.NewServer()

	grpc.NewGRPCHandler(grpcServer, svc, publisher)

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
