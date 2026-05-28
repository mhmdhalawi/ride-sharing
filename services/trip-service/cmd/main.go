package main

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"time"
)

func main() {
	ctx := context.Background()
	// Initialize the in-memory repository
	repo := repository.NewInMemoryRepository()
	// Create the trip service
	svc := service.NewTripService(repo)

	fare := &domain.RideFareModel{
		UserID:            "user123",
		TotalPriceInCents: 1500,
	}

	t, err := svc.CreateTrip(ctx, fare)
	if err != nil {
		log.Println(err)
	}
	log.Println("Created trip with ID:", t.ID.Hex())
	// remove later, just to keep the app runnning

	for {
		time.Sleep(2 * time.Second)
	}
}
