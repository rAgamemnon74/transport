package car

import (
	"fmt"
	"net/url"
	"strings"
)

// VehicleProfile contains fuel consumption data for a vehicle
type VehicleProfile struct {
	Name              string
	FuelType          string
	ShortDistanceRate float64 // L/100km for distances < 20km
	LongDistanceRate  float64 // L/100km for distances >= 20km
	ShortDistanceKm   float64 // Threshold for short distance
	TankSizeLiters    float64 // Fuel tank capacity
	ReservePercent    float64 // Don't go below this % of tank
}

// DefaultProfile returns the default vehicle profile (Tiguan Allspace 2018 Diesel)
func DefaultProfile() VehicleProfile {
	return VehicleProfile{
		Name:              "VW Tiguan Allspace 2018",
		FuelType:          "Diesel",
		ShortDistanceRate: 9.0,   // L/100km for short trips
		LongDistanceRate:  7.0,   // L/100km for long trips
		ShortDistanceKm:   20.0,  // Threshold in km
		TankSizeLiters:    58.0,  // Tiguan Allspace tank size
		ReservePercent:    15.0,  // Keep 15% reserve
	}
}

// CalculateFuel calculates fuel consumption for a given distance
func (v VehicleProfile) CalculateFuel(distanceKm float64) float64 {
	if distanceKm < v.ShortDistanceKm {
		return distanceKm * v.ShortDistanceRate / 100
	}
	return distanceKm * v.LongDistanceRate / 100
}

// GetFuelRate returns the appropriate fuel rate for a distance
func (v VehicleProfile) GetFuelRate(distanceKm float64) float64 {
	if distanceKm < v.ShortDistanceKm {
		return v.ShortDistanceRate
	}
	return v.LongDistanceRate
}

// UsableTank returns the usable fuel (accounting for reserve)
func (v VehicleProfile) UsableTank() float64 {
	return v.TankSizeLiters * (100 - v.ReservePercent) / 100
}

// MaxRange returns maximum range on a full tank (long distance rate)
func (v VehicleProfile) MaxRange() float64 {
	return v.UsableTank() / v.LongDistanceRate * 100
}

// FuelStop represents a recommended fuel stop
type FuelStop struct {
	AtKm         float64
	Location     string // General location description
	FuelUsed     float64
	FuelRemaining float64
}

// CalculateFuelStops calculates where to stop for fuel
func (v VehicleProfile) CalculateFuelStops(totalDistanceKm float64, startingFuelPercent float64) []FuelStop {
	if startingFuelPercent <= 0 {
		startingFuelPercent = 100
	}

	usableRange := v.MaxRange()
	startingRange := usableRange * startingFuelPercent / 100

	var stops []FuelStop

	// Use 85% of range as safety margin for stops
	safetyFactor := 0.85

	// Check if we can make it without stopping
	if totalDistanceKm <= startingRange*safetyFactor {
		return stops
	}

	// First stop based on starting fuel
	currentKm := startingRange * safetyFactor
	segmentStart := 0.0

	for currentKm < totalDistanceKm {
		// Don't add a stop if we're very close to destination (within 50km)
		if totalDistanceKm-currentKm < 50 {
			break
		}

		stop := FuelStop{
			AtKm:          currentKm,
			FuelUsed:      v.CalculateFuel(currentKm - segmentStart),
			FuelRemaining: v.TankSizeLiters * v.ReservePercent / 100,
		}

		// Describe approximate location
		percentComplete := currentKm / totalDistanceKm * 100
		if percentComplete < 30 {
			stop.Location = "ca 1/4 av vägen"
		} else if percentComplete < 45 {
			stop.Location = "ca 1/3 av vägen"
		} else if percentComplete < 55 {
			stop.Location = "halvvägs"
		} else if percentComplete < 70 {
			stop.Location = "ca 2/3 av vägen"
		} else {
			stop.Location = "ca 3/4 av vägen"
		}

		stops = append(stops, stop)

		// After refuel, can go another full range
		segmentStart = currentKm
		currentKm += usableRange * safetyFactor
	}

	return stops
}

