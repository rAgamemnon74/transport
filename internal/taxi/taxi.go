package taxi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	nominatimURL = "https://nominatim.openstreetmap.org/search"
	osrmURL      = "https://router.project-osrm.org/route/v1/driving"
	userAgent    = "transport-cli/1.0"
)

// Location represents a geocoded location
type Location struct {
	Name      string
	Lat       float64
	Lon       float64
	Address   string
	IsAirport bool
}

// Route represents a calculated route
type Route struct {
	From       Location
	To         Location
	DistanceKm float64
	DurationMin float64
}

// FareEstimate represents a taxi fare estimate
type FareEstimate struct {
	Company    string
	BaseFee    float64
	PerKmRate  float64
	PerHourRate float64
	Estimated  float64
	FixedPrice float64 // 0 if not applicable
	IsFixed    bool
	BookingURL string
	DeepLink   string // For apps with deep link support
}

// TaxiSearch contains search parameters
type TaxiSearch struct {
	From        string
	To          string
	Route       *Route
	Estimates   []FareEstimate
	Passengers  int
	IsWeekend   bool
}

// Geocode converts an address to coordinates using Nominatim
func Geocode(address string) (*Location, error) {
	params := url.Values{}
	params.Set("q", address+", Stockholm, Sweden")
	params.Set("format", "json")
	params.Set("limit", "1")
	params.Set("countrycodes", "se")

	req, err := http.NewRequest("GET", nominatimURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}
	defer resp.Body.Close()

	var results []struct {
		Lat         string `json:"lat"`
		Lon         string `json:"lon"`
		DisplayName string `json:"display_name"`
		Name        string `json:"name"`
		Class       string `json:"class"`
		Type        string `json:"type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode geocoding response: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("location not found: %s", address)
	}

	r := results[0]
	var lat, lon float64
	fmt.Sscanf(r.Lat, "%f", &lat)
	fmt.Sscanf(r.Lon, "%f", &lon)

	loc := &Location{
		Name:    r.Name,
		Lat:     lat,
		Lon:     lon,
		Address: r.DisplayName,
	}

	// Detect airports
	if r.Class == "aeroway" || r.Type == "aerodrome" ||
		strings.Contains(strings.ToLower(r.Name), "arlanda") ||
		strings.Contains(strings.ToLower(r.Name), "bromma") ||
		strings.Contains(strings.ToLower(r.Name), "skavsta") {
		loc.IsAirport = true
	}

	return loc, nil
}

// CalculateRoute calculates the route between two locations using OSRM
func CalculateRoute(from, to *Location) (*Route, error) {
	// OSRM uses lon,lat format
	coords := fmt.Sprintf("%f,%f;%f,%f", from.Lon, from.Lat, to.Lon, to.Lat)
	reqURL := fmt.Sprintf("%s/%s?overview=false", osrmURL, coords)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("routing failed: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Code   string `json:"code"`
		Routes []struct {
			Distance float64 `json:"distance"` // meters
			Duration float64 `json:"duration"` // seconds
		} `json:"routes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode routing response: %w", err)
	}

	if result.Code != "Ok" || len(result.Routes) == 0 {
		return nil, fmt.Errorf("no route found")
	}

	return &Route{
		From:        *from,
		To:          *to,
		DistanceKm:  result.Routes[0].Distance / 1000,
		DurationMin: result.Routes[0].Duration / 60,
	}, nil
}

// CalculateFare calculates the fare for a route
func CalculateFare(route *Route, baseFee, perKm, perHour float64) float64 {
	hours := route.DurationMin / 60
	return baseFee + (route.DistanceKm * perKm) + (hours * perHour)
}

