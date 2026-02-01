package flight

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// OurAirports CSV URL for Swedish airports
	airportsURL = "https://davidmegginson.github.io/ourairports-data/airports.csv"

	// Earth radius in kilometers
	earthRadiusKm = 6371.0
)

// Airport represents an airport from OurAirports database
type Airport struct {
	ID               string
	Ident            string  // ICAO code (e.g., "ESSA")
	Type             string  // large_airport, medium_airport, small_airport, heliport, closed
	Name             string
	Latitude         float64
	Longitude        float64
	Elevation        int     // feet
	Continent        string
	Country          string  // ISO country code
	Region           string  // ISO region code
	Municipality     string
	ScheduledService bool    // Has scheduled commercial flights
	GPSCode          string
	IATACode         string  // e.g., "ARN"
	LocalCode        string
	HomeLink         string
	WikipediaLink    string
	Keywords         string
	DistanceKm       float64 // Calculated distance from reference point
}

// NearbyResult contains an airport and its distance
type NearbyResult struct {
	Airport    Airport
	DistanceKm float64
}

// FetchSwedishAirports fetches all Swedish airports from OurAirports
func FetchSwedishAirports() ([]Airport, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(airportsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch airports: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	return parseAirportsCSV(resp.Body, "SE")
}

// parseAirportsCSV parses the OurAirports CSV format
func parseAirportsCSV(r io.Reader, countryFilter string) ([]Airport, error) {
	reader := csv.NewReader(r)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Build column index map
	colIndex := make(map[string]int)
	for i, name := range header {
		colIndex[name] = i
	}

	// Required columns
	required := []string{"id", "ident", "type", "name", "latitude_deg", "longitude_deg", "iso_country", "scheduled_service"}
	for _, col := range required {
		if _, ok := colIndex[col]; !ok {
			return nil, fmt.Errorf("missing required column: %s", col)
		}
	}

	var airports []Airport

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		// Filter by country
		country := getField(record, colIndex, "iso_country")
		if countryFilter != "" && country != countryFilter {
			continue
		}

		lat, _ := strconv.ParseFloat(getField(record, colIndex, "latitude_deg"), 64)
		lon, _ := strconv.ParseFloat(getField(record, colIndex, "longitude_deg"), 64)
		elev, _ := strconv.Atoi(getField(record, colIndex, "elevation_ft"))

		airport := Airport{
			ID:               getField(record, colIndex, "id"),
			Ident:            getField(record, colIndex, "ident"),
			Type:             getField(record, colIndex, "type"),
			Name:             getField(record, colIndex, "name"),
			Latitude:         lat,
			Longitude:        lon,
			Elevation:        elev,
			Continent:        getField(record, colIndex, "continent"),
			Country:          country,
			Region:           getField(record, colIndex, "iso_region"),
			Municipality:     getField(record, colIndex, "municipality"),
			ScheduledService: getField(record, colIndex, "scheduled_service") == "yes",
			GPSCode:          getField(record, colIndex, "gps_code"),
			IATACode:         getField(record, colIndex, "iata_code"),
			LocalCode:        getField(record, colIndex, "local_code"),
			HomeLink:         getField(record, colIndex, "home_link"),
			WikipediaLink:    getField(record, colIndex, "wikipedia_link"),
			Keywords:         getField(record, colIndex, "keywords"),
		}

		airports = append(airports, airport)
	}

	return airports, nil
}

func getField(record []string, colIndex map[string]int, name string) string {
	if idx, ok := colIndex[name]; ok && idx < len(record) {
		return record[idx]
	}
	return ""
}

// FindNearbyAirports finds airports within the specified radius
func FindNearbyAirports(airports []Airport, lat, lon, radiusKm float64, scheduledOnly bool) []Airport {
	var nearby []Airport

	for _, airport := range airports {
		// Filter closed airports
		if airport.Type == "closed" {
			continue
		}

		// Filter by scheduled service if requested
		if scheduledOnly && !airport.ScheduledService {
			continue
		}

		// Calculate distance
		distance := haversineDistance(lat, lon, airport.Latitude, airport.Longitude)

		if distance <= radiusKm {
			airport.DistanceKm = distance
			nearby = append(nearby, airport)
		}
	}

	// Sort by distance
	sort.Slice(nearby, func(i, j int) bool {
		return nearby[i].DistanceKm < nearby[j].DistanceKm
	})

	return nearby
}

// haversineDistance calculates the distance between two coordinates in km
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert to radians
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	// Haversine formula
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// Stockholm coordinates (central Stockholm)
var StockholmCoords = struct {
	Lat float64
	Lon float64
}{
	Lat: 59.3293,
	Lon: 18.0686,
}

