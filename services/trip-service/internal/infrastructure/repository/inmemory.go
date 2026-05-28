package repository

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"
)

type inMemoryRepository struct {
	trip      map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

func NewInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{
		trip:      make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *inMemoryRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	// Store the trip in the in-memory map
	r.trip[trip.ID.Hex()] = trip
	return trip, nil
}
