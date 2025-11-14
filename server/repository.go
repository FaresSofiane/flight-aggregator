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
	
	var flights []FlightServer1
	if err := json.NewDecoder(resp.Body).Decode(&flights); err != nil {
		return nil, err
	}
	
	result := make([]UnifiedFlight, 0, len(flights))
	for _, f := range flights {
		duration := f.ArrivalTime.Sub(f.DepartureTime)
		result = append(result, UnifiedFlight{
			BookingID:        f.BookingID,
			PassengerName:    f.PassengerName,
			DepartureAirport: f.DepartureAirport,
			ArrivalAirport:   f.ArrivalAirport,
			DepartureTime:    f.DepartureTime,
			ArrivalTime:      f.ArrivalTime,
			TravelTime:       duration,
			Price:            f.Price,
			Currency:         f.Currency,
			FlightNumbers:    []string{f.FlightNumber},
		})
	}
	
	return result, nil
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
	
	var flights []FlightServer2
	if err := json.NewDecoder(resp.Body).Decode(&flights); err != nil {
		return nil, err
	}
	
	result := make([]UnifiedFlight, 0)
	for _, f := range flights {
		if len(f.Segments) == 0 {
			continue
		}
		
		first := f.Segments[0]
		last := f.Segments[len(f.Segments)-1]
		
		start := first.Flight.Depart
		end := last.Flight.Arrive
		duration := end.Sub(start)
		
		numbers := make([]string, len(f.Segments))
		for i, seg := range f.Segments {
			numbers[i] = seg.Flight.Number
		}
		
		name := fmt.Sprintf("%s %s", f.Traveler.FirstName, f.Traveler.LastName)
		
		result = append(result, UnifiedFlight{
			BookingID:        f.Reference,
			PassengerName:    name,
			DepartureAirport: first.Flight.From,
			ArrivalAirport:   last.Flight.To,
			DepartureTime:    start,
			ArrivalTime:      end,
			TravelTime:       duration,
			Price:            f.Total.Amount,
			Currency:         f.Total.Currency,
			FlightNumbers:    numbers,
		})
	}
	
	return result, nil
}
