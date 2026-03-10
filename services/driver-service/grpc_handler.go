package main

import (
	"context"
	pb "ride-sharing/shared/proto/driver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	pb.UnimplementedDriverServiceServer

	service *service
}

func NewgRPCHandler(server *grpc.Server, service *service) {
	handler := &grpcHandler{
		service: service,
	}

	pb.RegisterDriverServiceServer(server, handler)
}

func (h *grpcHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	driverID := req.GetDriverID()
	packageSlug := req.GetPackageSlug()

	driver, err := h.service.RegisterDriver(driverID, packageSlug)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register driver")
	}

	return &pb.RegisterDriverResponse{
		Driver: driver,
	}, nil
}

func (h *grpcHandler) UnregisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	driverID := req.GetDriverID()

	driver, err := h.service.UnregisterDriver(driverID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to unregister driver")
	}

	return &pb.RegisterDriverResponse{
		Driver: driver,
	}, nil
}
