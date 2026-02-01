package api

import (
	"encoding/json"
	"strings"
)

// StopFinderResponse represents the response from /stop-finder endpoint
type StopFinderResponse struct {
	Locations      []Location      `json:"locations"`
	SystemMessages []SystemMessage `json:"systemMessages,omitempty"`
}

// Location represents a stop or place
type Location struct {
	ID               string    `json:"id"`
	IsGlobalID       bool      `json:"isGlobalId"`
	Name             string    `json:"name"`
	DisassembledName string    `json:"disassembledName"`
	Type             string    `json:"type"`
	Coord            []float64 `json:"coord"`
	Parent           *Parent   `json:"parent,omitempty"`
	ProductClasses   []int     `json:"productClasses,omitempty"`
	MatchQuality     int       `json:"matchQuality,omitempty"`
}

// Parent represents a parent location
type Parent struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	DisassembledName string  `json:"disassembledName,omitempty"`
	Type             string  `json:"type"`
	Parent           *Parent `json:"parent,omitempty"`
}

// TripsResponse represents the response from /trips endpoint
type TripsResponse struct {
	Journeys       []Journey       `json:"journeys"`
	SystemMessages []SystemMessage `json:"systemMessages,omitempty"`
}

// SystemMessage represents API status or error messages
type SystemMessage struct {
	Type    string `json:"type"`
	Module  string `json:"module"`
	Code    int    `json:"code"`
	Text    string `json:"text"`
	SubType string `json:"subType,omitempty"`
}

// Journey represents a complete trip from origin to destination
type Journey struct {
	TripDuration   int   `json:"tripDuration"`   // in seconds
	TripRTDuration int   `json:"tripRtDuration"` // real-time duration
	Rating         int   `json:"rating,omitempty"`
	IsAdditional   bool  `json:"isAdditional,omitempty"`
	Interchanges   int   `json:"interchanges"`
	Legs           []Leg `json:"legs"`
}

// Leg represents a segment of a journey
type Leg struct {
	Duration       int             `json:"duration"` // in seconds
	Origin         StopPoint       `json:"origin"`
	Destination    StopPoint       `json:"destination"`
	Transportation *Transportation `json:"transportation,omitempty"`
	Infos          json.RawMessage `json:"infos,omitempty"`
	StopSequence   []StopPoint     `json:"stopSequence,omitempty"`
	Coords         json.RawMessage `json:"coords,omitempty"`
	FootPathInfo   json.RawMessage `json:"footPathInfo,omitempty"`
}

// StopPoint represents a stop with timing information
type StopPoint struct {
	ID               string          `json:"id"`
	IsGlobalID       bool            `json:"isGlobalId"`
	Name             string          `json:"name"`
	DisassembledName string          `json:"disassembledName"`
	Type             string          `json:"type"`
	Coord            []float64       `json:"coord"`
	Parent           *Parent         `json:"parent,omitempty"`
	ProductClasses   []int           `json:"productClasses,omitempty"`
	Properties       json.RawMessage `json:"properties,omitempty"`

	// Departure times
	DepartureTimePlanned   string `json:"departureTimePlanned,omitempty"`
	DepartureTimeEstimated string `json:"departureTimeEstimated,omitempty"`

	// Arrival times
	ArrivalTimePlanned   string `json:"arrivalTimePlanned,omitempty"`
	ArrivalTimeEstimated string `json:"arrivalTimeEstimated,omitempty"`
}

// GetPlatform extracts platform from properties
func (s *StopPoint) GetPlatform() string {
	if s.Properties == nil {
		return ""
	}
	var props struct {
		Platform     string `json:"platform"`
		PlatformName string `json:"platformName"`
	}
	if err := json.Unmarshal(s.Properties, &props); err != nil {
		return ""
	}
	if props.PlatformName != "" {
		return props.PlatformName
	}
	return props.Platform
}

// GetStopName returns the clean stop name
func (s *StopPoint) GetStopName() string {
	// For platforms, get parent stop name
	if s.Type == "platform" && s.Parent != nil {
		if s.Parent.DisassembledName != "" {
			return s.Parent.DisassembledName
		}
		return s.Parent.Name
	}
	if s.DisassembledName != "" {
		return s.DisassembledName
	}
	return s.Name
}

// Transportation represents the vehicle/mode of transport
type Transportation struct {
	ID          string              `json:"id,omitempty"`
	Name        string              `json:"name"`
	Number      string              `json:"number,omitempty"`
	Description string              `json:"description,omitempty"`
	Product     *Product            `json:"product,omitempty"`
	Operator    *Operator           `json:"operator,omitempty"`
	Destination *TransportDest      `json:"destination,omitempty"`
	Properties  json.RawMessage     `json:"properties,omitempty"`
}

// Product represents product/line information
type Product struct {
	ID        int    `json:"id"`
	Class     int    `json:"class"`
	Name      string `json:"name"`
	IconID    int    `json:"iconId"`
	ShortName string `json:"shortName,omitempty"`
}

// Operator represents the transport operator
type Operator struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TransportDest represents where the vehicle is heading
type TransportDest struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	DisassembledName string `json:"disassembledName,omitempty"`
	Type             string `json:"type"`
}

// GetDirection returns the destination name for display
func (t *Transportation) GetDirection() string {
	if t.Destination == nil {
		return ""
	}
	if t.Destination.DisassembledName != "" {
		return t.Destination.DisassembledName
	}
	return t.Destination.Name
}

