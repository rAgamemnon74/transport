package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"transport/internal/api"
	"transport/internal/bus"
	"transport/internal/car"
	"transport/internal/display"
	"transport/internal/flight"
	"transport/internal/resrobot"
	"transport/internal/taxi"
)

var (
	version = "0.1.0"
)

func main() {
	if len(os.Args) > 1 {
		cmd := strings.ToLower(os.Args[1])

		if isNextCommand(cmd) {
			runNextCommand(os.Args[2:])
			return
		}

		if isCarCommand(cmd) {
			runCarCommand(os.Args[2:])
			return
		}

		if isFlightCommand(cmd) {
			runFlightCommand(os.Args[2:])
			return
		}

		if isFlyCommand(cmd) {
			runFlyCommand(os.Args[2:])
			return
		}

		if isTaxiCommand(cmd) {
			runTaxiCommand(os.Args[2:])
			return
		}

		if isBusCommand(cmd) {
			runBusCommand(os.Args[2:])
			return
		}
	}

	runTripCommand()
}

// isNextCommand checks if the argument is a "next" command (English or Swedish)
func isNextCommand(arg string) bool {
	switch strings.ToLower(arg) {
	case "next", "nästa":
		return true
	}
	return false
}

// isCarCommand checks if the argument is a "car" command (English or Swedish)
func isCarCommand(arg string) bool {
	switch strings.ToLower(arg) {
	case "car", "bil":
		return true
	}
	return false
}

// isFlightCommand checks if the argument is a "flight" command (English or Swedish)
func isFlightCommand(arg string) bool {
	switch strings.ToLower(arg) {
	case "flight", "flights", "flyg", "flygplats", "airport", "airports":
		return true
	}
	return false
}

// isFlyCommand checks if the argument is a "fly" command for flight search
func isFlyCommand(arg string) bool {
	switch strings.ToLower(arg) {
	case "fly", "ffly", "flyga":
		return true
	}
	return false
}

// isTaxiCommand checks if the argument is a "taxi" command
func isTaxiCommand(arg string) bool {
	switch strings.ToLower(arg) {
	case "taxi", "cab", "uber", "bolt":
		return true
	}
	return false
}

// isBusCommand checks if the argument is a "bus" command for long-distance buses
func isBusCommand(arg string) bool {
	switch strings.ToLower(arg) {
	case "buss", "bus", "flixbus", "vy":
		return true
	}
	return false
}

// normalizeMode converts Swedish/English transport mode names to API format
// Returns empty string if mode is invalid
func normalizeMode(mode string) string {
	switch mode {
	// Bus
	case "bus", "buss":
		return "bus"
	// Metro
	case "metro", "tunnelbana", "t-bana", "tbana":
		return "metro"
	// Train
	case "train", "tåg", "tag", "pendeltåg", "pendeltag":
		return "train"
	// Tram
	case "tram", "spårvagn", "sparvagn":
		return "tram"
	// Ship/Ferry
	case "ship", "båt", "bat", "färja", "farja", "ferry":
		return "ship"
	default:
		return ""
	}
}

func runCarCommand(args []string) {
	fs := flag.NewFlagSet("car", flag.ExitOnError)

	var (
		distance  float64
		startFuel float64
	)

	fs.Float64Var(&distance, "d", 0, "Distance in km")
	fs.Float64Var(&distance, "distance", 0, "Distance in km")
	fs.Float64Var(&startFuel, "f", 100, "Starting fuel level in % (default: 100 = full tank)")
	fs.Float64Var(&startFuel, "fuel", 100, "Starting fuel level in %")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Transport - Car directions / Bil vägbeskrivning\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  transport car|bil <from> <to>\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  from    Starting address/location\n")
		fmt.Fprintf(os.Stderr, "  to      Destination address/location\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  transport bil \"Styresman Sanders väg, Bromma\" \"Åre\"\n")
		fmt.Fprintf(os.Stderr, "  transport bil -d 620 \"Stockholm\" \"Åre\"\n")
		fmt.Fprintf(os.Stderr, "  transport bil -d 620 -f 50 \"Stockholm\" \"Åre\"  # Start with half tank\n")
		fmt.Fprintf(os.Stderr, "  transport bil \"Bromma\" \"Arlanda\"\n\n")
		fmt.Fprintf(os.Stderr, "Vehicle: VW Tiguan Allspace 2018 (Diesel)\n")
		fmt.Fprintf(os.Stderr, "  Tank:  58 L (range ~700 km)\n")
		fmt.Fprintf(os.Stderr, "  Short trips (<20km): 9.0 L/100km\n")
		fmt.Fprintf(os.Stderr, "  Long trips (≥20km):  7.0 L/100km\n")
	}

	fs.Parse(args)
	posArgs := fs.Args()

	if len(posArgs) < 2 {
		fs.Usage()
		os.Exit(1)
	}

	from := posArgs[0]
	to := posArgs[1]

	profile := car.DefaultProfile()
	output := car.FormatCarTrip(from, to, distance, startFuel, profile)
	fmt.Print(output)
}

