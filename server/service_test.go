package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFlightService_GetFlightsSortedBy(t *testing.T) {
	flights := []UnifiedFlight{
		{
			BookingID:     "B001",
			PassengerName: "Test User 1",
			DepartureTime: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
			ArrivalTime:   time.Date(2026, 1, 1, 18, 0, 0, 0, time.UTC),
			TravelTime:    6 * time.Hour,
			Price:         850.0,
			Currency:      "EUR",
		},
		{
			BookingID:     "B002",
			PassengerName: "Test User 2",
			DepartureTime: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
			ArrivalTime:   time.Date(2026, 1, 1, 20, 0, 0, 0, time.UTC),
			TravelTime:    10 * time.Hour,
			Price:         750.0,
			Currency:      "EUR",
		},
		{
			BookingID:     "B003",
			PassengerName: "Test User 3",
			DepartureTime: time.Date(2026, 1, 1, 14, 0, 0, 0, time.UTC),
			ArrivalTime:   time.Date(2026, 1, 1, 17, 0, 0, 0, time.UTC),
			TravelTime:    3 * time.Hour,
			Price:         950.0,
			Currency:      "EUR",
		},
	}

	mockRepo := &MockRepository{flights: flights}
	service := NewFlightService([]FlightRepository{mockRepo})
	ctx := context.Background()

	t.Run("Sort by price", func(t *testing.T) {
		sorted, err := service.GetFlightsSortedBy(ctx, "price")
		assert.NoError(t, err)
		assert.Len(t, sorted, 3)
		assert.Equal(t, 750.0, sorted[0].Price)
		assert.Equal(t, 850.0, sorted[1].Price)
		assert.Equal(t, 950.0, sorted[2].Price)
	})

	t.Run("Sort by departure_date", func(t *testing.T) {
		sorted, err := service.GetFlightsSortedBy(ctx, "departure_date")
		assert.NoError(t, err)
		assert.Len(t, sorted, 3)
		assert.Equal(t, "B002", sorted[0].BookingID)
		assert.Equal(t, "B001", sorted[1].BookingID)
		assert.Equal(t, "B003", sorted[2].BookingID)
	})

	t.Run("Sort by travel_time", func(t *testing.T) {
		sorted, err := service.GetFlightsSortedBy(ctx, "travel_time")
		assert.NoError(t, err)
		assert.Len(t, sorted, 3)
		assert.Equal(t, 3*time.Hour, sorted[0].TravelTime)
		assert.Equal(t, 6*time.Hour, sorted[1].TravelTime)
		assert.Equal(t, 10*time.Hour, sorted[2].TravelTime)
	})

	t.Run("Default sort (price)", func(t *testing.T) {
		sorted, err := service.GetFlightsSortedBy(ctx, "invalid")
		assert.NoError(t, err)
		assert.Len(t, sorted, 3)
		assert.Equal(t, 750.0, sorted[0].Price)
		assert.Equal(t, 850.0, sorted[1].Price)
		assert.Equal(t, 950.0, sorted[2].Price)
	})
}

type MockRepository struct {
	flights []UnifiedFlight
	err     error
}

func (m *MockRepository) GetFlights(ctx context.Context) ([]UnifiedFlight, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.flights, nil
}
