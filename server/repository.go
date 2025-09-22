package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FlightRepository interface {
	GetFlights(ctx context.Context) ([]UnifiedFlight, error)
}

type Server1Repository struct {
	BaseURL string
}

func NewServer1Repository(baseURL string) *Server1Repository {
	return &Server1Repository{
		BaseURL: baseURL,
	}
}

func (r *Server1Repository) GetFlights(ctx context.Context) ([]UnifiedFlight, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	
	req, err := http.NewRequestWithContext(ctx, "GET", r.BaseURL+"/flights", nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var flightResp FlightServer1Response
	if err := json.NewDecoder(resp.Body).Decode(&flightResp); err != nil {
		return nil, err
	}
	
	var unifiedFlights []UnifiedFlight
	for _, flight := range flightResp.Flights {
		travelTime := flight.ArrivalTime.Sub(flight.DepartureTime)
		
		unifiedFlights = append(unifiedFlights, UnifiedFlight{
			BookingID:        flight.BookingID,
			PassengerName:    flight.PassengerName,
			DepartureAirport: flight.DepartureAirport,
			ArrivalAirport:   flight.ArrivalAirport,
			DepartureTime:    flight.DepartureTime,
			ArrivalTime:      flight.ArrivalTime,
			TravelTime:       travelTime,
			Price:            flight.Price,
			Currency:         flight.Currency,
			FlightNumbers:    []string{flight.FlightNumber},
		})
	}
	
	return unifiedFlights, nil
}

type Server2Repository struct {
	BaseURL string
}

func NewServer2Repository(baseURL string) *Server2Repository {
	return &Server2Repository{
		BaseURL: baseURL,
	}
}

func (r *Server2Repository) GetFlights(ctx context.Context) ([]UnifiedFlight, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	
	req, err := http.NewRequestWithContext(ctx, "GET", r.BaseURL+"/flight_to_book", nil)
	if err != nil {
		return nil, err
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var flightResp FlightServer2Response
	if err := json.NewDecoder(resp.Body).Decode(&flightResp); err != nil {
		return nil, err
	}
	
	var unifiedFlights []UnifiedFlight
	for _, flight := range flightResp.FlightToBook {
		if len(flight.Segments) == 0 {
			continue
		}
		
		firstSegment := flight.Segments[0]
		lastSegment := flight.Segments[len(flight.Segments)-1]
		
		departureTime := firstSegment.Flight.Depart
		arrivalTime := lastSegment.Flight.Arrive
		travelTime := arrivalTime.Sub(departureTime)
		
		var flightNumbers []string
		for _, segment := range flight.Segments {
			flightNumbers = append(flightNumbers, segment.Flight.Number)
		}
		
		passengerName := fmt.Sprintf("%s %s", flight.Traveler.FirstName, flight.Traveler.LastName)
		
		unifiedFlights = append(unifiedFlights, UnifiedFlight{
			BookingID:        flight.Reference,
			PassengerName:    passengerName,
			DepartureAirport: firstSegment.Flight.From,
			ArrivalAirport:   lastSegment.Flight.To,
			DepartureTime:    departureTime,
			ArrivalTime:      arrivalTime,
			TravelTime:       travelTime,
			Price:            flight.Total.Amount,
			Currency:         flight.Total.Currency,
			FlightNumbers:    flightNumbers,
		})
	}
	
	return unifiedFlights, nil
}
