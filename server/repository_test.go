package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer1Repository_GetFlights(t *testing.T) {
	mockData := []FlightServer1{
		{
			BookingID:        "A10001",
			Status:           "confirmed",
			PassengerName:    "Marie Curie",
			FlightNumber:     "JL046",
			DepartureAirport: "CDG",
			ArrivalAirport:   "HND",
			DepartureTime:    time.Date(2026, 1, 1, 13, 0, 0, 0, time.UTC),
			ArrivalTime:      time.Date(2026, 1, 2, 8, 30, 0, 0, time.UTC),
			Price:            850.0,
			Currency:         "EUR",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/flights", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	repo := NewServer1Repository(server.URL)
	flights, err := repo.GetFlights(context.Background())

	assert.NoError(t, err)
	assert.Len(t, flights, 1)
	assert.Equal(t, "A10001", flights[0].BookingID)
	assert.Equal(t, "Marie Curie", flights[0].PassengerName)
	assert.Equal(t, 850.0, flights[0].Price)
	assert.Equal(t, []string{"JL046"}, flights[0].FlightNumbers)
}

func TestServer2Repository_GetFlights(t *testing.T) {
	mockData := []FlightServer2{
		{
			Reference: "B30001",
			Status:    "confirmed",
			Traveler: Traveler{
				FirstName: "Marie",
				LastName:  "Curie",
			},
			Segments: []Segment{
				{
					Flight: Flight{
						Number: "AF276",
						From:   "CDG",
						To:     "HND",
						Depart: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
						Arrive: time.Date(2026, 1, 1, 23, 0, 0, 0, time.UTC),
					},
				},
			},
			Total: Total{
				Amount:   950.0,
				Currency: "EUR",
			},
		},
		{
			Reference: "B30004",
			Status:    "confirmed",
			Traveler: Traveler{
				FirstName: "Marie",
				LastName:  "Curie",
			},
			Segments: []Segment{
				{
					Flight: Flight{
						Number: "KE902",
						From:   "CDG",
						To:     "ICN",
						Depart: time.Date(2026, 1, 1, 9, 30, 0, 0, time.UTC),
						Arrive: time.Date(2026, 1, 1, 18, 0, 0, 0, time.UTC),
					},
				},
				{
					Flight: Flight{
						Number: "KE711",
						From:   "ICN",
						To:     "HND",
						Depart: time.Date(2026, 1, 1, 20, 0, 0, 0, time.UTC),
						Arrive: time.Date(2026, 1, 2, 0, 30, 0, 0, time.UTC),
					},
				},
			},
			Total: Total{
				Amount:   880.0,
				Currency: "EUR",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/flight_to_book", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	repo := NewServer2Repository(server.URL)
	flights, err := repo.GetFlights(context.Background())

	assert.NoError(t, err)
	assert.Len(t, flights, 2)
	
	assert.Equal(t, "B30001", flights[0].BookingID)
	assert.Equal(t, "Marie Curie", flights[0].PassengerName)
	assert.Equal(t, 950.0, flights[0].Price)
	assert.Equal(t, []string{"AF276"}, flights[0].FlightNumbers)
	
	assert.Equal(t, "B30004", flights[1].BookingID)
	assert.Equal(t, 880.0, flights[1].Price)
	assert.Equal(t, []string{"KE902", "KE711"}, flights[1].FlightNumbers)
}
