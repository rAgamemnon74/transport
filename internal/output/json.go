package output

import (
	"encoding/json"
	"time"

	"transport/internal/api"
)

// JSONOutput represents the standardized JSON output format
type JSONOutput struct {
	Type      string      `json:"type"`       // trip, departures, car, flight, taxi, bus
	Timestamp string      `json:"timestamp"`
	Origin    string      `json:"origin,omitempty"`
	Dest      string      `json:"destination,omitempty"`
	Data      interface{} `json:"data"`
}

// TripResult represents trip planning results
type TripResult struct {
	Trips []Trip `json:"trips"`
}

// Trip represents a single trip option
type Trip struct {
	Departure       string `json:"departure"`        // HH:MM
	Arrival         string `json:"arrival"`          // HH:MM
	DurationMinutes int    `json:"duration_minutes"`
	Changes         int    `json:"changes"`
	Legs            []Leg  `json:"legs"`
	GoogleMapsURL   string `json:"google_maps_url,omitempty"`
}

// Leg represents one segment of a trip
type Leg struct {
	Mode      string   `json:"mode"`       // metro, bus, train, tram, ship, walk
	Line      string   `json:"line,omitempty"`
	Direction string   `json:"direction,omitempty"`
	From      StopInfo `json:"from"`
	To        StopInfo `json:"to"`
	Departure string   `json:"departure"`
	Arrival   string   `json:"arrival"`
	Duration  int      `json:"duration_minutes"`
	Coords    []Coord  `json:"coords,omitempty"`
}

// StopInfo represents a stop/station
type StopInfo struct {
	Name     string  `json:"name"`
	Platform string  `json:"platform,omitempty"`
	Lat      float64 `json:"lat,omitempty"`
	Lon      float64 `json:"lon,omitempty"`
}

// Coord represents a coordinate point
type Coord struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// DeparturesResult represents departure board results
type DeparturesResult struct {
	StopName   string      `json:"stop_name"`
	Departures []Departure `json:"departures"`
}

// Departure represents a single departure
type Departure struct {
	Line        string `json:"line"`
	Destination string `json:"destination"`
	Departure   string `json:"departure"`      // HH:MM
	Expected    string `json:"expected"`       // HH:MM (real-time)
	MinutesAway int    `json:"minutes_away"`
	Platform    string `json:"platform,omitempty"`
	Mode        string `json:"mode"`           // bus, metro, train, etc.
	Delayed     bool   `json:"delayed"`
}

// CarResult represents car journey results
type CarResult struct {
	DistanceKm    float64    `json:"distance_km"`
	DurationMin   int        `json:"duration_minutes"`
	FuelNeeded    float64    `json:"fuel_needed_liters"`
	FuelCost      float64    `json:"fuel_cost_sek"`
	GoogleMapsURL string     `json:"google_maps_url"`
	FuelStops     []FuelStop `json:"fuel_stops,omitempty"`
}

// FuelStop represents a recommended fuel stop
type FuelStop struct {
	Name      string  `json:"name"`
	AtKm      float64 `json:"at_km"`
	FuelLevel float64 `json:"fuel_level_percent"`
}

// FlightResult represents flight search results
type FlightResult struct {
	Flights []FlightOption `json:"flights"`
}

// FlightOption represents a flight booking option
type FlightOption struct {
	Airline    string `json:"airline"`
	From       string `json:"from_airport"`
	To         string `json:"to_airport"`
	Price      string `json:"price,omitempty"`
	BookingURL string `json:"booking_url"`
	Direct     bool   `json:"direct"`
}

// TaxiResult represents taxi fare estimation results
type TaxiResult struct {
	DistanceKm  float64        `json:"distance_km"`
	DurationMin int            `json:"duration_minutes"`
	Estimates   []TaxiEstimate `json:"estimates"`
}

// TaxiEstimate represents a taxi company estimate
type TaxiEstimate struct {
	Company    string  `json:"company"`
	Estimated  float64 `json:"estimated_sek"`
	FixedPrice float64 `json:"fixed_price_sek,omitempty"`
	BookingURL string  `json:"booking_url"`
	DeepLink   string  `json:"deep_link,omitempty"`
}

// BusResult represents long-distance bus results
type BusResult struct {
	Routes []BusRoute `json:"routes"`
}

// BusRoute represents a bus route option
type BusRoute struct {
	Operator   string `json:"operator"`
	Departure  string `json:"departure"`
	Arrival    string `json:"arrival"`
	Duration   string `json:"duration"`
	Price      string `json:"price,omitempty"`
	BookingURL string `json:"booking_url"`
}

