package output

import (
	"fmt"
	"net/url"

	"transport/internal/bus"
	"transport/internal/car"
	"transport/internal/flight"
	"transport/internal/resrobot"
	"transport/internal/taxi"
)

// FormatCarJSON converts car trip results to JSON format
func FormatCarJSON(from, to string, distanceKm, startFuel float64, profile car.VehicleProfile) string {
	output := NewOutput("car", from, to)

	fuelNeeded := profile.CalculateFuel(distanceKm)
	fuelCost := fuelNeeded * 19.5 // approximate diesel price in SEK/L

	carResult := CarResult{
		DistanceKm:    distanceKm,
		DurationMin:   int(distanceKm / 80 * 60), // approximate at 80 km/h average
		FuelNeeded:    fuelNeeded,
		FuelCost:      fuelCost,
		GoogleMapsURL: car.GenerateGoogleMapsURL(from, to),
	}

	// Calculate fuel stops if needed
	stops := profile.CalculateFuelStops(distanceKm, startFuel)
	for _, stop := range stops {
		carResult.FuelStops = append(carResult.FuelStops, FuelStop{
			Name:      stop.Location,
			AtKm:      stop.AtKm,
			FuelLevel: stop.FuelRemaining / profile.TankSizeLiters * 100,
		})
	}

	output.Data = carResult
	result, _ := output.Marshal()
	return result
}

// FormatFlightJSON converts flight search results to JSON format
func FormatFlightJSON(search flight.FlightSearch) string {
	output := NewOutput("flight", search.OriginCity, search.DestinationCity)

	flights := make([]FlightOption, 0)

	// Generate booking links for various airlines
	dateStr := ""
	if !search.DepartureDate.IsZero() {
		dateStr = search.DepartureDate.Format("2006-01-02")
	}

	// SAS
	flights = append(flights, FlightOption{
		Airline:    "SAS",
		From:       search.OriginCode,
		To:         search.DestinationCode,
		BookingURL: fmt.Sprintf("https://www.flysas.com/search/?from=%s&to=%s&date=%s", search.OriginCode, search.DestinationCode, dateStr),
		Direct:     true,
	})

	// Norwegian
	flights = append(flights, FlightOption{
		Airline:    "Norwegian",
		From:       search.OriginCode,
		To:         search.DestinationCode,
		BookingURL: fmt.Sprintf("https://www.norwegian.com/flight-search/?from=%s&to=%s&date=%s", search.OriginCode, search.DestinationCode, dateStr),
		Direct:     true,
	})

	// Google Flights
	flights = append(flights, FlightOption{
		Airline:    "Google Flights",
		From:       search.OriginCode,
		To:         search.DestinationCode,
		BookingURL: fmt.Sprintf("https://www.google.com/travel/flights?q=flights%%20from%%20%s%%20to%%20%s%%20%s", search.OriginCode, search.DestinationCode, url.QueryEscape(dateStr)),
		Direct:     false,
	})

	output.Data = FlightResult{Flights: flights}
	result, _ := output.Marshal()
	return result
}

// FormatTaxiJSON converts taxi search results to JSON format
func FormatTaxiJSON(search taxi.TaxiSearch) string {
	output := NewOutput("taxi", search.From, search.To)

	estimates := make([]TaxiEstimate, 0)
	for _, e := range search.Estimates {
		estimates = append(estimates, TaxiEstimate{
			Company:    e.Company,
			Estimated:  e.Estimated,
			FixedPrice: e.FixedPrice,
			BookingURL: e.BookingURL,
			DeepLink:   e.DeepLink,
		})
	}

	taxiResult := TaxiResult{
		DistanceKm:  search.Route.DistanceKm,
		DurationMin: int(search.Route.DurationMin),
		Estimates:   estimates,
	}

	output.Data = taxiResult
	result, _ := output.Marshal()
	return result
}

// FormatBusJSON converts bus search results to JSON format
func FormatBusJSON(search bus.BusSearch) string {
	output := NewOutput("bus", search.From, search.To)

	routes := make([]BusRoute, 0)
	for _, r := range search.Routes {
		routes = append(routes, BusRoute{
			Operator:   r.Operator,
			Departure:  "",  // Not available in static data
			Arrival:    "",
			Duration:   r.Duration,
			Price:      fmt.Sprintf("%d-%d SEK", r.PriceFrom, r.PriceTo),
			BookingURL: r.BookingURL,
		})
	}

	output.Data = BusResult{Routes: routes}
	result, _ := output.Marshal()
	return result
}

// FormatResRobotJSON converts ResRobot trip results to JSON format
func FormatResRobotJSON(origin, dest string, trips []resrobot.ParsedTrip) string {
	output := NewOutput("trip", origin, dest)

	tripResults := make([]Trip, 0, len(trips))
	for _, t := range trips {
		trip := Trip{
			DurationMinutes: int(t.Duration.Minutes()),
			Changes:         t.Interchanges,
			Legs:            make([]Leg, 0, len(t.Legs)),
		}

		for i, l := range t.Legs {
			mode := "unknown"
			switch l.Category {
			case resrobot.CatMetro:
				mode = "metro"
			case resrobot.CatBus, resrobot.CatNattbus:
				mode = "bus"
			case resrobot.CatPendeltag, resrobot.CatRegionaltag, resrobot.CatSJ, resrobot.CatSnabbtag, resrobot.CatNorrtag:
				mode = "train"
			case resrobot.CatSparvagn:
				mode = "tram"
			case resrobot.CatFerry:
				mode = "ship"
			}
			if l.IsWalk {
				mode = "walk"
			}

			leg := Leg{
				Mode:      mode,
				Line:      l.Line,
				Direction: l.Direction,
				From: StopInfo{
					Name:     l.Origin,
					Platform: l.OriginTrack,
				},
				To: StopInfo{
					Name:     l.Destination,
					Platform: l.DestTrack,
				},
				Departure: l.DepartureTime.Format("15:04"),
				Arrival:   l.ArrivalTime.Format("15:04"),
				Duration:  int(l.ArrivalTime.Sub(l.DepartureTime).Minutes()),
			}

			// Set first leg departure as trip departure
			if i == 0 {
				trip.Departure = leg.Departure
			}
			// Set last leg arrival as trip arrival
			if i == len(t.Legs)-1 {
				trip.Arrival = leg.Arrival
			}

			trip.Legs = append(trip.Legs, leg)
		}

		tripResults = append(tripResults, trip)
	}

	output.Data = TripResult{Trips: tripResults}
	result, _ := output.Marshal()
	return result
}
