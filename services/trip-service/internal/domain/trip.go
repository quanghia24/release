package domain

import (
	"context"
	"ride-sharing/shared/types"

	tripTypes "ride-sharing/services/trip-service/pkg/types"

	pb "ride-sharing/shared/proto/trip"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID       primitive.ObjectID
	UserID   string
	Status   string
	RideFare *RideFareModel
	Driver   *pb.TripDriver
}

func (t *TripModel) ToProto() *pb.CreateTripResponse {
	return &pb.CreateTripResponse{}
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRideFare(ctx context.Context, fare *RideFareModel) error
	GetRideFareByID(ctx context.Context, fareID string) (*RideFareModel, error)
}
type TripService interface {
	CreateTrip(ctx context.Context, rideFare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OSRMRouteServiceResponse, error)
	EstimatePackagesWithRoute(route *tripTypes.OSRMRouteServiceResponse) []*RideFareModel
	GenerateTripFares(rideFares []*RideFareModel, userID string, route *tripTypes.OSRMRouteServiceResponse) ([]*RideFareModel, error)
	GetAndValidateRideFare(ctx context.Context, fareID, userID string) (*RideFareModel, error)
}
