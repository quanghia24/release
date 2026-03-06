package types

import (
	pb "ride-sharing/shared/proto/trip"
)

type OSRMRouteServiceResponse struct {
	Routes []struct {
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
		Duration float64 `json:"duration"`
		Distance float64 `json:"distance"`
	} `json:"routes"`
}

func (r *OSRMRouteServiceResponse) ToProto() *pb.Route {
	routes := r.Routes[0]
	geometries := routes.Geometry.Coordinates
	coordinates := make([]*pb.Coordinate, len(geometries))
	for i, coord := range geometries {
		coordinates[i] = &pb.Coordinate{
			Latitude:  coord[0],
			Longitude: coord[1],
		}
	}
	return &pb.Route{
		Geometry: []*pb.Geometry{
			{
				Coordinates: coordinates,
			},
		},
		Duration: routes.Duration,
		Distance: routes.Distance,
	}
}

type PricingConfig struct {
	PricePerUnitOfDistance float64
	PricePerMinute         float64
}