func runFlyCommand(args []string) {
	fs := flag.NewFlagSet("fly", flag.ExitOnError)

	var (
		dateFlag    string
		returnFlag  string
		showPrivate bool
	)

	fs.StringVar(&dateFlag, "d", "", "Departure date (YYYY-MM-DD)")
	fs.StringVar(&dateFlag, "date", "", "Departure date (YYYY-MM-DD)")
	fs.StringVar(&returnFlag, "r", "", "Return date (YYYY-MM-DD) for round-trip")
	fs.StringVar(&returnFlag, "return", "", "Return date (YYYY-MM-DD)")
	fs.BoolVar(&showPrivate, "p", false, "Show private jet & helicopter options")
	fs.BoolVar(&showPrivate, "private", false, "Show private jet & helicopter options")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Transport - Flight search / Sök flygresor\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  transport fly [options] from <origin> to <destination>\n")
		fmt.Fprintf(os.Stderr, "  transport fly [options] <origin> <destination>\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  origin       Departure city or IATA code (e.g., Stockholm, ARN)\n")
		fmt.Fprintf(os.Stderr, "  destination  Arrival city or IATA code (e.g., Vilnius, VNO)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  transport fly from stockholm to vilnius\n")
		fmt.Fprintf(os.Stderr, "  transport fly stockholm vilnius\n")
		fmt.Fprintf(os.Stderr, "  transport fly -d 2025-03-15 from göteborg to barcelona\n")
		fmt.Fprintf(os.Stderr, "  transport fly -d 2025-06-01 -r 2025-06-08 ARN BCN\n")
		fmt.Fprintf(os.Stderr, "  transport fly -p from bromma to visby     # Private jets & heli\n")
		fmt.Fprintf(os.Stderr, "  transport flyga från malmö till london\n")
	}

	fs.Parse(args)
	posArgs := fs.Args()

	if len(posArgs) < 2 {
		fs.Usage()
		os.Exit(1)
	}

	// Parse origin and destination from args
	// Supports: "from X to Y", "från X till Y", or just "X Y"
	origin, dest := parseFlightRoute(posArgs)

	if origin == "" || dest == "" {
		fmt.Fprintf(os.Stderr, "Error: Could not parse origin and destination\n")
		fmt.Fprintf(os.Stderr, "Usage: transport fly from <origin> to <destination>\n")
		os.Exit(1)
	}

	// Look up IATA codes
	originCode, originFound := flight.LookupAirportCode(origin)
	if !originFound {
		fmt.Fprintf(os.Stderr, "Error: Unknown airport/city '%s'\n", origin)
		fmt.Fprintf(os.Stderr, "Use a city name (e.g., Stockholm) or IATA code (e.g., ARN)\n")
		os.Exit(1)
	}

	destCode, destFound := flight.LookupAirportCode(dest)
	if !destFound {
		fmt.Fprintf(os.Stderr, "Error: Unknown airport/city '%s'\n", dest)
		fmt.Fprintf(os.Stderr, "Use a city name (e.g., Vilnius) or IATA code (e.g., VNO)\n")
		os.Exit(1)
	}

	// Build search
	search := flight.FlightSearch{
		OriginCode:      originCode,
		OriginCity:      capitalizeFirst(origin),
		DestinationCode: destCode,
		DestinationCity: capitalizeFirst(dest),
		Passengers:      1,
		ShowPrivate:     showPrivate,
	}

	// Parse dates
	if dateFlag != "" {
		parsed, err := time.Parse("2006-01-02", dateFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid date format '%s' (use YYYY-MM-DD)\n", dateFlag)
			os.Exit(1)
		}
		search.DepartureDate = parsed
	}

	if returnFlag != "" {
		parsed, err := time.Parse("2006-01-02", returnFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid return date format '%s' (use YYYY-MM-DD)\n", returnFlag)
			os.Exit(1)
		}
		search.ReturnDate = parsed
	}

	// Format and display
	output := flight.FormatFlightSearch(search)
	fmt.Print(output)
}

