package flight

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// FlightSearch contains parameters for a flight search
type FlightSearch struct {
	OriginCode      string
	OriginCity      string
	DestinationCode string
	DestinationCity string
	DepartureDate   time.Time
	ReturnDate      time.Time // Zero value = one-way
	Passengers      int
	ShowPrivate     bool      // Show private jet/helicopter options
}

// CityToAirport maps city names to their main IATA airport codes
var CityToAirport = map[string]string{
	// Swedish cities
	"stockholm":   "ARN",
	"göteborg":    "GOT",
	"gothenburg":  "GOT",
	"malmö":       "MMX",
	"malmo":       "MMX",
	"luleå":       "LLA",
	"lulea":       "LLA",
	"umeå":        "UME",
	"umea":        "UME",
	"kiruna":      "KRN",
	"sundsvall":   "SDL",
	"östersund":   "OSD",
	"ostersund":   "OSD",
	"växjö":       "VXO",
	"vaxjo":       "VXO",
	"kalmar":      "KLR",
	"visby":       "VBY",
	"karlstad":    "KSD",
	"linköping":   "LPI",
	"linkoping":   "LPI",
	"norrköping":  "NRK",
	"norrkoping":  "NRK",
	"örebro":      "ORB",
	"orebro":      "ORB",
	"jönköping":   "JKG",
	"jonkoping":   "JKG",
	"halmstad":    "HAD",
	"ängelholm":   "AGH",
	"angelholm":   "AGH",
	"ronneby":     "RNB",
	"trollhättan": "THN",
	"trollhattan": "THN",
	"arvidsjaur":  "AJR",
	"gällivare":   "GEV",
	"gallivare":   "GEV",
	"hemavan":     "HMV",
	"lycksele":    "LYC",
	"mora":        "MXX",
	"pajala":      "PJA",
	"skellefteå":  "SFT",
	"skelleftea":  "SFT",
	"storuman":    "SQO",
	"sveg":        "EVG",
	"vilhelmina":  "VHM",
	"åre":         "ARE",
	"are":         "ARE",
	"bromma":      "BMA",
	"arlanda":     "ARN",
	"landvetter":  "GOT",
	"sturup":      "MMX",
	"skavsta":     "NYO",
	"västerås":    "VST",
	"vasteras":    "VST",

	// Nordic capitals
	"oslo":        "OSL",
	"copenhagen":  "CPH",
	"köpenhamn":   "CPH",
	"kopenhamn":   "CPH",
	"helsinki":    "HEL",
	"helsingfors": "HEL",
	"reykjavik":   "KEF",

	// Baltic
	"vilnius":   "VNO",
	"riga":      "RIX",
	"tallinn":   "TLL",
	"kaunas":    "KUN",
	"palanga":   "PLQ",

	// Major European cities
	"london":      "LHR",
	"paris":       "CDG",
	"amsterdam":   "AMS",
	"frankfurt":   "FRA",
	"munich":      "MUC",
	"münchen":     "MUC",
	"munchen":     "MUC",
	"berlin":      "BER",
	"hamburg":     "HAM",
	"düsseldorf":  "DUS",
	"dusseldorf":  "DUS",
	"zürich":      "ZRH",
	"zurich":      "ZRH",
	"geneva":      "GVA",
	"geneve":      "GVA",
	"wien":        "VIE",
	"vienna":      "VIE",
	"brussels":    "BRU",
	"bryssel":     "BRU",
	"rome":        "FCO",
	"roma":        "FCO",
	"rom":         "FCO",
	"milan":       "MXP",
	"milano":      "MXP",
	"barcelona":   "BCN",
	"madrid":      "MAD",
	"lisbon":      "LIS",
	"lissabon":    "LIS",
	"dublin":      "DUB",
	"edinburgh":   "EDI",
	"manchester":  "MAN",
	"prague":      "PRG",
	"prag":        "PRG",
	"warsaw":      "WAW",
	"warszawa":    "WAW",
	"krakow":      "KRK",
	"krakau":      "KRK",
	"budapest":    "BUD",
	"athens":      "ATH",
	"aten":        "ATH",
	"istanbul":    "IST",
	"nice":        "NCE",
	"lyon":        "LYS",
	"marseille":   "MRS",

	// South America
	"buenos aires": "EZE",
	"sao paulo":    "GRU",
	"são paulo":    "GRU",
	"rio de janeiro": "GIG",
	"rio":          "GIG",
	"lima":         "LIM",
	"bogota":       "BOG",
	"bogotá":       "BOG",
	"santiago":     "SCL",
	"caracas":      "CCS",
	"medellin":     "MDE",
	"medellín":     "MDE",
	"quito":        "UIO",
	"montevideo":   "MVD",

	// Central America & Caribbean
	"managua":      "MGA",
	"panama city":  "PTY",
	"panama":       "PTY",
	"san jose":     "SJO",
	"san josé":     "SJO",
	"guatemala city": "GUA",
	"guatemala":    "GUA",
	"havana":       "HAV",
	"havanna":      "HAV",
	"kingston":     "KIN",
	"santo domingo": "SDQ",
	"san juan":     "SJU",
	"puerto rico":  "SJU",
	"punta cana":   "PUJ",
	"montego bay":  "MBJ",
	"aruba":        "AUA",
	"curacao":      "CUR",
	"curaçao":      "CUR",
	"nassau":       "NAS",
	"bahamas":      "NAS",
	"barbados":     "BGI",
	"belize city":  "BZE",
	"belize":       "BZE",
	"tegucigalpa":  "TGU",
	"honduras":     "TGU",
	"san salvador": "SAL",
	"el salvador":  "SAL",

	// Africa
	"casablanca":   "CMN",
	"marrakech":    "RAK",
	"tunis":        "TUN",
	"algiers":      "ALG",
	"cairo":        "CAI",
	"kairo":        "CAI",
	"johannesburg": "JNB",
	"cape town":    "CPT",
	"kapstaden":    "CPT",
	"nairobi":      "NBO",
	"addis ababa":  "ADD",
	"lagos":        "LOS",
	"accra":        "ACC",
	"dakar":        "DSS",
	"dar es salaam": "DAR",
	"zanzibar":     "ZNZ",
	"mauritius":    "MRU",
	"seychelles":   "SEZ",
	"kigali":       "KGL",
	"windhoek":     "WDH",
	"luanda":       "LAD",
	"maputo":       "MPM",
	"durban":       "DUR",
	"mombasa":      "MBA",
	"kilimanjaro":  "JRO",
	"victoria falls": "VFA",
	"livingstone":  "LVI",

	// Middle East
	"dubai":        "DXB",
	"abu dhabi":    "AUH",
	"doha":         "DOH",
	"qatar":        "DOH",
	"riyadh":       "RUH",
	"jeddah":       "JED",
	"kuwait":       "KWI",
	"muscat":       "MCT",
	"oman":         "MCT",
	"bahrain":      "BAH",
	"amman":        "AMM",
	"jordan":       "AMM",
	"beirut":       "BEY",
	"tel aviv":     "TLV",
	"telaviv":      "TLV",
	"jerusalem":    "TLV",

	// Asia
	"bangkok":      "BKK",
	"singapore":    "SIN",
	"tokyo":        "NRT",
	"osaka":        "KIX",
	"kyoto":        "KIX",
	"hong kong":    "HKG",
	"hongkong":     "HKG",
	"seoul":        "ICN",
	"beijing":      "PEK",
	"peking":       "PEK",
	"shanghai":     "PVG",
	"taipei":       "TPE",
	"taiwan":       "TPE",
	"manila":       "MNL",
	"kuala lumpur": "KUL",
	"jakarta":      "CGK",
	"ho chi minh":  "SGN",
	"saigon":       "SGN",
	"hanoi":        "HAN",
	"mumbai":       "BOM",
	"bombay":       "BOM",
	"delhi":        "DEL",
	"new delhi":    "DEL",
	"bangalore":    "BLR",
	"chennai":      "MAA",
	"kolkata":      "CCU",
	"goa":          "GOI",
	"kathmandu":    "KTM",
	"nepal":        "KTM",
	"colombo":      "CMB",
	"sri lanka":    "CMB",
	"maldives":     "MLE",
	"male":         "MLE",
	"phuket":       "HKT",
	"bali":         "DPS",
	"denpasar":     "DPS",
	"krabi":        "KBV",
	"chiang mai":   "CNX",
	"siem reap":    "REP",
	"phnom penh":   "PNH",
	"yangon":       "RGN",
	"myanmar":      "RGN",
	"dhaka":        "DAC",
	"bangladesh":   "DAC",
	"karachi":      "KHI",
	"islamabad":    "ISB",
	"lahore":       "LHE",

	// Australia & Pacific
	"sydney":       "SYD",
	"melbourne":    "MEL",
	"brisbane":     "BNE",
	"perth":        "PER",
	"adelaide":     "ADL",
	"auckland":     "AKL",
	"wellington":   "WLG",
	"christchurch": "CHC",
	"queenstown":   "ZQN",
	"fiji":         "NAN",
	"tahiti":       "PPT",
	"hawaii":       "HNL",
	"honolulu":     "HNL",

	// North America
	"new york":     "JFK",
	"newyork":      "JFK",
	"nyc":          "JFK",
	"los angeles":  "LAX",
	"la":           "LAX",
	"chicago":      "ORD",
	"miami":        "MIA",
	"san francisco": "SFO",
	"boston":       "BOS",
	"washington":   "IAD",
	"seattle":      "SEA",
	"las vegas":    "LAS",
	"denver":       "DEN",
	"atlanta":      "ATL",
	"dallas":       "DFW",
	"houston":      "IAH",
	"phoenix":      "PHX",
	"philadelphia": "PHL",
	"san diego":    "SAN",
	"orlando":      "MCO",
	"new orleans":  "MSY",
	"portland":     "PDX",
	"detroit":      "DTW",
	"minneapolis":  "MSP",
	"salt lake city": "SLC",
	"austin":       "AUS",
	"nashville":    "BNA",
	"toronto":      "YYZ",
	"vancouver":    "YVR",
	"montreal":     "YUL",
	"calgary":      "YYC",
	"ottawa":       "YOW",
	"mexico city":  "MEX",
	"cancun":       "CUN",
	"guadalajara":  "GDL",
	"monterrey":    "MTY",
	"puerto vallarta": "PVR",
	"los cabos":    "SJD",
	"cabo":         "SJD",

	// Popular vacation destinations
	"alicante":    "ALC",
	"malaga":      "AGP",
	"palma":       "PMI",
	"mallorca":    "PMI",
	"ibiza":       "IBZ",
	"tenerife":    "TFS",
	"gran canaria": "LPA",
	"las palmas":  "LPA",
	"fuerteventura": "FUE",
	"lanzarote":   "ACE",
	"rhodes":      "RHO",
	"rhodos":      "RHO",
	"kreta":       "HER",
	"crete":       "HER",
	"heraklion":   "HER",
	"corfu":       "CFU",
	"korfu":       "CFU",
	"santorini":   "JTR",
	"antalya":     "AYT",
	"bodrum":      "BJV",
	"split":       "SPU",
	"dubrovnik":   "DBV",
	"pula":        "PUY",
	"zadar":       "ZAD",
	"malta":       "MLA",
	"larnaca":     "LCA",
	"paphos":      "PFO",
	"faro":        "FAO",
	"funchal":     "FNC",
	"madeira":     "FNC",
	"island":      "KEF",
	"iceland":     "KEF",
	"hurghada":    "HRG",
	"sharm el sheikh": "SSH",
}

