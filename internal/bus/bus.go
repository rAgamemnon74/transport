package bus

import (
	"fmt"
	"net/url"
	"strings"
)

// BusRoute represents a long-distance bus route
type BusRoute struct {
	Operator    string
	From        string
	To          string
	Duration    string // approximate
	PriceFrom   int    // SEK, minimum price
	PriceTo     int    // SEK, typical price
	Frequency   string // e.g., "8-12 per dag"
	BookingURL  string
	HasWifi     bool
	HasPower    bool
	HasToilet   bool
}

// BusSearch contains search parameters and results
type BusSearch struct {
	From      string
	To        string
	FromCity  string
	ToCity    string
	Routes    []BusRoute
	IsAirport bool
}

// City represents a bus station/city
type City struct {
	Name       string
	FlixBusID  string // FlixBus city ID
	VyStop     string // Vy Bus4You stop name
	IsAirport  bool
	AirportCode string
}

// Swedish cities with FlixBus IDs and Vy stops
var Cities = map[string]City{
	// Major cities
	"stockholm": {
		Name:      "Stockholm",
		FlixBusID: "40dfdbe7-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Stockholm",
	},
	"göteborg": {
		Name:      "Göteborg",
		FlixBusID: "40de87a6-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Göteborg",
	},
	"gothenburg": {
		Name:      "Göteborg",
		FlixBusID: "40de87a6-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Göteborg",
	},
	"malmö": {
		Name:      "Malmö",
		FlixBusID: "40de8c24-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Malmö",
	},
	"malmo": {
		Name:      "Malmö",
		FlixBusID: "40de8c24-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Malmö",
	},
	"uppsala": {
		Name:      "Uppsala",
		FlixBusID: "40de9066-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Uppsala",
	},
	"linköping": {
		Name:      "Linköping",
		FlixBusID: "40de8aea-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Linköping",
	},
	"linkoping": {
		Name:      "Linköping",
		FlixBusID: "40de8aea-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Linköping",
	},
	"norrköping": {
		Name:      "Norrköping",
		FlixBusID: "40de8d64-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Norrköping",
	},
	"norrkoping": {
		Name:      "Norrköping",
		FlixBusID: "40de8d64-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Norrköping",
	},
	"jönköping": {
		Name:      "Jönköping",
		FlixBusID: "40de8940-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Jönköping",
	},
	"jonkoping": {
		Name:      "Jönköping",
		FlixBusID: "40de8940-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Jönköping",
	},
	"örebro": {
		Name:      "Örebro",
		FlixBusID: "40de90e8-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Örebro",
	},
	"orebro": {
		Name:      "Örebro",
		FlixBusID: "40de90e8-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Örebro",
	},
	"västerås": {
		Name:      "Västerås",
		FlixBusID: "40de9156-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Västerås",
	},
	"vasteras": {
		Name:      "Västerås",
		FlixBusID: "40de9156-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Västerås",
	},
	"karlstad": {
		Name:      "Karlstad",
		FlixBusID: "40de89b8-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Karlstad",
	},
	"borås": {
		Name:      "Borås",
		FlixBusID: "40de85ee-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Borås",
	},
	"boras": {
		Name:      "Borås",
		FlixBusID: "40de85ee-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Borås",
	},
	"helsingborg": {
		Name:      "Helsingborg",
		FlixBusID: "40de882a-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Helsingborg",
	},
	"lund": {
		Name:      "Lund",
		FlixBusID: "40de8b9e-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Lund",
	},
	"umeå": {
		Name:      "Umeå",
		FlixBusID: "40de8ff8-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Umeå",
	},
	"umea": {
		Name:      "Umeå",
		FlixBusID: "40de8ff8-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Umeå",
	},
	"sundsvall": {
		Name:      "Sundsvall",
		FlixBusID: "40de8efe-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Sundsvall",
	},
	"gävle": {
		Name:      "Gävle",
		FlixBusID: "40de879c-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Gävle",
	},
	"gavle": {
		Name:      "Gävle",
		FlixBusID: "40de879c-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Gävle",
	},
	"kalmar": {
		Name:      "Kalmar",
		FlixBusID: "40de88a2-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Kalmar",
	},
	"växjö": {
		Name:      "Växjö",
		FlixBusID: "40de91c4-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Växjö",
	},
	"vaxjo": {
		Name:      "Växjö",
		FlixBusID: "40de91c4-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Växjö",
	},
	"halmstad": {
		Name:      "Halmstad",
		FlixBusID: "40de87ce-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Halmstad",
	},
	"kristianstad": {
		Name:      "Kristianstad",
		FlixBusID: "40de8a50-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Kristianstad",
	},
	"karlskrona": {
		Name:      "Karlskrona",
		FlixBusID: "40de8a28-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Karlskrona",
	},
	"luleå": {
		Name:      "Luleå",
		FlixBusID: "40de8b76-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Luleå",
	},
	"lulea": {
		Name:      "Luleå",
		FlixBusID: "40de8b76-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Luleå",
	},
	// Northern Sweden / Transit towns
	"ånge": {
		Name:      "Ånge",
		FlixBusID: "",
		VyStop:    "Ånge",
	},
	"ange": {
		Name:      "Ånge",
		FlixBusID: "",
		VyStop:    "Ånge",
	},
	"svenstavik": {
		Name:      "Svenstavik",
		FlixBusID: "",
		VyStop:    "Svenstavik",
	},
	"östersund": {
		Name:      "Östersund",
		FlixBusID: "40de9124-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Östersund",
	},
	"ostersund": {
		Name:      "Östersund",
		FlixBusID: "40de9124-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Östersund",
	},
	"mora": {
		Name:      "Mora",
		FlixBusID: "40de8cf6-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Mora",
	},
	"falun": {
		Name:      "Falun",
		FlixBusID: "40de86ca-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Falun",
	},
	"borlänge": {
		Name:      "Borlänge",
		FlixBusID: "40de85c6-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Borlänge",
	},
	"borlange": {
		Name:      "Borlänge",
		FlixBusID: "40de85c6-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Borlänge",
	},
	// Småland / Shopping
	"ullared": {
		Name:      "Ullared",
		FlixBusID: "",
		VyStop:    "Ullared",
	},
	"värnamo": {
		Name:      "Värnamo",
		FlixBusID: "40de9188-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Värnamo",
	},
	"varnamo": {
		Name:      "Värnamo",
		FlixBusID: "40de9188-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Värnamo",
	},
	"ljungby": {
		Name:      "Ljungby",
		FlixBusID: "40de8afe-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Ljungby",
	},
	"nässjö": {
		Name:      "Nässjö",
		FlixBusID: "40de8d8c-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Nässjö",
	},
	"nassjo": {
		Name:      "Nässjö",
		FlixBusID: "40de8d8c-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Nässjö",
	},
	"vetlanda": {
		Name:      "Vetlanda",
		FlixBusID: "",
		VyStop:    "Vetlanda",
	},
	"eksjö": {
		Name:      "Eksjö",
		FlixBusID: "",
		VyStop:    "Eksjö",
	},
	"eksjo": {
		Name:      "Eksjö",
		FlixBusID: "",
		VyStop:    "Eksjö",
	},
	// Ski resorts
	"åre": {
		Name:      "Åre",
		FlixBusID: "40de9232-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Åre",
	},
	"are": {
		Name:      "Åre",
		FlixBusID: "40de9232-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Åre",
	},
	"sälen": {
		Name:      "Sälen",
		FlixBusID: "40de8e72-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Sälen",
	},
	"salen": {
		Name:      "Sälen",
		FlixBusID: "40de8e72-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Sälen",
	},
	"vemdalen": {
		Name:      "Vemdalen",
		FlixBusID: "0f2869d8-d001-42e3-8f28-df360bbfa313",
		VyStop:    "Vemdalen",
	},
	"idre": {
		Name:      "Idre",
		FlixBusID: "40de88c0-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Idre",
	},
	"funäsdalen": {
		Name:      "Funäsdalen",
		FlixBusID: "40de8756-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Funäsdalen",
	},
	"funasdalen": {
		Name:      "Funäsdalen",
		FlixBusID: "40de8756-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Funäsdalen",
	},
	"trysil": {
		Name:      "Trysil",
		FlixBusID: "40de7e68-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Trysil",
	},
	"hemavan": {
		Name:      "Hemavan",
		FlixBusID: "40de8846-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Hemavan",
	},
	"riksgränsen": {
		Name:      "Riksgränsen",
		FlixBusID: "40de8dc8-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Riksgränsen",
	},
	"riksgransen": {
		Name:      "Riksgränsen",
		FlixBusID: "40de8dc8-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Riksgränsen",
	},
	// International
	"oslo": {
		Name:      "Oslo",
		FlixBusID: "40de7d0a-8646-11e6-9066-549f350fcb0c",
		VyStop:    "Oslo",
	},
	"copenhagen": {
		Name:      "Köpenhamn",
		FlixBusID: "40de5cda-8646-11e6-9066-549f350fcb0c",
		VyStop:    "København",
	},
	"köpenhamn": {
		Name:      "Köpenhamn",
		FlixBusID: "40de5cda-8646-11e6-9066-549f350fcb0c",
		VyStop:    "København",
	},
	"kopenhamn": {
		Name:      "Köpenhamn",
		FlixBusID: "40de5cda-8646-11e6-9066-549f350fcb0c",
		VyStop:    "København",
	},
	// Airports
	"arlanda": {
		Name:        "Stockholm Arlanda Airport",
		FlixBusID:   "40dea650-8646-11e6-9066-549f350fcb0c",
		VyStop:      "Arlanda",
		IsAirport:   true,
		AirportCode: "ARN",
	},
	"landvetter": {
		Name:        "Göteborg Landvetter Airport",
		FlixBusID:   "40dea754-8646-11e6-9066-549f350fcb0c",
		VyStop:      "Landvetter",
		IsAirport:   true,
		AirportCode: "GOT",
	},
}

