package main

import "time"

type FlightServer1Response struct {
	Flights []FlightServer1 `json:"flights"`
}

type FlightServer1 struct {
	BookingID        string    `json:"bookingId"`
	Status           string    `json:"status"`
	PassengerName    string    `json:"passengerName"`
	FlightNumber     string    `json:"flightNumber"`
	DepartureAirport string    `json:"departureAirport"`
	ArrivalAirport   string    `json:"arrivalAirport"`
	DepartureTime    time.Time `json:"departureTime"`
	ArrivalTime      time.Time `json:"arrivalTime"`
	Price            float64   `json:"price"`
	Currency         string    `json:"currency"`
}

type FlightServer2Response struct {
	FlightToBook []FlightServer2 `json:"flight_to_book"`
}

type FlightServer2 struct {
	Reference string    `json:"reference"`
	Status    string    `json:"status"`
	Traveler  Traveler  `json:"traveler"`
	Segments  []Segment `json:"segments"`
	Total     Total     `json:"total"`
}

type Traveler struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Segment struct {
	Flight Flight `json:"flight"`
}

type Flight struct {
	Number string    `json:"number"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Depart time.Time `json:"depart"`
	Arrive time.Time `json:"arrive"`
}

type Total struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type UnifiedFlight struct {
	BookingID        string
	PassengerName    string
	DepartureAirport string
	ArrivalAirport   string
	DepartureTime    time.Time
	ArrivalTime      time.Time
	TravelTime       time.Duration
	Price            float64
	Currency         string
	FlightNumbers    []string
}