func runTaxiCommand(args []string) {
	fs := flag.NewFlagSet("taxi", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Transport - Taxi fare estimation / Taxipris\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  transport taxi [options] <from> <to>\n")
		fmt.Fprintf(os.Stderr, "  transport taxi [options] from <origin> to <destination>\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  from    Pickup location (address or place name)\n")
		fmt.Fprintf(os.Stderr, "  to      Dropoff location (address or place name)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  transport taxi Slussen Arlanda\n")
		fmt.Fprintf(os.Stderr, "  transport taxi from T-Centralen to Bromma Airport\n")
		fmt.Fprintf(os.Stderr, "  transport taxi \"Kungsgatan 1\" \"Arlanda Terminal 5\"\n")
	}

	fs.Parse(args)
	posArgs := fs.Args()

	if len(posArgs) < 2 {
		fs.Usage()
		os.Exit(1)
	}

	// Parse origin and destination (reuse flight route parsing)
	from, to := parseFlightRoute(posArgs)
	if from == "" || to == "" {
		// Simple two-argument format
		from = posArgs[0]
		to = posArgs[1]
	}

	fmt.Fprintf(os.Stderr, "Söker rutt från %s till %s...\n", from, to)

	// Geocode locations
	fromLoc, err := taxi.Geocode(from)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: kunde inte hitta '%s': %v\n", from, err)
		os.Exit(1)
	}

	toLoc, err := taxi.Geocode(to)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: kunde inte hitta '%s': %v\n", to, err)
		os.Exit(1)
	}

	// Calculate route
	route, err := taxi.CalculateRoute(fromLoc, toLoc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: kunde inte beräkna rutt: %v\n", err)
		os.Exit(1)
	}

	// Get fare estimates
	estimates := taxi.GetFareEstimates(route)

	// Build search result
	search := taxi.TaxiSearch{
		From:      from,
		To:        to,
		Route:     route,
		Estimates: estimates,
	}

	// Format and display
	output := taxi.FormatTaxiSearch(search)
	fmt.Print(output)
}

func runBusCommand(args []string) {
	fs := flag.NewFlagSet("buss", flag.ExitOnError)

	var dateFlag string

	fs.StringVar(&dateFlag, "d", "", "Departure date (YYYY-MM-DD)")
	fs.StringVar(&dateFlag, "date", "", "Departure date (YYYY-MM-DD)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Transport - Long-distance bus / Långfärdsbuss\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  transport buss [options] <from> <to>\n")
		fmt.Fprintf(os.Stderr, "  transport buss [options] from <origin> to <destination>\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  from    Origin city (e.g., Stockholm, Göteborg, Malmö)\n")
		fmt.Fprintf(os.Stderr, "  to      Destination city\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nOperators:\n")
		fmt.Fprintf(os.Stderr, "  FlixBus      - Europe's largest bus network\n")
		fmt.Fprintf(os.Stderr, "  Vy Bus4You   - Premium Swedish/Norwegian routes\n")
		fmt.Fprintf(os.Stderr, "  Flygbussarna - Airport buses (Arlanda, Landvetter, etc.)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  transport buss Stockholm Göteborg\n")
		fmt.Fprintf(os.Stderr, "  transport buss from Malmö to Oslo\n")
		fmt.Fprintf(os.Stderr, "  transport buss -d 2026-03-15 Stockholm Arlanda\n")
		fmt.Fprintf(os.Stderr, "  transport buss Göteborg Köpenhamn\n")
	}

	fs.Parse(args)
	posArgs := fs.Args()

	if len(posArgs) < 2 {
		fs.Usage()
		os.Exit(1)
	}

	// Parse origin and destination (reuse flight route parsing)
	from, to := parseFlightRoute(posArgs)
	if from == "" || to == "" {
		// Simple two-argument format
		from = posArgs[0]
		to = posArgs[1]
	}

	// Look up cities
	fromCity, fromFound := bus.LookupCity(from)
	if !fromFound {
		fmt.Fprintf(os.Stderr, "Error: Okänd stad '%s'\n", from)
		fmt.Fprintf(os.Stderr, "Kända städer: Stockholm, Göteborg, Malmö, Uppsala, Linköping, etc.\n")
		os.Exit(1)
	}

	toCity, toFound := bus.LookupCity(to)
	if !toFound {
		fmt.Fprintf(os.Stderr, "Error: Okänd stad '%s'\n", to)
		fmt.Fprintf(os.Stderr, "Kända städer: Stockholm, Göteborg, Malmö, Uppsala, Linköping, etc.\n")
		os.Exit(1)
	}

	// Get bus routes
	routes := bus.GetBusRoutes(fromCity, toCity, dateFlag)

	// Build search result
	search := bus.BusSearch{
		From:      from,
		To:        to,
		FromCity:  fromCity.Name,
		ToCity:    toCity.Name,
		Routes:    routes,
		IsAirport: fromCity.IsAirport || toCity.IsAirport,
	}

	// Format and display
	output := bus.FormatBusSearch(search)
	fmt.Print(output)
}