// GetFareEstimates returns fare estimates for all taxi companies
func GetFareEstimates(route *Route) []FareEstimate {
	estimates := []FareEstimate{}

	// Check for fixed airport prices
	isArlanda := strings.Contains(strings.ToLower(route.To.Name), "arlanda") ||
		strings.Contains(strings.ToLower(route.From.Name), "arlanda") ||
		strings.Contains(strings.ToLower(route.To.Address), "arlanda") ||
		strings.Contains(strings.ToLower(route.From.Address), "arlanda")

	// Taxi Stockholm
	tsEstimate := FareEstimate{
		Company:     "Taxi Stockholm",
		BaseFee:     59,
		PerKmRate:   14.90,
		PerHourRate: 565,
		BookingURL:  "https://www.taxistockholm.se/en/booking/",
	}
	tsEstimate.Estimated = CalculateFare(route, tsEstimate.BaseFee, tsEstimate.PerKmRate, tsEstimate.PerHourRate)
	if isArlanda {
		tsEstimate.FixedPrice = 700 // Approximate fixed price
		tsEstimate.IsFixed = true
	}
	estimates = append(estimates, tsEstimate)

	// Taxi Kurir
	tkEstimate := FareEstimate{
		Company:     "Taxi Kurir",
		BaseFee:     55,
		PerKmRate:   14.60,
		PerHourRate: 576,
		BookingURL:  "https://www.taxikurir.se/boka",
	}
	tkEstimate.Estimated = CalculateFare(route, tkEstimate.BaseFee, tkEstimate.PerKmRate, tkEstimate.PerHourRate)
	if isArlanda {
		tkEstimate.FixedPrice = 695
		tkEstimate.IsFixed = true
	}
	estimates = append(estimates, tkEstimate)

	// Uber (approximate rates)
	uberEstimate := FareEstimate{
		Company:     "Uber",
		BaseFee:     30,
		PerKmRate:   10,
		PerHourRate: 120, // ~2 SEK/min
		DeepLink:    GenerateUberDeepLink(route.From.Lat, route.From.Lon, route.To.Lat, route.To.Lon, route.To.Name),
	}
	uberEstimate.Estimated = CalculateFare(route, uberEstimate.BaseFee, uberEstimate.PerKmRate, uberEstimate.PerHourRate)
	estimates = append(estimates, uberEstimate)

	// Bolt (approximate rates, similar to Uber)
	boltEstimate := FareEstimate{
		Company:     "Bolt",
		BaseFee:     25,
		PerKmRate:   9,
		PerHourRate: 108, // ~1.8 SEK/min
		BookingURL:  "https://bolt.eu/",
	}
	boltEstimate.Estimated = CalculateFare(route, boltEstimate.BaseFee, boltEstimate.PerKmRate, boltEstimate.PerHourRate)
	estimates = append(estimates, boltEstimate)

	return estimates
}

// GenerateUberDeepLink creates a deep link for the Uber app
func GenerateUberDeepLink(pickupLat, pickupLon, dropLat, dropLon float64, dropName string) string {
	params := url.Values{}
	params.Set("action", "setPickup")
	params.Set("pickup", "my_location")
	params.Set("dropoff[latitude]", fmt.Sprintf("%f", dropLat))
	params.Set("dropoff[longitude]", fmt.Sprintf("%f", dropLon))
	params.Set("dropoff[nickname]", dropName)

	return "https://m.uber.com/ul/?" + params.Encode()
}

// GenerateGoogleMapsURL creates a Google Maps directions URL
func GenerateGoogleMapsURL(from, to string) string {
	params := url.Values{}
	params.Set("api", "1")
	params.Set("origin", from)
	params.Set("destination", to)
	params.Set("travelmode", "driving")

	return "https://www.google.com/maps/dir/?" + params.Encode()
}

// FormatTaxiSearch formats the taxi search results for display
func FormatTaxiSearch(search TaxiSearch) string {
	var sb strings.Builder

	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	sb.WriteString(fmt.Sprintf(" 🚕 Taxi: %s → %s\n", search.From, search.To))
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	if search.Route != nil {
		sb.WriteString(fmt.Sprintf("  Avstånd:    %.1f km\n", search.Route.DistanceKm))
		sb.WriteString(fmt.Sprintf("  Restid:     %.0f min\n", search.Route.DurationMin))
		sb.WriteString("\n")
	}

	// Fare estimates
	sb.WriteString("  Prisuppskattning:\n")
	sb.WriteString("  ─────────────────────────────────────────────────────────────────\n")

	for _, est := range search.Estimates {
		if est.IsFixed && est.FixedPrice > 0 {
			sb.WriteString(fmt.Sprintf("  🚖 %-15s  ~%4.0f kr  (fast pris: %.0f kr)\n",
				est.Company, est.Estimated, est.FixedPrice))
		} else {
			sb.WriteString(fmt.Sprintf("  🚖 %-15s  ~%4.0f kr\n", est.Company, est.Estimated))
		}
	}

	sb.WriteString("\n")
	sb.WriteString("  📊 Taxameter: grundavgift + kr/km + kr/tim (priserna är ungefärliga)\n")
	sb.WriteString("\n")

	// Booking links
	sb.WriteString("  Boka taxi:\n")
	sb.WriteString("  ─────────────────────────────────────────────────────────────────\n")

	for _, est := range search.Estimates {
		if est.DeepLink != "" {
			sb.WriteString(fmt.Sprintf("  📱 %s (app):\n     %s\n\n", est.Company, est.DeepLink))
		} else if est.BookingURL != "" {
			sb.WriteString(fmt.Sprintf("  🌐 %s:\n     %s\n\n", est.Company, est.BookingURL))
		}
	}

	// Google Maps link
	sb.WriteString("  🗺️  Visa rutt i Google Maps:\n")
	sb.WriteString(fmt.Sprintf("     %s\n\n", GenerateGoogleMapsURL(search.From, search.To)))

	// Tips
	sb.WriteString("  💡 Tips:\n")
	sb.WriteString("     • Jämförpris (10 km, 15 min) finns på bilens dörr\n")
	sb.WriteString("     • Fråga om fast pris innan du stiger in\n")
	sb.WriteString("     • Godkända taxibilar har gula nummerskyltar\n")

	sb.WriteString("\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	return sb.String()
}
