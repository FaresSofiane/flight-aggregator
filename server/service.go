package main

import (
	"context"
	"sort"
)

type FlightService struct {
	repositories []FlightRepository
}

func NewFlightService(repositories []FlightRepository) *FlightService {
	return &FlightService{
		repositories: repositories,
	}
}

func (s *FlightService) GetAllFlights(ctx context.Context) ([]UnifiedFlight, error) {
	results := []UnifiedFlight{}
	
	for _, repo := range s.repositories {
		flights, err := repo.GetFlights(ctx)
		if err != nil {
			return nil, err
		}
		results = append(results, flights...)
	}
	
	return results, nil
}

func (s *FlightService) GetFlightsSortedBy(ctx context.Context, sortBy string) ([]UnifiedFlight, error) {
	flights, err := s.GetAllFlights(ctx)
	if err != nil {
		return nil, err
	}
	
	switch sortBy {
	case "price":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].Price < flights[j].Price
		})
	case "departure_date":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].DepartureTime.Before(flights[j].DepartureTime)
		})
	case "travel_time":
		sort.Slice(flights, func(i, j int) bool {
			return flights[i].TravelTime < flights[j].TravelTime
		})
	}
	
	return flights, nil
}