// GetLineName returns the line name/number for display
func (t *Transportation) GetLineName() string {
	// Try to extract from the number field (e.g., "tunnelbanans gröna linje 18" -> "18")
	if t.Number != "" {
		lineNum := extractTrailingNumber(t.Number)
		if lineNum != "" {
			return lineNum
		}
	}

	// Try from ID (e.g., "tfs:02018: :R:y01" -> "18")
	if t.ID != "" {
		lineNum := extractLineFromID(t.ID)
		if lineNum != "" {
			return lineNum
		}
	}

	if t.Product != nil && t.Product.ShortName != "" {
		return t.Product.ShortName
	}
	if t.Number != "" {
		return t.Number
	}
	return t.Name
}

// extractTrailingNumber extracts the trailing number from strings like
// "tunnelbanans gröna linje 18" -> "18"
// "Buss 4" -> "4"
func extractTrailingNumber(s string) string {
	// Find the last sequence of digits at the end of the string
	end := len(s)
	for end > 0 && s[end-1] == ' ' {
		end--
	}

	start := end
	for start > 0 && s[start-1] >= '0' && s[start-1] <= '9' {
		start--
	}

	if start == end {
		return ""
	}

	result := s[start:end]
	// Remove leading zeros
	for len(result) > 1 && result[0] == '0' {
		result = result[1:]
	}

	return result
}

// extractLineFromID extracts line number from ID format "tfs:02018: :R:y01" -> "18"
func extractLineFromID(id string) string {
	// Look for pattern "tfs:0XXXX:" where XXXX is the line
	if len(id) < 10 {
		return ""
	}
	if id[:4] != "tfs:" {
		return ""
	}

	// Find the end of the number part
	start := 4
	end := start
	for end < len(id) && id[end] != ':' {
		end++
	}

	if end == start {
		return ""
	}

	num := id[start:end]
	// Remove leading zeros
	for len(num) > 1 && num[0] == '0' {
		num = num[1:]
	}

	return num
}

// ProductClass constants for transport types
const (
	ProductClassTrain    = 1  // Pendeltåg
	ProductClassMetro    = 2  // Tunnelbana
	ProductClassTram     = 4  // Spårvagn/Tvärbanan
	ProductClassBus      = 5  // Buss
	ProductClassFerry    = 9  // Båt
	ProductClassFootpath = 99 // Walking/footpath
)

// TransportMode constants for SL Transport API
const (
	TransportModeBus   = "BUS"
	TransportModeMetro = "METRO"
	TransportModeTrain = "TRAIN"
	TransportModeTram  = "TRAM"
	TransportModeShip  = "SHIP"
)

// Site represents a stop/station from the sites API
type Site struct {
	ID           int      `json:"id"`
	GID          int64    `json:"gid"`
	Name         string   `json:"name"`
	Alias        []string `json:"alias,omitempty"`
	Abbreviation string   `json:"abbreviation,omitempty"`
	Note         string   `json:"note,omitempty"`
	Lat          float64  `json:"lat"`
	Lon          float64  `json:"lon"`
}

// MatchesQuery checks if the site matches a search query
func (s *Site) MatchesQuery(query string) bool {
	queryLower := strings.ToLower(query)

	// Check name
	if strings.Contains(strings.ToLower(s.Name), queryLower) {
		return true
	}

	// Check aliases
	for _, alias := range s.Alias {
		if strings.Contains(strings.ToLower(alias), queryLower) {
			return true
		}
	}

	// Check abbreviation
	if s.Abbreviation != "" && strings.EqualFold(s.Abbreviation, query) {
		return true
	}

	return false
}

// DeparturesResponse represents the response from the departures endpoint
type DeparturesResponse struct {
	Departures     []Departure     `json:"departures"`
	StopDeviations []StopDeviation `json:"stop_deviations,omitempty"`
}

// Departure represents a single departure
type Departure struct {
	Destination   string        `json:"destination"`
	DirectionCode int           `json:"direction_code"`
	Direction     string        `json:"direction"`
	State         string        `json:"state"`
	Display       string        `json:"display"`
	Scheduled     string        `json:"scheduled"`
	Expected      string        `json:"expected"`
	Journey       JourneyInfo   `json:"journey"`
	StopArea      StopAreaInfo  `json:"stop_area"`
	StopPoint     StopPointInfo `json:"stop_point"`
	Line          LineInfo      `json:"line"`
	Deviations    []Deviation   `json:"deviations,omitempty"`
}

// JourneyInfo contains journey state information
type JourneyInfo struct {
	ID              int64  `json:"id"`
	State           string `json:"state"`
	PredictionState string `json:"prediction_state"`
}

// StopAreaInfo contains stop area details
type StopAreaInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// StopPointInfo contains stop point details
type StopPointInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Designation string `json:"designation"`
}

// LineInfo contains line details
type LineInfo struct {
	ID                   int    `json:"id"`
	Designation          string `json:"designation"`
	TransportAuthorityID int    `json:"transport_authority_id"`
	TransportMode        string `json:"transport_mode"`
	GroupOfLines         string `json:"group_of_lines,omitempty"`
}

// Deviation represents a service deviation
type Deviation struct {
	Importance  int    `json:"importance"`
	Consequence string `json:"consequence"`
	Message     string `json:"message"`
}

// StopDeviation represents a stop-level deviation
type StopDeviation struct {
	Importance  int    `json:"importance"`
	Consequence string `json:"consequence"`
	Message     string `json:"message"`
}

// IsWalking returns true if this leg is a walking segment
func (t *Transportation) IsWalking() bool {
	if t == nil {
		return true
	}
	if t.Product != nil && t.Product.Class == ProductClassFootpath {
		return true
	}
	if t.Product != nil && t.Product.Name == "footpath" {
		return true
	}
	return false
}