// CommonLocations maps location names to coordinates
var CommonLocations = map[string]struct {
	Lat float64
	Lon float64
}{
	"stockholm":  {59.3293, 18.0686},
	"göteborg":   {57.7089, 11.9746},
	"gothenburg": {57.7089, 11.9746},
	"malmö":      {55.6050, 13.0038},
	"malmo":      {55.6050, 13.0038},
	"uppsala":    {59.8586, 17.6389},
	"linköping":  {58.4108, 15.6214},
	"linkoping":  {58.4108, 15.6214},
	"örebro":     {59.2753, 15.2134},
	"orebro":     {59.2753, 15.2134},
	"västerås":   {59.6099, 16.5448},
	"vasteras":   {59.6099, 16.5448},
	"norrköping": {58.5877, 16.1924},
	"norrkoping": {58.5877, 16.1924},
	"lund":       {55.7047, 13.1910},
	"umeå":       {63.8258, 20.2630},
	"umea":       {63.8258, 20.2630},
	"jönköping":  {57.7826, 14.1618},
	"jonkoping":  {57.7826, 14.1618},
	"luleå":      {65.5848, 22.1547},
	"lulea":      {65.5848, 22.1547},
	"kiruna":     {67.8558, 20.2253},
	"sundsvall":  {62.3908, 17.3069},
	"gävle":      {60.6749, 17.1413},
	"gavle":      {60.6749, 17.1413},
	"karlstad":   {59.3793, 13.5036},
	"växjö":      {56.8777, 14.8091},
	"vaxjo":      {56.8777, 14.8091},
	"halmstad":   {56.6745, 12.8578},
	"kalmar":     {56.6634, 16.3566},
	"visby":      {57.6348, 18.2948},
	"åre":        {63.3988, 13.0814},
	"are":        {63.3988, 13.0814},
}

// GetCoordinates returns coordinates for a location name
func GetCoordinates(location string) (lat, lon float64, found bool) {
	coords, ok := CommonLocations[strings.ToLower(location)]
	if ok {
		return coords.Lat, coords.Lon, true
	}
	return 0, 0, false
}

// GetAirportTypeLabel returns a Swedish label for airport type
func GetAirportTypeLabel(airportType string) string {
	switch airportType {
	case "large_airport":
		return "Stor flygplats"
	case "medium_airport":
		return "Mellanstor flygplats"
	case "small_airport":
		return "Liten flygplats"
	case "heliport":
		return "Heliport"
	case "seaplane_base":
		return "Sjöflygplats"
	default:
		return airportType
	}
}

// GetAirportIcon returns an icon for the airport type
func GetAirportIcon(airport Airport) string {
	if airport.ScheduledService {
		return "✈️"
	}
	switch airport.Type {
	case "heliport":
		return "🚁"
	case "seaplane_base":
		return "🛩️"
	default:
		return "🛫"
	}
}

// FormatNearbyAirports formats the nearby airports for display
func FormatNearbyAirports(location string, radiusKm float64, airports []Airport) string {
	var sb strings.Builder

	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	sb.WriteString(fmt.Sprintf(" ✈️  Flygplatser inom %.0f km från %s\n", radiusKm, location))
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	if len(airports) == 0 {
		sb.WriteString("  Inga flygplatser hittades inom angivet avstånd.\n")
	} else {
		// Group by scheduled service
		var scheduled, other []Airport
		for _, a := range airports {
			if a.ScheduledService {
				scheduled = append(scheduled, a)
			} else {
				other = append(other, a)
			}
		}

		if len(scheduled) > 0 {
			sb.WriteString("  Med reguljär trafik:\n")
			sb.WriteString("  ─────────────────────────────────────────────────────────────────\n")
			for _, a := range scheduled {
				iata := a.IATACode
				if iata == "" {
					iata = a.Ident
				}
				sb.WriteString(fmt.Sprintf("  %s %-4s  %-35s  %5.0f km\n",
					GetAirportIcon(a),
					iata,
					truncateString(a.Name, 35),
					a.DistanceKm))
				if a.Municipality != "" {
					sb.WriteString(fmt.Sprintf("          %s\n", a.Municipality))
				}
			}
			sb.WriteString("\n")
		}

		if len(other) > 0 {
			sb.WriteString("  Övriga flygplatser:\n")
			sb.WriteString("  ─────────────────────────────────────────────────────────────────\n")
			for _, a := range other {
				code := a.IATACode
				if code == "" {
					code = a.Ident
				}
				sb.WriteString(fmt.Sprintf("  %s %-6s  %-33s  %5.0f km  (%s)\n",
					GetAirportIcon(a),
					code,
					truncateString(a.Name, 33),
					a.DistanceKm,
					GetAirportTypeLabel(a.Type)))
			}
		}
	}

	sb.WriteString("\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	return sb.String()
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
