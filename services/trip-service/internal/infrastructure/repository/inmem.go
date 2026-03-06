package repository

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
)

type InMemoryTripRepository struct {
	trips     map[string]*domain.TripModel
	rideFares map[string]*domain.RideFareModel
}

func NewInMemoryTripRepository() *InMemoryTripRepository {
	return &InMemoryTripRepository{
		trips:     make(map[string]*domain.TripModel),
		rideFares: make(map[string]*domain.RideFareModel),
	}
}

func (r *InMemoryTripRepository) CreateTrip(ctx context.Context, trip *domain.TripModel) (*domain.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip
	return trip, nil
}

func (r *InMemoryTripRepository) SaveRideFare(ctx context.Context, fare *domain.RideFareModel) error {
	r.rideFares[fare.ID.Hex()] = fare
	return nil
}

func (r *InMemoryTripRepository) GetRideFareByID(ctx context.Context, fareID string) (*domain.RideFareModel, error) {
	fare, exists := r.rideFares[fareID]
	if !exists {
		return nil, fmt.Errorf("ride fare with ID %s not found", fareID)
	}

	return fare, nil
}
