package grpc

import (
	"context"
	"log"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/events"
	pb "ride-sharing/shared/proto/trip"

	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer

	service   domain.TripService
	publisher *events.TripEventPublisher
}

func NewgRPCHandler(server *grpc.Server, service domain.TripService, publisher *events.TripEventPublisher) {
	handler := &gRPCHandler{
		service:   service,
		publisher: publisher,
	}

	pb.RegisterTripServiceServer(server, handler)
}

func (h *gRPCHandler) TripPreview(ctx context.Context, req *pb.TripPreviewRequest) (*pb.TripPreviewResponse, error) {
	pickup := req.GetPickup()
	destination := req.GetDestination()

	pickupCoord := &types.Coordinate{
		Latitude:  pickup.GetLatitude(),
		Longitude: pickup.GetLongitude(),
	}

	destinationCoord := &types.Coordinate{
		Latitude:  destination.GetLatitude(),
		Longitude: destination.GetLongitude(),
	}

	userID := req.GetUserID()

	route, err := h.service.GetRoute(
		ctx,
		pickupCoord,
		destinationCoord,
	)
	if err != nil {
		log.Println(err)

		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	// 1. Estimate the ride fares price based on the route (ex: distance, duration, etc.)
	estimatedFares := h.service.EstimatePackagesWithRoute(route)

	// 2. Store the ride fares in the database and return the response to the client
	fares, err := h.service.GenerateTripFares(estimatedFares, userID, route)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate trip fares: %v", err)
	}

	return &pb.TripPreviewResponse{
		Route:     route.ToProto(),
		RideFares: domain.ToRideFaresProto(fares),
	}, nil
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	fareID := req.GetRideFareID()
	userID := req.GetUserID()

	fare, err := h.service.GetAndValidateRideFare(ctx, fareID, userID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to validate ride fare: %v", err)
	}

	trip, err := h.service.CreateTrip(ctx, fare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create trip: %v", err)
	}

	// send to async processing
	if err := h.publisher.PublishTripCreated(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish trip created event: %v", err)
	}

	return &pb.CreateTripResponse{
		TripID: trip.ID.Hex(),
		Trip: &pb.Trip{
			Id:               trip.ID.Hex(),
			UserID:           trip.UserID,
			Status:           trip.Status,
			SelectedRideFare: fare.ToProto(),
			TripDriver:       trip.Driver,
		},
	}, nil
}