// LookupCity finds a city by name
func LookupCity(name string) (*City, bool) {
	city, ok := Cities[strings.ToLower(strings.TrimSpace(name))]
	if ok {
		return &city, true
	}
	return nil, false
}

// RouteInfo contains pre-defined route information
type RouteInfo struct {
	Duration  string
	PriceFrom int
	PriceTo   int
	Frequency string
	HasVy     bool // Vy Bus4You operates this route
}

// Common routes with approximate info
var CommonRoutes = map[string]RouteInfo{
	"stockholm-göteborg":    {"4-5 tim", 99, 299, "8-12/dag", true},
	"stockholm-malmö":       {"6-8 tim", 149, 399, "6-10/dag", false},
	"göteborg-malmö":        {"3-4 tim", 99, 249, "8-12/dag", false},
	"stockholm-oslo":        {"6-7 tim", 149, 349, "4-6/dag", true},
	"stockholm-köpenhamn":   {"8-9 tim", 199, 449, "4-6/dag", false},
	"göteborg-köpenhamn":    {"4-5 tim", 149, 299, "4-6/dag", true},
	"stockholm-linköping":   {"2-3 tim", 79, 199, "10-15/dag", true},
	"stockholm-norrköping":  {"2 tim", 79, 179, "10-15/dag", true},
	"stockholm-jönköping":   {"3-4 tim", 99, 249, "6-8/dag", true},
	"stockholm-örebro":      {"2.5 tim", 79, 199, "6-8/dag", true},
	"stockholm-västerås":    {"1-1.5 tim", 59, 149, "8-10/dag", true},
	"stockholm-karlstad":    {"3-4 tim", 99, 249, "4-6/dag", true},
	"stockholm-uppsala":     {"45 min", 49, 99, "Många/dag", false},
	"göteborg-borås":        {"1 tim", 49, 99, "Många/dag", true},
	"malmö-lund":            {"20 min", 39, 59, "Många/dag", false},
	"malmö-helsingborg":     {"1 tim", 49, 99, "Många/dag", false},
	"stockholm-arlanda":     {"45 min", 99, 139, "Var 10 min", false},
	"göteborg-landvetter":   {"30 min", 99, 119, "Var 15 min", false},
	// Ski resorts
	"stockholm-åre":         {"7-8 tim", 299, 599, "2-4/dag", false},
	"stockholm-sälen":       {"5-6 tim", 249, 499, "2-4/dag", false},
	"stockholm-vemdalen":    {"5-6 tim", 249, 499, "2-3/dag", false},
	"stockholm-idre":        {"5 tim", 249, 449, "1-2/dag", false},
	"göteborg-åre":          {"8-9 tim", 349, 649, "1-2/dag", false},
	"göteborg-sälen":        {"5-6 tim", 249, 499, "1-2/dag", false},
	"oslo-trysil":           {"2.5 tim", 149, 299, "3-4/dag", false},
	"stockholm-funäsdalen":  {"6 tim", 279, 529, "1-2/dag", false},
	"stockholm-östersund":   {"6 tim", 249, 499, "3-4/dag", false},
	"stockholm-mora":        {"4 tim", 199, 399, "2-3/dag", false},
	"stockholm-falun":       {"3 tim", 149, 299, "4-6/dag", false},
	"göteborg-ullared":      {"1.5 tim", 99, 199, "4-6/dag", false},
	"stockholm-värnamo":     {"4 tim", 149, 349, "2-3/dag", false},
	"göteborg-värnamo":      {"2 tim", 99, 199, "3-4/dag", false},
	"malmö-värnamo":         {"2.5 tim", 99, 249, "2-3/dag", false},
}