// parseFlightRoute extracts origin and destination from args
// Supports: "from X to Y", "från X till Y", or just "X Y"
func parseFlightRoute(args []string) (origin, dest string) {
	if len(args) < 2 {
		return "", ""
	}

	// Look for "from ... to ..." or "från ... till ..." pattern
	fromIdx := -1
	toIdx := -1

	for i, arg := range args {
		lower := strings.ToLower(arg)
		if lower == "from" || lower == "från" {
			fromIdx = i
		}
		if lower == "to" || lower == "till" {
			toIdx = i
		}
	}

	if fromIdx >= 0 && toIdx > fromIdx {
		// Extract words between "from" and "to"
		originParts := args[fromIdx+1 : toIdx]
		origin = strings.Join(originParts, " ")

		// Extract words after "to"
		destParts := args[toIdx+1:]
		dest = strings.Join(destParts, " ")

		return origin, dest
	}

	// Simple "X Y" format
	if len(args) == 2 {
		return args[0], args[1]
	}

	// Try to split on common words
	return args[0], args[len(args)-1]
}

func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	// If it looks like an IATA code (3 uppercase letters), keep as-is
	if len(s) == 3 && strings.ToUpper(s) == s {
		return s
	}
	// Handle multi-word strings
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

func runFlightCommand(args []string) {
	fs := flag.NewFlagSet("flight", flag.ExitOnError)

	var (
		radius        float64
		scheduledOnly bool
	)

	fs.Float64Var(&radius, "r", 100, "Search radius in km (default: 100)")
	fs.Float64Var(&radius, "radius", 100, "Search radius in km")
	fs.BoolVar(&scheduledOnly, "s", false, "Only show airports with scheduled service")
	fs.BoolVar(&scheduledOnly, "scheduled", false, "Only show airports with scheduled service")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Transport - Nearby airports / Flygplatser i närheten\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  transport flight|flyg [options] [location]\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  location    City name (default: Stockholm)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  transport flyg                       # Airports near Stockholm\n")
		fmt.Fprintf(os.Stderr, "  transport flyg -r 150 Göteborg       # 150km radius from Gothenburg\n")
		fmt.Fprintf(os.Stderr, "  transport flyg -s Malmö              # Only scheduled service\n")
		fmt.Fprintf(os.Stderr, "  transport flight -r 200 Kiruna       # 200km radius from Kiruna\n")
	}

	fs.Parse(args)
	posArgs := fs.Args()

	// Default location is Stockholm
	location := "Stockholm"
	if len(posArgs) >= 1 {
		location = posArgs[0]
	}

	// Get coordinates for location
	lat, lon, found := flight.GetCoordinates(location)
	if !found {
		fmt.Fprintf(os.Stderr, "Error: Unknown location '%s'\n", location)
		fmt.Fprintf(os.Stderr, "Known locations: Stockholm, Göteborg, Malmö, Uppsala, Linköping, Örebro, Västerås,\n")
		fmt.Fprintf(os.Stderr, "  Norrköping, Lund, Umeå, Jönköping, Luleå, Kiruna, Sundsvall, Gävle, Karlstad,\n")
		fmt.Fprintf(os.Stderr, "  Växjö, Halmstad, Kalmar, Visby, Åre\n")
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Hämtar flygplatsdata...\n")

	// Fetch Swedish airports
	airports, err := flight.FetchSwedishAirports()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Find nearby airports
	nearby := flight.FindNearbyAirports(airports, lat, lon, radius, scheduledOnly)

	// Format and display
	output := flight.FormatNearbyAirports(location, radius, nearby)
	fmt.Print(output)
}

