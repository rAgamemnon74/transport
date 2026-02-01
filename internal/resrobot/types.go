package resrobot

import "time"

// LocationResponse represents the response from location.name endpoint
type LocationResponse struct {
	StopLocationOrCoordLocation []StopLocation `json:"stopLocationOrCoordLocation"`
}

// StopLocation represents a stop/station
type StopLocation struct {
	StopLocation *StopLocationData `json:"StopLocation,omitempty"`
}

// StopLocationData contains the actual stop data
type StopLocationData struct {
	ID          string   `json:"id"`           // e.g., "740000001" (Central Station)
	ExtID       string   `json:"extId"`        // External ID
	Name        string   `json:"name"`         // e.g., "Stockholm Centralstation"
	Lon         float64  `json:"lon"`
	Lat         float64  `json:"lat"`
	Weight      int      `json:"weight"`
	Products    int      `json:"products"`     // Bitmask of available transport types
	TimezoneOffset int   `json:"timezoneOffset"`
}

// TripResponse represents the response from trip endpoint
type TripResponse struct {
	Trip []Trip `json:"Trip"`
}

// Trip represents a journey from origin to destination
type Trip struct {
	Duration       string    `json:"duration"`        // e.g., "PT1H30M"
	TripStatus     string    `json:"tripStatus"`
	Origin         TripStop  `json:"Origin"`
	Destination    TripStop  `json:"Destination"`
	LegList        LegList   `json:"LegList"`
	ServiceDays    []ServiceDay `json:"ServiceDays,omitempty"`
}

// LegList contains the legs of a trip
type LegList struct {
	Leg []Leg `json:"Leg"`
}

// Leg represents a segment of a journey
type Leg struct {
	Origin       TripStop    `json:"Origin"`
	Destination  TripStop    `json:"Destination"`
	Notes        *Notes      `json:"Notes,omitempty"`
	Product      *Product    `json:"Product,omitempty"`
	Name         string      `json:"name"`           // e.g., "Pendeltåg 41"
	Type         string      `json:"type"`           // e.g., "JNY" (journey), "WALK"
	Direction    string      `json:"direction"`
	Duration     string      `json:"duration"`       // e.g., "PT45M"
	Dist         int         `json:"dist,omitempty"` // Distance in meters
	Category     string      `json:"category"`       // e.g., "PEN", "MET", "BUS"
	Number       string      `json:"number"`         // Line number
	Operator     string      `json:"operator"`       // Operator name
	OperatorCode string      `json:"operatorCode"`
}

// TripStop represents a stop in a trip
type TripStop struct {
	Name           string `json:"name"`
	ID             string `json:"id"`
	ExtID          string `json:"extId"`
	Lon            float64 `json:"lon"`
	Lat            float64 `json:"lat"`
	Time           string `json:"time"`           // "HH:MM:SS"
	Date           string `json:"date"`           // "YYYY-MM-DD"
	Track          string `json:"track,omitempty"`
	RtTime         string `json:"rtTime,omitempty"`  // Real-time
	RtDate         string `json:"rtDate,omitempty"`
	RtTrack        string `json:"rtTrack,omitempty"`
	Prognosistype  string `json:"prognosisType,omitempty"`
}

// Product represents transport product info
type Product struct {
	Name         string `json:"name"`
	Num          string `json:"num"`
	Line         string `json:"line"`
	CatOut       string `json:"catOut"`       // Category output (e.g., "PEN")
	CatOutS      string `json:"catOutS"`      // Short category
	CatOutL      string `json:"catOutL"`      // Long category (e.g., "Pendeltåg")
	CatIn        string `json:"catIn"`
	CatCode      string `json:"catCode"`
	Operator     string `json:"operator"`
	OperatorCode string `json:"operatorCode"`
}

// Notes contains additional notes for a leg
type Notes struct {
	Note []Note `json:"Note"`
}

// Note represents a single note
type Note struct {
	Value    string `json:"value"`
	Key      string `json:"key"`
	Type     string `json:"type"`
	Priority int    `json:"priority"`
}

// ServiceDay represents service day information
type ServiceDay struct {
	PlanningPeriodBegin string `json:"planningPeriodBegin"`
	PlanningPeriodEnd   string `json:"planningPeriodEnd"`
	SDaysR              string `json:"sDaysR"`
	SDaysB              string `json:"sDaysB"`
}

