package main

import (
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type meters int

type distancer interface {
	Distance(origin string, destination string) (meters, error)
}

type googleMapsDistancer struct {
	mapsClient *maps.Client
}

func (g *googleMapsDistancer) Distance(origin string, destination string) (meters, error) {
	req := &maps.DistanceMatrixRequest{
		Origins:      []string{origin},
		Destinations: []string{destination},
		Mode:         "ModeWalking",
	}
	resp, err := g.mapsClient.DistanceMatrix(context.Background(), req)

	if err != nil {
		return meters(0), err
	}

	return meters(resp.Rows[0].Elements[0].Distance.Meters), nil
}