// GetRouteInfo returns info for a route
func GetRouteInfo(from, to string) *RouteInfo {
	key1 := strings.ToLower(from) + "-" + strings.ToLower(to)
	key2 := strings.ToLower(to) + "-" + strings.ToLower(from)

	if info, ok := CommonRoutes[key1]; ok {
		return &info
	}
	if info, ok := CommonRoutes[key2]; ok {
		return &info
	}
	return nil
}

// GenerateFlixBusURL creates a FlixBus search URL
func GenerateFlixBusURL(from, to *City, date string) string {
	baseURL := "https://shop.flixbus.se/search"
	params := url.Values{}
	params.Set("departureCity", from.FlixBusID)
	params.Set("arrivalCity", to.FlixBusID)
	params.Set("route", from.Name+"-"+to.Name)
	if date != "" {
		// Convert YYYY-MM-DD to DD.MM.YYYY
		parts := strings.Split(date, "-")
		if len(parts) == 3 {
			params.Set("rideDate", parts[2]+"."+parts[1]+"."+parts[0])
		}
	}
	params.Set("adult", "1")
	params.Set("_locale", "sv")
	params.Set("departureCountryCode", "SE")
	params.Set("arrivalCountryCode", "SE")

	return baseURL + "?" + params.Encode()
}