// DepartureBoard represents the response from departureBoard endpoint
type DepartureBoardResponse struct {
	Departure []DepartureBoardItem `json:"Departure"`
}

// DepartureBoardItem represents a single departure
type DepartureBoardItem struct {
	Product       Product `json:"Product"`
	Stops         *Stops  `json:"Stops,omitempty"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	Stop          string  `json:"stop"`
	StopID        string  `json:"stopid"`
	StopExtID     string  `json:"stopExtId"`
	Time          string  `json:"time"`
	Date          string  `json:"date"`
	RtTime        string  `json:"rtTime,omitempty"`
	RtDate        string  `json:"rtDate,omitempty"`
	Direction     string  `json:"direction"`
	DirectionFlag string  `json:"directionFlag"`
	Track         string  `json:"track,omitempty"`
	RtTrack       string  `json:"rtTrack,omitempty"`
}

// Stops contains stop sequence
type Stops struct {
	Stop []StopSequence `json:"Stop"`
}

// StopSequence represents a stop in the sequence
type StopSequence struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	ExtID     string `json:"extId"`
	Lon       float64 `json:"lon"`
	Lat       float64 `json:"lat"`
	ArrTime   string `json:"arrTime,omitempty"`
	ArrDate   string `json:"arrDate,omitempty"`
	DepTime   string `json:"depTime,omitempty"`
	DepDate   string `json:"depDate,omitempty"`
}

// ParsedTrip is our internal representation of a trip
type ParsedTrip struct {
	Duration     time.Duration
	Interchanges int
	Legs         []ParsedLeg
}

// ParsedLeg is our internal representation of a leg
type ParsedLeg struct {
	Origin        string
	OriginTrack   string
	Destination   string
	DestTrack     string
	DepartureTime time.Time
	ArrivalTime   time.Time
	RtDeparture   *time.Time // Real-time departure
	RtArrival     *time.Time // Real-time arrival
	Line          string
	Category      string     // PEN, MET, BUS, REG, SJ, etc.
	Operator      string
	Direction     string
	IsWalk        bool
	Distance      int        // meters, for walking
}

// Transport category codes
const (
	CatPendeltag    = "PEN"  // Pendeltåg (commuter rail)
	CatMetro        = "MET"  // Tunnelbana
	CatBus          = "BUS"  // Buss
	CatRegionaltag  = "REG"  // Regionaltåg
	CatSJ           = "SJ"   // SJ (long-distance train)
	CatSnabbtag     = "SNT"  // Snabbtåg
	CatNorrtag      = "NRT"  // Norrtåg
	CatSparvagn     = "SPN"  // Spårvagn (tram)
	CatFerry        = "BAT"  // Båt/Färja
	CatFlygbuss     = "FLY"  // Flygbuss
	CatNattbus      = "NAT"  // Nattbuss
)

// GetCategoryEmoji returns an emoji for the transport category
func GetCategoryEmoji(cat string) string {
	switch cat {
	case CatPendeltag, CatRegionaltag, CatSJ, CatSnabbtag, CatNorrtag:
		return "🚆"
	case CatMetro:
		return "🚇"
	case CatBus, CatNattbus:
		return "🚌"
	case CatSparvagn:
		return "🚊"
	case CatFerry:
		return "⛴️"
	case CatFlygbuss:
		return "🚐"
	default:
		return "🚍"
	}
}

// GetCategoryName returns a readable name for the category
func GetCategoryName(cat string) string {
	switch cat {
	case CatPendeltag:
		return "Pendeltåg"
	case CatMetro:
		return "Tunnelbana"
	case CatBus:
		return "Buss"
	case CatRegionaltag:
		return "Regionaltåg"
	case CatSJ:
		return "SJ"
	case CatSnabbtag:
		return "Snabbtåg"
	case CatNorrtag:
		return "Norrtåg"
	case CatSparvagn:
		return "Spårvagn"
	case CatFerry:
		return "Färja"
	case CatFlygbuss:
		return "Flygbuss"
	case CatNattbus:
		return "Nattbuss"
	default:
		return cat
	}
}