// NewOutput creates a new JSON output wrapper
func NewOutput(outputType, origin, dest string) *JSONOutput {
	return &JSONOutput{
		Type:      outputType,
		Timestamp: time.Now().Format(time.RFC3339),
		Origin:    origin,
		Dest:      dest,
	}
}

// Marshal converts the output to JSON string
func (o *JSONOutput) Marshal() (string, error) {
	data, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FormatJourneysJSON converts API journeys to JSON format
func FormatJourneysJSON(origin, dest string, journeys []api.Journey, mapsURLGen func(api.Journey) string) string {
	output := NewOutput("trip", origin, dest)

	trips := make([]Trip, 0, len(journeys))
	for _, j := range journeys {
		trip := Trip{
			DurationMinutes: j.TripDuration / 60,
			Changes:         j.Interchanges,
			Legs:            make([]Leg, 0, len(j.Legs)),
		}

		if mapsURLGen != nil {
			trip.GoogleMapsURL = mapsURLGen(j)
		}

		for i, l := range j.Legs {
			leg := Leg{
				From: StopInfo{
					Name:     l.Origin.GetStopName(),
					Platform: l.Origin.GetPlatform(),
				},
				To: StopInfo{
					Name:     l.Destination.GetStopName(),
					Platform: l.Destination.GetPlatform(),
				},
				Departure: extractTime(l.Origin.DepartureTimePlanned),
				Arrival:   extractTime(l.Destination.ArrivalTimePlanned),
				Duration:  l.Duration / 60,
			}

			// Add coordinates if available
			if len(l.Origin.Coord) >= 2 {
				leg.From.Lat = l.Origin.Coord[0]
				leg.From.Lon = l.Origin.Coord[1]
			}
			if len(l.Destination.Coord) >= 2 {
				leg.To.Lat = l.Destination.Coord[0]
				leg.To.Lon = l.Destination.Coord[1]
			}

			// Set first leg departure as trip departure
			if i == 0 {
				trip.Departure = leg.Departure
			}
			// Set last leg arrival as trip arrival
			if i == len(j.Legs)-1 {
				trip.Arrival = leg.Arrival
			}

			// Determine mode and line
			if l.Transportation != nil && !l.Transportation.IsWalking() {
				leg.Line = l.Transportation.GetLineName()
				leg.Direction = l.Transportation.GetDirection()
				leg.Mode = getModeString(l.Transportation)
			} else {
				leg.Mode = "walk"
			}

			trip.Legs = append(trip.Legs, leg)
		}

		trips = append(trips, trip)
	}

	output.Data = TripResult{Trips: trips}

	result, _ := output.Marshal()
	return result
}

// FormatDeparturesJSON converts departures to JSON format
func FormatDeparturesJSON(siteName string, departures []api.Departure) string {
	output := NewOutput("departures", siteName, "")

	deps := make([]Departure, 0, len(departures))
	now := time.Now()

	for _, d := range departures {
		dep := Departure{
			Line:        d.Line.Designation,
			Destination: d.Destination,
			Departure:   d.Scheduled,
			Expected:    d.Expected,
			Platform:    d.StopPoint.Designation,
			Mode:        d.Line.TransportMode,
		}

		// Calculate minutes away
		if expectedTime, err := time.Parse("15:04:05", d.Expected); err == nil {
			expectedToday := time.Date(now.Year(), now.Month(), now.Day(),
				expectedTime.Hour(), expectedTime.Minute(), expectedTime.Second(), 0, now.Location())
			dep.MinutesAway = int(expectedToday.Sub(now).Minutes())
			if dep.MinutesAway < 0 {
				dep.MinutesAway = 0
			}
		}

		// Check if delayed
		dep.Delayed = d.Scheduled != d.Expected

		deps = append(deps, dep)
	}

	output.Data = DeparturesResult{
		StopName:   siteName,
		Departures: deps,
	}

	result, _ := output.Marshal()
	return result
}

// Helper to extract time from ISO format
func extractTime(isoTime string) string {
	if len(isoTime) < 16 {
		return isoTime
	}
	// Format: 2024-01-20T14:35:00
	return isoTime[11:16]
}

// Helper to get mode string from transportation
func getModeString(t *api.Transportation) string {
	if t == nil || t.Product == nil {
		return "unknown"
	}
	switch t.Product.Class {
	case api.ProductClassMetro:
		return "metro"
	case api.ProductClassBus:
		return "bus"
	case api.ProductClassTrain:
		return "train"
	case api.ProductClassTram:
		return "tram"
	case api.ProductClassFerry:
		return "ship"
	default:
		return "unknown"
	}
}