// GenerateVyURL creates a Vy Bus4You search URL
func GenerateVyURL(from, to *City) string {
	return "https://www.vy.se/en/traffic-and-routes/buses"
}

// GenerateFlygbussarnaURL creates a Flygbussarna URL for airport routes
func GenerateFlygbussarnaURL(airportCode string) string {
	switch airportCode {
	case "ARN":
		return "https://www.flygbussarna.se/en/arlanda"
	case "GOT":
		return "https://www.flygbussarna.se/en/landvetter"
	case "BMA":
		return "https://www.flygbussarna.se/en/bromma"
	case "MMX":
		return "https://www.flygbussarna.se/en/sturup"
	default:
		return "https://www.flygbussarna.se/en"
	}
}

// GenerateOmioURL creates an Omio aggregator search URL
func GenerateOmioURL(from, to string) string {
	params := url.Values{}
	params.Set("departurePosition", from)
	params.Set("arrivalPosition", to)

	return "https://www.omio.com/search-frontend/results?" + params.Encode()
}

// GetBusRoutes returns available bus routes between two cities
func GetBusRoutes(from, to *City, date string) []BusRoute {
	var routes []BusRoute

	routeInfo := GetRouteInfo(from.Name, to.Name)

	// FlixBus (if both cities have FlixBus IDs)
	if from.FlixBusID != "" && to.FlixBusID != "" {
		flixbus := BusRoute{
			Operator:   "FlixBus",
			From:       from.Name,
			To:         to.Name,
			HasWifi:    true,
			HasPower:   true,
			HasToilet:  true,
			BookingURL: GenerateFlixBusURL(from, to, date),
		}
		if routeInfo != nil {
			flixbus.Duration = routeInfo.Duration
			flixbus.PriceFrom = routeInfo.PriceFrom
			flixbus.PriceTo = routeInfo.PriceTo
			flixbus.Frequency = routeInfo.Frequency
		} else {
			flixbus.Duration = "Varierande"
			flixbus.PriceFrom = 99
			flixbus.PriceTo = 399
			flixbus.Frequency = "Se hemsida"
		}
		routes = append(routes, flixbus)
	}

	// Vy Bus4You (selected routes)
	if routeInfo != nil && routeInfo.HasVy {
		vy := BusRoute{
			Operator:   "Vy Bus4You",
			From:       from.Name,
			To:         to.Name,
			Duration:   routeInfo.Duration,
			PriceFrom:  routeInfo.PriceFrom + 50, // Vy typically slightly more expensive
			PriceTo:    routeInfo.PriceTo + 100,
			Frequency:  routeInfo.Frequency,
			HasWifi:    true,
			HasPower:   true,
			HasToilet:  true,
			BookingURL: GenerateVyURL(from, to),
		}
		routes = append(routes, vy)
	}

	// Flygbussarna for airport routes
	if from.IsAirport {
		flygbuss := BusRoute{
			Operator:   "Flygbussarna",
			From:       from.Name,
			To:         to.Name,
			Duration:   "30-45 min",
			PriceFrom:  99,
			PriceTo:    139,
			Frequency:  "Var 10-15 min",
			HasWifi:    true,
			HasPower:   true,
			HasToilet:  false,
			BookingURL: GenerateFlygbussarnaURL(from.AirportCode),
		}
		routes = append(routes, flygbuss)
	} else if to.IsAirport {
		flygbuss := BusRoute{
			Operator:   "Flygbussarna",
			From:       from.Name,
			To:         to.Name,
			Duration:   "30-45 min",
			PriceFrom:  99,
			PriceTo:    139,
			Frequency:  "Var 10-15 min",
			HasWifi:    true,
			HasPower:   true,
			HasToilet:  false,
			BookingURL: GenerateFlygbussarnaURL(to.AirportCode),
		}
		routes = append(routes, flygbuss)
	}

	return routes
}

