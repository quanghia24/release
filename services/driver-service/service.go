package main

import (
	"fmt"
	math "math/rand/v2"
	pb "ride-sharing/shared/proto/driver"
	"ride-sharing/shared/util"
	"sync"

	"github.com/mmcloughlin/geohash"
)

type service struct {
	drivers []*driverInMap

	mu sync.RWMutex
}

type driverInMap struct {
	Driver *pb.Driver
	// Index
	// route
}

func NewService() *service {
	return &service{
		drivers: []*driverInMap{},
	}
}

func (s *service) RegisterDriver(driverID, packageSlug string) (*pb.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	randomIndex := math.IntN(len(PredefinedRoutes))
	randomRoute := PredefinedRoutes[randomIndex]

	randomAvatar := util.GetRandomAvatar(randomIndex)
	randomPlate := GenerateRandomPlate()

	geohash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	driver := &pb.Driver{
		Id:             driverID,
		Geohash:        geohash,
		Location:       &pb.Location{Latitude: randomRoute[0][0], Longitude: randomRoute[0][1]},
		Name:           "Lec Le",
		PackageSlug:    packageSlug,
		ProfilePicture: randomAvatar,
		CarPlate:       randomPlate,
	}

	s.drivers = append(s.drivers, &driverInMap{
		Driver: driver,
	})

	return driver, nil
}

func (s *service) UnregisterDriver(driverID string) (*pb.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, d := range s.drivers {
		if d.Driver.Id == driverID {
			s.drivers = append(s.drivers[:i], s.drivers[i+1:]...)
			return &pb.Driver{
				Id: d.Driver.Id,
			}, nil
		}
	}

	return nil, fmt.Errorf("driver not found")
}
