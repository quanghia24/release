package main

import (
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"
)

type tripReviewRequest struct {
	UserID      string           `json:"userID"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (r *tripReviewRequest) toProto() *pb.TripPreviewRequest {
	return &pb.TripPreviewRequest{
		UserID: r.UserID,
		Pickup: &pb.Coordinate{
			Latitude:  r.Pickup.Latitude,
			Longitude: r.Pickup.Longitude,
		},
		Destination: &pb.Coordinate{
			Latitude:  r.Destination.Latitude,
			Longitude: r.Destination.Longitude,
		},
	}
}

type startTripRequest struct {
	RideFareID string `json:"rideFareID"`
	UserID     string `json:"userID"`
}

func (c *startTripRequest) toProto() *pb.CreateTripRequest {
	return &pb.CreateTripRequest{
		RideFareID: c.RideFareID,
		UserID:     c.UserID,
	}
}
