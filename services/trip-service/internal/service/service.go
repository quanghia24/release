package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTrip(ctx context.Context, rideFare *domain.RideFareModel) (*domain.TripModel, error) {
	trip := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   rideFare.UserID,
		Status:   "pending",
		RideFare: rideFare,
		Driver:   &pb.TripDriver{},
	}
	return s.repo.CreateTrip(ctx, trip)
}

func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OSRMRouteServiceResponse, error) {
	// baseUrl := "https://osrm.selfmadeengineer.com"
	baseUrl := "http://router.project-osrm.org"
	url := fmt.Sprintf(
		"%s/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		baseUrl,
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM APIs %w", err)
	}
	defer resp.Body.Close()

	var routeResp tripTypes.OSRMRouteServiceResponse

	if err := json.NewDecoder(resp.Body).Decode(&routeResp); err != nil {
		return nil, fmt.Errorf("failed to decode OSRM response: %w", err)
	}

	return &routeResp, nil
}

func (s *service) EstimatePackagesWithRoute(route *tripTypes.OSRMRouteServiceResponse) []*domain.RideFareModel {
	baseFares := getBaseFares()
	estimatedFares := make([]*domain.RideFareModel, len(baseFares))

	for i, fare := range baseFares {
		estimatedFares[i] = estimateFareRoute(fare, route)
	}

	return estimatedFares
}

func (s *service) GenerateTripFares(rideFares []*domain.RideFareModel, userID string, route *tripTypes.OSRMRouteServiceResponse) ([]*domain.RideFareModel, error) {
	fares := make([]*domain.RideFareModel, len(rideFares))
	for i, f := range rideFares {
		id := primitive.NewObjectID()

		fare := &domain.RideFareModel{
			ID:                id,
			UserID:            userID,
			PackageSlug:       f.PackageSlug,
			TotalPriceInCents: f.TotalPriceInCents,
			Route: 					 route,
		}

		if err := s.repo.SaveRideFare(context.Background(), fare); err != nil {
			return nil, fmt.Errorf("failed to save trip fare: %w", err)
		}

		fares[i] = fare
	}

	return fares, nil
}

func (s *service) GetAndValidateRideFare(ctx context.Context, fareID, userID string) (*domain.RideFareModel, error) {
	fare, err := s.repo.GetRideFareByID(ctx, fareID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ride fare by ID: %w", err)
	}

	if fare == nil {
		return nil, fmt.Errorf("ride fare with ID %s not found", fareID)
	}

	if fare.UserID != userID {
		return nil, fmt.Errorf("fare doesn't belong to user ID %s", userID)
	}

	return fare, nil
}

func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 150,
		},
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 350,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 550,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 1000,
		},
	}
}

func estimateFareRoute(fare *domain.RideFareModel, route *tripTypes.OSRMRouteServiceResponse) *domain.RideFareModel {
	priceConfig := defaultPricingConfig()

	distance := fare.TotalPriceInCents * priceConfig.PricePerUnitOfDistance
	duration := route.Routes[0].Duration * priceConfig.PricePerMinute
	total := distance + duration

	return &domain.RideFareModel{
		PackageSlug:       fare.PackageSlug,
		TotalPriceInCents: total,
	}
}

func defaultPricingConfig() *tripTypes.PricingConfig {
	return &tripTypes.PricingConfig{
		PricePerUnitOfDistance: 1.5,
		PricePerMinute:         0.25,
	}
}