func runNextCommand(args []string) {
	fs := flag.NewFlagSet("next", flag.ExitOnError)

	var (
		count int
		lang  string
	)

	fs.IntVar(&count, "n", 3, "Number of departures to show")
	fs.StringVar(&lang, "l", "sv", "Language (sv/en)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Transport - Next departures / Nästa avgång\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  transport next|nästa <mode> <location> [towards]\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  mode      Transport mode:\n")
		fmt.Fprintf(os.Stderr, "              bus, buss\n")
		fmt.Fprintf(os.Stderr, "              metro, tunnelbana, t-bana\n")
		fmt.Fprintf(os.Stderr, "              train, tåg, pendeltåg\n")
		fmt.Fprintf(os.Stderr, "              tram, spårvagn\n")
		fmt.Fprintf(os.Stderr, "              ship, båt, färja\n")
		fmt.Fprintf(os.Stderr, "  location  Stop/station name to depart from\n")
		fmt.Fprintf(os.Stderr, "  towards   (Optional) Destination to filter by\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  transport nästa buss Odenplan\n")
		fmt.Fprintf(os.Stderr, "  transport nästa buss \"Spånga station\" Brommaplan\n")
		fmt.Fprintf(os.Stderr, "  transport nästa tunnelbana Slussen\n")
		fmt.Fprintf(os.Stderr, "  transport nästa tåg \"Stockholm Central\"\n")
		fmt.Fprintf(os.Stderr, "  transport next -n 5 bus Odenplan\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}

	fs.Parse(args)
	posArgs := fs.Args()

	if len(posArgs) < 2 {
		fs.Usage()
		os.Exit(1)
	}

	mode := posArgs[0]
	location := posArgs[1]
	towards := ""
	if len(posArgs) >= 3 {
		towards = posArgs[2]
	}

	// Normalize and validate mode (supports Swedish and English)
	modeLower := normalizeMode(strings.ToLower(mode))
	if modeLower == "" {
		fmt.Fprintf(os.Stderr, "Error: invalid transport mode '%s'\n", mode)
		fmt.Fprintf(os.Stderr, "Valid modes: bus/buss, metro/tunnelbana/t-bana, train/tåg, tram/spårvagn, ship/båt/färja\n")
		os.Exit(1)
	}

	// Create API client
	client := api.NewClient()

	if towards != "" {
		fmt.Fprintf(os.Stderr, "Söker %s från %s mot %s...\n", modeLower, location, towards)
	} else {
		fmt.Fprintf(os.Stderr, "Söker %s från %s...\n", modeLower, location)
	}

	departures, site, err := client.GetNextDepartures(location, modeLower, towards, count)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Format and display
	formatter := display.NewFormatter(lang)
	output := formatter.FormatDepartures(site, modeLower, towards, departures)
	fmt.Print(output)
}

func runTripCommand() {
	// Flags
	var (
		timeFlag    string
		dateFlag    string
		arriveBy    bool
		maxChanges  int
		numResults  int
		lang        string
		showVersion bool
		jsonOutput  bool
		nationwide  bool
	)

	flag.StringVar(&timeFlag, "time", "", "Departure time (HH:MM)")
	flag.StringVar(&timeFlag, "t", "", "Departure time (HH:MM) (shorthand)")
	flag.StringVar(&dateFlag, "date", "", "Departure date (YYYY-MM-DD)")
	flag.StringVar(&dateFlag, "d", "", "Departure date (YYYY-MM-DD) (shorthand)")
	flag.BoolVar(&arriveBy, "arrive", false, "Search by arrival time instead of departure")
	flag.BoolVar(&arriveBy, "a", false, "Search by arrival time (shorthand)")
	flag.IntVar(&maxChanges, "changes", -1, "Maximum number of changes (0-9, -1 for unlimited)")
	flag.IntVar(&maxChanges, "c", -1, "Maximum changes (shorthand)")
	flag.IntVar(&numResults, "results", 3, "Number of results (1-6)")
	flag.IntVar(&numResults, "n", 3, "Number of results (shorthand)")
	flag.StringVar(&lang, "lang", "sv", "Language (sv/en)")
	flag.StringVar(&lang, "l", "sv", "Language (shorthand)")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showVersion, "v", false, "Show version (shorthand)")
	flag.BoolVar(&jsonOutput, "json", false, "Output raw JSON")
	flag.BoolVar(&jsonOutput, "j", false, "Output JSON (shorthand)")
	flag.BoolVar(&nationwide, "se", false, "Search all of Sweden (ResRobot)")
	flag.BoolVar(&nationwide, "sweden", false, "Search all of Sweden (ResRobot)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Transport - Sweden public transport planner\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  transport [options] <destination>\n")
		fmt.Fprintf(os.Stderr, "  transport [options] <origin> <destination>\n")
		fmt.Fprintf(os.Stderr, "  transport next|nästa <mode> <location> [towards]\n")
		fmt.Fprintf(os.Stderr, "  transport car|bil <from> <to>\n")
		fmt.Fprintf(os.Stderr, "  transport fly|flyga from <origin> to <destination>\n")
		fmt.Fprintf(os.Stderr, "  transport flight|flyg [location]\n")
		fmt.Fprintf(os.Stderr, "  transport taxi <from> <to>\n")
		fmt.Fprintf(os.Stderr, "  transport buss <from> <to>\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  (default)    Plan a trip between two locations (public transport)\n")
		fmt.Fprintf(os.Stderr, "  next, nästa  Show next departures (see: transport next --help)\n")
		fmt.Fprintf(os.Stderr, "  car, bil     Car directions with fuel calculation\n")
		fmt.Fprintf(os.Stderr, "  fly, flyga   Search for flights (booking links)\n")
		fmt.Fprintf(os.Stderr, "  flight, flyg Find nearby airports\n")
		fmt.Fprintf(os.Stderr, "  taxi         Taxi fare estimation & booking\n")
		fmt.Fprintf(os.Stderr, "  buss         Long-distance buses (FlixBus, Vy, Flygbussarna)\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  transport Odenplan                           # Stockholm (SL)\n")
		fmt.Fprintf(os.Stderr, "  transport Slussen Odenplan                   # Stockholm (SL)\n")
		fmt.Fprintf(os.Stderr, "  transport -se Sundsvall Ånge                 # Nationwide (ResRobot)\n")
		fmt.Fprintf(os.Stderr, "  transport -se Göteborg \"Stockholm Central\"   # Nationwide (ResRobot)\n")
		fmt.Fprintf(os.Stderr, "  transport -t 08:30 Slussen T-Centralen\n")
		fmt.Fprintf(os.Stderr, "  transport nästa buss Odenplan\n")
		fmt.Fprintf(os.Stderr, "  transport bil \"Stockholm\" \"Åre\"\n")
		fmt.Fprintf(os.Stderr, "  transport fly from stockholm to vilnius\n")
		fmt.Fprintf(os.Stderr, "  transport flyg Göteborg\n")
		fmt.Fprintf(os.Stderr, "  transport taxi Slussen Arlanda\n")
		fmt.Fprintf(os.Stderr, "  transport buss Stockholm Göteborg\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nNationwide search (-se) requires RESROBOT_API_KEY.\n")
		fmt.Fprintf(os.Stderr, "Get a free key at: https://www.trafiklab.se/api/trafiklab-apis/resrobot-v21/\n")
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("transport version %s\n", version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	var origin, dest string
	if len(args) == 1 {
		// Only destination provided, use default origin
		origin = getDefaultLocation()
		dest = args[0]
	} else {
		origin = args[0]
		dest = args[1]
	}

	// Parse time
	searchTime := time.Now()
	if dateFlag != "" {
		parsed, err := time.Parse("2006-01-02", dateFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid date format '%s' (use YYYY-MM-DD)\n", dateFlag)
			os.Exit(1)
		}
		searchTime = time.Date(parsed.Year(), parsed.Month(), parsed.Day(),
			searchTime.Hour(), searchTime.Minute(), 0, 0, searchTime.Location())
	}

	if timeFlag != "" {
		parsed, err := time.Parse("15:04", timeFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid time format '%s' (use HH:MM)\n", timeFlag)
			os.Exit(1)
		}
		searchTime = time.Date(searchTime.Year(), searchTime.Month(), searchTime.Day(),
			parsed.Hour(), parsed.Minute(), 0, 0, searchTime.Location())
	}

	// Use ResRobot for nationwide search
	if nationwide {
		runResRobotSearch(origin, dest, searchTime, arriveBy, numResults)
		return
	}

	// Use SL for Stockholm region
	opts := api.DefaultTripOptions()
	opts.NumResults = numResults
	opts.MaxChanges = maxChanges
	opts.Language = lang
	opts.ArriveBy = arriveBy
	opts.Time = searchTime

	// Create API client
	client := api.NewClient()

	// Search for stops
	fmt.Fprintf(os.Stderr, "Söker resor från %s till %s...\n", origin, dest)

	journeys, err := client.PlanTripByName(origin, dest, opts)
	if err != nil {
		handleError(err, origin, dest, client)
		os.Exit(1)
	}

	if len(journeys) == 0 {
		fmt.Fprintln(os.Stderr, "Inga resor hittades.")
		os.Exit(1)
	}

	// Format and display results
	formatter := display.NewFormatter(lang)
	output := formatter.FormatJourneys(origin, dest, journeys)
	fmt.Print(output)
}

// runResRobotSearch performs a nationwide search using ResRobot
func runResRobotSearch(origin, dest string, searchTime time.Time, arriveBy bool, numResults int) {
	client := resrobot.NewClient()

	if !client.HasAPIKey() {
		fmt.Fprintln(os.Stderr, "Error: RESROBOT_API_KEY not set")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Get a free API key at:")
		fmt.Fprintln(os.Stderr, "  https://www.trafiklab.se/api/trafiklab-apis/resrobot-v21/")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Then set it:")
		fmt.Fprintln(os.Stderr, "  export RESROBOT_API_KEY=\"your-key-here\"")
		os.Exit(1)
	}

	opts := resrobot.TripOptions{
		Time:       searchTime,
		ArriveBy:   arriveBy,
		NumResults: numResults,
	}

	fmt.Fprintf(os.Stderr, "Söker resor från %s till %s (hela Sverige)...\n", origin, dest)

	trips, err := client.PlanTripByName(origin, dest, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(trips) == 0 {
		fmt.Fprintln(os.Stderr, "Inga resor hittades.")
		os.Exit(1)
	}

	output := resrobot.FormatTrips(origin, dest, trips)
	fmt.Print(output)
}

// getDefaultLocation returns the default origin location
func getDefaultLocation() string {
	// Check environment variable
	if loc := os.Getenv("TRANSPORT_DEFAULT_LOCATION"); loc != "" {
		return loc
	}

	// TODO: Read from config file
	// For now, prompt user
	fmt.Fprintln(os.Stderr, "Error: no origin specified and no default location set")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Set a default location with:")
	fmt.Fprintln(os.Stderr, "  export TRANSPORT_DEFAULT_LOCATION=\"Slussen\"")
	os.Exit(1)
	return ""
}

// handleError provides helpful error messages
func handleError(err error, origin, dest string, client *api.Client) {
	errMsg := err.Error()

	if strings.Contains(errMsg, "no stops found") {
		location := origin
		if strings.Contains(errMsg, "destination") {
			location = dest
		} else if strings.Contains(errMsg, "origin") {
			location = origin
		}

		fmt.Fprintf(os.Stderr, "Error: Could not find '%s'\n\n", location)

		// Try to suggest alternatives
		stops, searchErr := client.SearchStops(location)
		if searchErr == nil && len(stops) > 0 {
			fmt.Fprintln(os.Stderr, "Did you mean:")
			for i, stop := range stops {
				if i >= 5 {
					break
				}
				fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, stop.Name)
			}
		}
		return
	}

	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}
