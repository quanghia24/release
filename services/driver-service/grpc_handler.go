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
	return nil, status.Error(codes.Unimplemented, "method RegisterDriver not implemented")
}
func (h *grpcHandler) UnregisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method UnregisterDriver not implemented")
}
