package main

import pb "ride-sharing/shared/proto/driver"

type service struct {
	driver []*driverInMap
}

type driverInMap struct {
	Driver *pb.Driver
	// Index
	// route
}

func NewService() *service {
	return &service{
		driver: []*driverInMap{},
	}
}