// GenerateGoogleMapsURL creates a Google Maps directions URL
func GenerateGoogleMapsURL(from, to string) string {
	baseURL := "https://www.google.com/maps/dir/"
	fromEncoded := url.PathEscape(from)
	toEncoded := url.PathEscape(to)
	return baseURL + fromEncoded + "/" + toEncoded
}

// GenerateGoogleMapsURLWithParams creates a Google Maps URL with driving mode
func GenerateGoogleMapsURLWithParams(from, to string, waypoints []string) string {
	params := url.Values{}
	params.Set("api", "1")
	params.Set("origin", from)
	params.Set("destination", to)
	params.Set("travelmode", "driving")

	if len(waypoints) > 0 {
		params.Set("waypoints", strings.Join(waypoints, "|"))
	}

	return "https://www.google.com/maps/dir/?" + params.Encode()
}

// GenerateGasSearchURL creates a Google Maps search URL for gas stations
func GenerateGasSearchURL(nearLocation string, fuelType string) string {
	query := fuelType + " tankstation nära " + nearLocation
	params := url.Values{}
	params.Set("api", "1")
	params.Set("query", query)
	return "https://www.google.com/maps/search/?" + params.Encode()
}

// FormatCarTrip formats car trip information for display
func FormatCarTrip(from, to string, distanceKm float64, startFuel float64, profile VehicleProfile) string {
	var sb strings.Builder

	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	sb.WriteString(fmt.Sprintf(" 🚗 Bil: %s → %s\n", from, to))
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	sb.WriteString(fmt.Sprintf("  Fordon:    %s (%s)\n", profile.Name, profile.FuelType))
	sb.WriteString(fmt.Sprintf("  Tank:      %.0f liter (räckvidd ~%.0f km)\n", profile.TankSizeLiters, profile.MaxRange()))
	sb.WriteString("\n")

	if distanceKm > 0 {
		fuelNeeded := profile.CalculateFuel(distanceKm)
		rate := profile.GetFuelRate(distanceKm)

		sb.WriteString(fmt.Sprintf("  Avstånd:       %.0f km\n", distanceKm))
		sb.WriteString(fmt.Sprintf("  Förbrukning:   %.1f L/100km\n", rate))
		sb.WriteString(fmt.Sprintf("  Bränsle:       %.1f liter %s\n", fuelNeeded, strings.ToLower(profile.FuelType)))

		// Check if fuel stops are needed
		stops := profile.CalculateFuelStops(distanceKm, startFuel)
		if len(stops) > 0 {
			sb.WriteString("\n")
			sb.WriteString(fmt.Sprintf("  ⛽ Tankstopp behövs (%d st):\n", len(stops)))
			for i, stop := range stops {
				sb.WriteString(fmt.Sprintf("     %d. Efter ~%.0f km (%s)\n", i+1, stop.AtKm, stop.Location))
			}
			sb.WriteString("\n")
			sb.WriteString("  🔍 Sök tankstation längs rutten:\n")
			sb.WriteString(fmt.Sprintf("     %s\n", GenerateGasSearchURL("E4 mot "+to, profile.FuelType)))
		} else {
			sb.WriteString("\n")
			sb.WriteString("  ✓ Ingen tankning behövs under resan\n")
		}
	} else {
		sb.WriteString(fmt.Sprintf("  Förbrukning:   %.1f L/100km (< %.0f km)\n", profile.ShortDistanceRate, profile.ShortDistanceKm))
		sb.WriteString(fmt.Sprintf("                 %.1f L/100km (≥ %.0f km)\n", profile.LongDistanceRate, profile.ShortDistanceKm))
		sb.WriteString("\n")
		sb.WriteString("  💡 Ange avstånd med -d <km> för bränsleberäkning\n")
		sb.WriteString("     Ange tankläge med -f <procent> (t.ex. -f 50 för halvfull tank)\n")
	}

	sb.WriteString("\n")
	sb.WriteString("  🗺️  Google Maps:\n")
	sb.WriteString(fmt.Sprintf("     %s\n", GenerateGoogleMapsURLWithParams(from, to, nil)))
	sb.WriteString("\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	return sb.String()
}