// FormatBusSearch formats the bus search results for display
func FormatBusSearch(search BusSearch) string {
	var sb strings.Builder

	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	sb.WriteString(fmt.Sprintf(" 🚌 Buss: %s → %s\n", search.FromCity, search.ToCity))
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	if len(search.Routes) == 0 {
		sb.WriteString("  Inga direktbokningar tillgängliga.\n\n")
		sb.WriteString("  🔍 Sök på aggregatorer:\n")
		sb.WriteString(fmt.Sprintf("     Omio:    %s\n", GenerateOmioURL(search.FromCity, search.ToCity)))
		sb.WriteString(fmt.Sprintf("     FlixBus: https://www.flixbus.se/\n"))
		sb.WriteString("\n")
	} else {
		// Price comparison
		sb.WriteString("  Prisöversikt:\n")
		sb.WriteString("  ─────────────────────────────────────────────────────────────────\n")
		for _, route := range search.Routes {
			amenities := ""
			if route.HasWifi {
				amenities += "📶"
			}
			if route.HasPower {
				amenities += "🔌"
			}
			if route.HasToilet {
				amenities += "🚻"
			}

			sb.WriteString(fmt.Sprintf("  🚌 %-14s  %3d-%d kr  %-10s  %s\n",
				route.Operator,
				route.PriceFrom, route.PriceTo,
				route.Duration,
				amenities))
		}
		sb.WriteString("\n")

		// Route details
		sb.WriteString("  Boka biljett:\n")
		sb.WriteString("  ─────────────────────────────────────────────────────────────────\n")
		for _, route := range search.Routes {
			sb.WriteString(fmt.Sprintf("  🎫 %s:\n", route.Operator))
			sb.WriteString(fmt.Sprintf("     %s\n", route.BookingURL))
			if route.Frequency != "" {
				sb.WriteString(fmt.Sprintf("     Avgångar: %s\n", route.Frequency))
			}
			sb.WriteString("\n")
		}

		// Aggregator
		sb.WriteString("  🔍 Jämför alla operatörer:\n")
		sb.WriteString(fmt.Sprintf("     %s\n\n", GenerateOmioURL(search.FromCity, search.ToCity)))
	}

	// Tips
	sb.WriteString("  💡 Tips:\n")
	sb.WriteString("     • Boka i förväg för lägsta pris\n")
	sb.WriteString("     • FlixBus: ändra bokning upp till 15 min före avgång\n")
	sb.WriteString("     • Vy Bus4You: \"Sveriges nöjdaste kunder\" 12 år i rad\n")

	sb.WriteString("\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	return sb.String()
}