// LookupAirportCode looks up IATA code for a city name
// Returns the code and whether it was found
func LookupAirportCode(city string) (string, bool) {
	// First check if it's already an IATA code (3 uppercase letters)
	upper := strings.ToUpper(strings.TrimSpace(city))
	if len(upper) == 3 && isAllLetters(upper) {
		return upper, true
	}

	// Look up by city name
	code, ok := CityToAirport[strings.ToLower(strings.TrimSpace(city))]
	return code, ok
}

func isAllLetters(s string) bool {
	for _, c := range s {
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}

// formatAirportDisplay formats airport for display
// If city name equals the code, just show code
// Otherwise show "City (CODE)"
func formatAirportDisplay(city, code string) string {
	if city == "" || strings.ToUpper(city) == code {
		return code
	}
	return fmt.Sprintf("%s (%s)", city, code)
}

// GenerateGoogleFlightsURL creates a Google Flights search URL
func GenerateGoogleFlightsURL(search FlightSearch) string {
	// Google Flights URL format
	// https://www.google.com/travel/flights/search?tfs=CBwQ...
	// Simpler format using query params:
	baseURL := "https://www.google.com/travel/flights"

	params := url.Values{}

	// Build the search query
	q := fmt.Sprintf("flights from %s to %s", search.OriginCode, search.DestinationCode)
	if !search.DepartureDate.IsZero() {
		q += " on " + search.DepartureDate.Format("Jan 2")
	}
	params.Set("q", q)

	return baseURL + "?" + params.Encode()
}

// GenerateSkyscannerURL creates a Skyscanner search URL
func GenerateSkyscannerURL(search FlightSearch) string {
	// Skyscanner URL format:
	// https://www.skyscanner.se/transport/flights/arn/vno/250215/
	// Or for round-trip: /arn/vno/250215/250220/

	origin := strings.ToLower(search.OriginCode)
	dest := strings.ToLower(search.DestinationCode)

	path := fmt.Sprintf("/transport/flights/%s/%s", origin, dest)

	if !search.DepartureDate.IsZero() {
		path += "/" + search.DepartureDate.Format("060102")
		if !search.ReturnDate.IsZero() {
			path += "/" + search.ReturnDate.Format("060102")
		}
	}

	return "https://www.skyscanner.se" + path + "/"
}

// GenerateMomondoURL creates a Momondo search URL
func GenerateMomondoURL(search FlightSearch) string {
	// Momondo URL format:
	// https://www.momondo.se/flight-search/ARN-VNO/2025-02-15

	path := fmt.Sprintf("/flight-search/%s-%s", search.OriginCode, search.DestinationCode)

	if !search.DepartureDate.IsZero() {
		path += "/" + search.DepartureDate.Format("2006-01-02")
		if !search.ReturnDate.IsZero() {
			path += "/" + search.ReturnDate.Format("2006-01-02")
		}
	}

	return "https://www.momondo.se" + path
}

// GenerateKayakURL creates a Kayak search URL
func GenerateKayakURL(search FlightSearch) string {
	// Kayak URL format:
	// https://www.kayak.se/flights/ARN-VNO/2025-02-15

	path := fmt.Sprintf("/flights/%s-%s", search.OriginCode, search.DestinationCode)

	if !search.DepartureDate.IsZero() {
		path += "/" + search.DepartureDate.Format("2006-01-02")
		if !search.ReturnDate.IsZero() {
			path += "/" + search.ReturnDate.Format("2006-01-02")
		}
	}

	return "https://www.kayak.se" + path
}

// GenerateTUIURL creates a TUI charter search URL
func GenerateTUIURL(search FlightSearch) string {
	// TUI URL format:
	// https://www.tui.se/resa/sok/?departureAirportCodes=ARN&destinationCodes=LPA&duration=7-7
	params := url.Values{}
	params.Set("departureAirportCodes", search.OriginCode)
	params.Set("destinationCodes", search.DestinationCode)
	params.Set("flexibleDays", "3")

	if !search.DepartureDate.IsZero() {
		params.Set("departureDate", search.DepartureDate.Format("2006-01-02"))
	}

	return "https://www.tui.se/resa/sok/?" + params.Encode()
}

// GenerateApolloURL creates an Apollo charter search URL
func GenerateApolloURL(search FlightSearch) string {
	// Apollo URL format:
	// https://www.apollo.se/resor/searchresult?DestinationAirportCodes=LPA&DepartureAirportCode=ARN
	params := url.Values{}
	params.Set("DepartureAirportCode", search.OriginCode)
	params.Set("DestinationAirportCodes", search.DestinationCode)
	params.Set("CategoryCodes", "FlightHotel,Flight")
	params.Set("Adults", "1")

	if !search.DepartureDate.IsZero() {
		params.Set("DepartureDate", search.DepartureDate.Format("2006-01-02"))
	}

	return "https://www.apollo.se/resor/searchresult?" + params.Encode()
}

// GenerateVingURL creates a Ving charter search URL
func GenerateVingURL(search FlightSearch) string {
	// Ving URL format similar to Apollo (same company group)
	params := url.Values{}
	params.Set("DepartureAirportCode", search.OriginCode)
	params.Set("DestinationAirportCodes", search.DestinationCode)
	params.Set("CategoryCodes", "FlightHotel,Flight")
	params.Set("Adults", "1")

	if !search.DepartureDate.IsZero() {
		params.Set("DepartureDate", search.DepartureDate.Format("2006-01-02"))
	}

	return "https://www.ving.se/resor/searchresult?" + params.Encode()
}

// GenerateTicketURL creates a Ticket.se search URL
func GenerateTicketURL(search FlightSearch) string {
	// Ticket.se format
	params := url.Values{}
	params.Set("from", search.OriginCode)
	params.Set("to", search.DestinationCode)
	params.Set("adults", "1")

	if !search.DepartureDate.IsZero() {
		params.Set("outDate", search.DepartureDate.Format("2006-01-02"))
	}
	if !search.ReturnDate.IsZero() {
		params.Set("inDate", search.ReturnDate.Format("2006-01-02"))
	}

	return "https://www.ticket.se/flight/search?" + params.Encode()
}

// GenerateGrafairURL creates a Grafair private jet inquiry URL (Bromma-based)
func GenerateGrafairURL(search FlightSearch) string {
	// Grafair is based at Bromma - link to their quote request
	params := url.Values{}
	params.Set("from", search.OriginCode)
	params.Set("to", search.DestinationCode)
	if !search.DepartureDate.IsZero() {
		params.Set("date", search.DepartureDate.Format("2006-01-02"))
	}
	// Grafair doesn't have deep linking, so link to contact/quote page
	return "https://www.grafair.se/offertforfragan/"
}

// GeneratePrivateFlyURL creates a PrivateFly search URL
func GeneratePrivateFlyURL(search FlightSearch) string {
	// PrivateFly quote request format
	params := url.Values{}
	params.Set("departure", search.OriginCode)
	params.Set("arrival", search.DestinationCode)
	params.Set("pax", "1")
	if !search.DepartureDate.IsZero() {
		params.Set("date", search.DepartureDate.Format("02/01/2006"))
	}
	return "https://www.privatefly.com/sv/privat-jet-priser/priser-offert.html?" + params.Encode()
}

// GenerateVictorURL creates a Victor private jet search URL
func GenerateVictorURL(search FlightSearch) string {
	// Victor (FlyVictor) format
	return fmt.Sprintf("https://www.flyvictor.com/en-gb/quote/?from=%s&to=%s",
		search.OriginCode, search.DestinationCode)
}

// GenerateLunaJetsURL creates a LunaJets search URL
func GenerateLunaJetsURL(search FlightSearch) string {
	params := url.Values{}
	params.Set("departure_1", search.OriginCode)
	params.Set("arrival_1", search.DestinationCode)
	if !search.DepartureDate.IsZero() {
		params.Set("date_1", search.DepartureDate.Format("2006-01-02"))
	}
	return "https://www.lunajets.com/en/instant-estimate/?" + params.Encode()
}

// GenerateHeliAirURL creates HeliAir Sweden inquiry URL
func GenerateHeliAirURL() string {
	return "https://heliair.se/boka/"
}

// GenerateHelipadyURL creates Helipady helicopter booking URL
func GenerateHelipadyURL(search FlightSearch) string {
	// Helipady is a helicopter booking platform
	return fmt.Sprintf("https://www.helipady.com/search?from=%s&to=%s",
		search.OriginCode, search.DestinationCode)
}

// GenerateNorwegianURL creates a Norwegian Air search URL
func GenerateNorwegianURL(search FlightSearch) string {
	// Norwegian URL format:
	// https://www.norwegian.com/se/booking/fly/low-fare/?D_City=ARN&A_City=VNO&D_Day=15&D_Month=202502

	params := url.Values{}
	params.Set("D_City", search.OriginCode)
	params.Set("A_City", search.DestinationCode)
	params.Set("AdultCount", "1")
	params.Set("TripType", "1") // 1 = one-way, 2 = round-trip

	if !search.DepartureDate.IsZero() {
		params.Set("D_Day", search.DepartureDate.Format("02"))
		params.Set("D_Month", search.DepartureDate.Format("200601"))
	}

	if !search.ReturnDate.IsZero() {
		params.Set("TripType", "2")
		params.Set("R_Day", search.ReturnDate.Format("02"))
		params.Set("R_Month", search.ReturnDate.Format("200601"))
	}

	return "https://www.norwegian.com/se/booking/fly/low-fare/?" + params.Encode()
}

// GenerateSASURL creates an SAS search URL
func GenerateSASURL(search FlightSearch) string {
	// SAS URL format:
	// https://www.flysas.com/se-sv/book/flights?from=ARN&to=VNO&out=2025-02-15

	params := url.Values{}
	params.Set("from", search.OriginCode)
	params.Set("to", search.DestinationCode)
	params.Set("adt", "1")

	if !search.DepartureDate.IsZero() {
		params.Set("out", search.DepartureDate.Format("2006-01-02"))
	}

	if !search.ReturnDate.IsZero() {
		params.Set("in", search.ReturnDate.Format("2006-01-02"))
	}

	return "https://www.flysas.com/se-sv/book/flights?" + params.Encode()
}

// FormatFlightSearch formats flight search results for display
func FormatFlightSearch(search FlightSearch) string {
	var sb strings.Builder

	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Header with route
	originDisplay := formatAirportDisplay(search.OriginCity, search.OriginCode)
	destDisplay := formatAirportDisplay(search.DestinationCity, search.DestinationCode)

	sb.WriteString(fmt.Sprintf(" ✈️  Flyg: %s → %s\n", originDisplay, destDisplay))
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Date info
	if !search.DepartureDate.IsZero() {
		sb.WriteString(fmt.Sprintf("  Datum:       %s\n", search.DepartureDate.Format("Mon 2 Jan 2006")))
		if !search.ReturnDate.IsZero() {
			sb.WriteString(fmt.Sprintf("  Retur:       %s\n", search.ReturnDate.Format("Mon 2 Jan 2006")))
		} else {
			sb.WriteString("  Retur:       Enkel resa\n")
		}
		sb.WriteString("\n")
	}

	// Flight search URLs
	sb.WriteString("  Sök flyg:\n")
	sb.WriteString("  ─────────────────────────────────────────────────────────────────\n")

	sb.WriteString(fmt.Sprintf("  🔍 Google Flights:\n     %s\n\n", GenerateGoogleFlightsURL(search)))
	sb.WriteString(fmt.Sprintf("  🔍 Skyscanner:\n     %s\n\n", GenerateSkyscannerURL(search)))
	sb.WriteString(fmt.Sprintf("  🔍 Momondo:\n     %s\n\n", GenerateMomondoURL(search)))
	sb.WriteString(fmt.Sprintf("  🔍 Kayak:\n     %s\n\n", GenerateKayakURL(search)))

	// Direct airline links
	sb.WriteString("  Flygbolag:\n")
	sb.WriteString("  ─────────────────────────────────────────────────────────────────\n")
	sb.WriteString(fmt.Sprintf("  🇸🇪 SAS:\n     %s\n\n", GenerateSASURL(search)))
	sb.WriteString(fmt.Sprintf("  🇳🇴 Norwegian:\n     %s\n\n", GenerateNorwegianURL(search)))

	// Charter companies
	sb.WriteString("  Charterbolag:\n")
	sb.WriteString("  ─────────────────────────────────────────────────────────────────\n")
	sb.WriteString(fmt.Sprintf("  🌴 TUI:\n     %s\n\n", GenerateTUIURL(search)))
	sb.WriteString(fmt.Sprintf("  🌴 Apollo:\n     %s\n\n", GenerateApolloURL(search)))
	sb.WriteString(fmt.Sprintf("  🌴 Ving:\n     %s\n\n", GenerateVingURL(search)))
	sb.WriteString(fmt.Sprintf("  🎫 Ticket:\n     %s\n\n", GenerateTicketURL(search)))

	// Private aviation section
	if search.ShowPrivate {
		sb.WriteString("  Privatflyg & Helikopter:\n")
		sb.WriteString("  ─────────────────────────────────────────────────────────────────\n")
		sb.WriteString(fmt.Sprintf("  🛩️  Grafair (Bromma):\n     %s\n\n", GenerateGrafairURL(search)))
		sb.WriteString(fmt.Sprintf("  🛩️  PrivateFly:\n     %s\n\n", GeneratePrivateFlyURL(search)))
		sb.WriteString(fmt.Sprintf("  🛩️  Victor:\n     %s\n\n", GenerateVictorURL(search)))
		sb.WriteString(fmt.Sprintf("  🛩️  LunaJets:\n     %s\n\n", GenerateLunaJetsURL(search)))
		sb.WriteString(fmt.Sprintf("  🚁 HeliAir Sweden:\n     %s\n\n", GenerateHeliAirURL()))
		sb.WriteString(fmt.Sprintf("  🚁 Helipady:\n     %s\n\n", GenerateHelipadyURL(search)))
	}

	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	return sb.String()
}
