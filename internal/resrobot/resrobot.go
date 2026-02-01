package resrobot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL        = "https://api.resrobot.se/v2.1"
	defaultTimeout = 15 * time.Second
)

// Client handles communication with the ResRobot API
type Client struct {
	httpClient *http.Client
	apiKey     string
}

// NewClient creates a new ResRobot API client
func NewClient() *Client {
	apiKey := os.Getenv("RESROBOT_API_KEY")
	if apiKey == "" {
		// Try alternate env var name
		apiKey = os.Getenv("TRAFIKLAB_RESROBOT_KEY")
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		apiKey: apiKey,
	}
}

// HasAPIKey returns true if an API key is configured
func (c *Client) HasAPIKey() bool {
	return c.apiKey != ""
}

// SearchStops finds stops matching the given query
func (c *Client) SearchStops(query string) ([]StopLocationData, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("RESROBOT_API_KEY not set. Get a free key at https://www.trafiklab.se/api/trafiklab-apis/resrobot-v21/")
	}

	params := url.Values{}
	params.Set("input", query)
	params.Set("format", "json")
	params.Set("accessId", c.apiKey)

	reqURL := fmt.Sprintf("%s/location.name?%s", baseURL, params.Encode())

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search stops: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var result LocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract stop locations
	var stops []StopLocationData
	for _, item := range result.StopLocationOrCoordLocation {
		if item.StopLocation != nil {
			stops = append(stops, *item.StopLocation)
		}
	}

	return stops, nil
}

// TripOptions contains options for trip planning
type TripOptions struct {
	Time       time.Time // Departure or arrival time
	ArriveBy   bool      // If true, time is arrival time
	NumResults int       // Number of results to return
	Via        string    // Optional via station
}

// DefaultTripOptions returns default options
func DefaultTripOptions() TripOptions {
	return TripOptions{
		Time:       time.Now(),
		ArriveBy:   false,
		NumResults: 5,
	}
}

// PlanTrip finds journeys between origin and destination
func (c *Client) PlanTrip(originID, destID string, opts TripOptions) ([]ParsedTrip, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("RESROBOT_API_KEY not set. Get a free key at https://www.trafiklab.se/api/trafiklab-apis/resrobot-v21/")
	}

	params := url.Values{}
	params.Set("originId", originID)
	params.Set("destId", destID)
	params.Set("format", "json")
	params.Set("accessId", c.apiKey)
	params.Set("numF", fmt.Sprintf("%d", opts.NumResults))
	params.Set("passlist", "0") // Don't include intermediate stops

	if !opts.Time.IsZero() {
		params.Set("date", opts.Time.Format("2006-01-02"))
		params.Set("time", opts.Time.Format("15:04"))
	}

	if opts.ArriveBy {
		params.Set("searchForArrival", "1")
	}

	reqURL := fmt.Sprintf("%s/trip?%s", baseURL, params.Encode())

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to plan trip: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var result TripResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Parse trips into our format
	var trips []ParsedTrip
	for _, trip := range result.Trip {
		parsed := parseTrip(trip)
		trips = append(trips, parsed)
	}

	return trips, nil
}

// PlanTripByName resolves stop names and plans the trip
func (c *Client) PlanTripByName(origin, dest string, opts TripOptions) ([]ParsedTrip, error) {
	// Resolve origin
	originStops, err := c.SearchStops(origin)
	if err != nil {
		return nil, fmt.Errorf("failed to find origin '%s': %w", origin, err)
	}
	if len(originStops) == 0 {
		return nil, fmt.Errorf("no stops found for origin '%s'", origin)
	}

	// Resolve destination
	destStops, err := c.SearchStops(dest)
	if err != nil {
		return nil, fmt.Errorf("failed to find destination '%s': %w", dest, err)
	}
	if len(destStops) == 0 {
		return nil, fmt.Errorf("no stops found for destination '%s'", dest)
	}

	return c.PlanTrip(originStops[0].ID, destStops[0].ID, opts)
}

// parseTrip converts API trip to our parsed format
func parseTrip(trip Trip) ParsedTrip {
	parsed := ParsedTrip{
		Duration:     parseDuration(trip.Duration),
		Interchanges: len(trip.LegList.Leg) - 1,
	}

	for _, leg := range trip.LegList.Leg {
		parsedLeg := parseLeg(leg)
		parsed.Legs = append(parsed.Legs, parsedLeg)
	}

	// Adjust interchange count (don't count walks)
	walkCount := 0
	for _, leg := range parsed.Legs {
		if leg.IsWalk {
			walkCount++
		}
	}
	parsed.Interchanges = len(parsed.Legs) - 1 - walkCount
	if parsed.Interchanges < 0 {
		parsed.Interchanges = 0
	}

	return parsed
}

// parseLeg converts API leg to our parsed format
func parseLeg(leg Leg) ParsedLeg {
	parsed := ParsedLeg{
		Origin:      cleanStopName(leg.Origin.Name),
		OriginTrack: leg.Origin.Track,
		Destination: cleanStopName(leg.Destination.Name),
		DestTrack:   leg.Destination.Track,
		IsWalk:      leg.Type == "WALK" || leg.Type == "TRSF",
		Distance:    leg.Dist,
		Direction:   leg.Direction,
	}

	// Parse times
	parsed.DepartureTime = parseDateTime(leg.Origin.Date, leg.Origin.Time)
	parsed.ArrivalTime = parseDateTime(leg.Destination.Date, leg.Destination.Time)

	// Real-time times
	if leg.Origin.RtTime != "" {
		rt := parseDateTime(leg.Origin.RtDate, leg.Origin.RtTime)
		parsed.RtDeparture = &rt
	}
	if leg.Destination.RtTime != "" {
		rt := parseDateTime(leg.Destination.RtDate, leg.Destination.RtTime)
		parsed.RtArrival = &rt
	}

	// Real-time tracks
	if leg.Origin.RtTrack != "" {
		parsed.OriginTrack = leg.Origin.RtTrack
	}
	if leg.Destination.RtTrack != "" {
		parsed.DestTrack = leg.Destination.RtTrack
	}

	// Line and category info
	if leg.Product != nil {
		parsed.Line = leg.Product.Line
		if parsed.Line == "" {
			parsed.Line = leg.Product.Num
		}
		parsed.Category = leg.Product.CatOut
		parsed.Operator = leg.Product.Operator
	} else {
		// Fallback to leg info
		parsed.Line = leg.Number
		parsed.Category = leg.Category
		parsed.Operator = leg.Operator
	}

	return parsed
}

// parseDuration parses ISO 8601 duration (e.g., "PT1H30M")
func parseDuration(s string) time.Duration {
	if s == "" {
		return 0
	}

	// Remove "PT" prefix
	s = strings.TrimPrefix(s, "PT")

	var total time.Duration

	// Parse hours
	if idx := strings.Index(s, "H"); idx >= 0 {
		if hours, err := strconv.Atoi(s[:idx]); err == nil {
			total += time.Duration(hours) * time.Hour
		}
		s = s[idx+1:]
	}

	// Parse minutes
	if idx := strings.Index(s, "M"); idx >= 0 {
		if mins, err := strconv.Atoi(s[:idx]); err == nil {
			total += time.Duration(mins) * time.Minute
		}
		s = s[idx+1:]
	}

	// Parse seconds
	if idx := strings.Index(s, "S"); idx >= 0 {
		if secs, err := strconv.Atoi(s[:idx]); err == nil {
			total += time.Duration(secs) * time.Second
		}
	}

	return total
}

// parseDateTime parses date and time strings
func parseDateTime(date, timeStr string) time.Time {
	if date == "" || timeStr == "" {
		return time.Time{}
	}

	// Try parsing with different formats
	combined := date + " " + timeStr

	// Try "2006-01-02 15:04:05"
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", combined, time.Local); err == nil {
		return t
	}

	// Try "2006-01-02 15:04"
	if t, err := time.ParseInLocation("2006-01-02 15:04", combined, time.Local); err == nil {
		return t
	}

	return time.Time{}
}

// cleanStopName removes redundant suffixes from stop names
func cleanStopName(name string) string {
	// Remove common suffixes like "(Sundsvall kn)"
	re := regexp.MustCompile(`\s*\([^)]+\)\s*$`)
	name = re.ReplaceAllString(name, "")

	// Trim whitespace
	return strings.TrimSpace(name)
}

// FormatTrips formats trips for display
func FormatTrips(origin, dest string, trips []ParsedTrip) string {
	var sb strings.Builder

	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	sb.WriteString(fmt.Sprintf(" 🚆 %s → %s\n", origin, dest))
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	if len(trips) == 0 {
		sb.WriteString("  Inga resor hittades.\n")
		sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		return sb.String()
	}

	for i, trip := range trips {
		if i > 0 {
			sb.WriteString("\n  ─────────────────────────────────────────────────────────────────\n\n")
		}

		// Trip header
		firstLeg := trip.Legs[0]
		lastLeg := trip.Legs[len(trip.Legs)-1]

		depTime := firstLeg.DepartureTime.Format("15:04")
		arrTime := lastLeg.ArrivalTime.Format("15:04")

		// Check for delays
		depStr := depTime
		if firstLeg.RtDeparture != nil && !firstLeg.RtDeparture.Equal(firstLeg.DepartureTime) {
			delay := firstLeg.RtDeparture.Sub(firstLeg.DepartureTime).Minutes()
			if delay > 0 {
				depStr = fmt.Sprintf("%s (+%.0f)", depTime, delay)
			}
		}

		arrStr := arrTime
		if lastLeg.RtArrival != nil && !lastLeg.RtArrival.Equal(lastLeg.ArrivalTime) {
			delay := lastLeg.RtArrival.Sub(lastLeg.ArrivalTime).Minutes()
			if delay > 0 {
				arrStr = fmt.Sprintf("%s (+%.0f)", arrTime, delay)
			}
		}

		durationMin := int(trip.Duration.Minutes())
		hours := durationMin / 60
		mins := durationMin % 60
		var durationStr string
		if hours > 0 {
			durationStr = fmt.Sprintf("%d tim %d min", hours, mins)
		} else {
			durationStr = fmt.Sprintf("%d min", mins)
		}

		changes := ""
		if trip.Interchanges == 0 {
			changes = "direkt"
		} else if trip.Interchanges == 1 {
			changes = "1 byte"
		} else {
			changes = fmt.Sprintf("%d byten", trip.Interchanges)
		}

		sb.WriteString(fmt.Sprintf("  %s → %s   (%s, %s)\n\n", depStr, arrStr, durationStr, changes))

		// Legs
		for _, leg := range trip.Legs {
			if leg.IsWalk {
				walkMin := int(leg.ArrivalTime.Sub(leg.DepartureTime).Minutes())
				if leg.Distance > 0 {
					sb.WriteString(fmt.Sprintf("  🚶 Gå %d m (%d min)\n", leg.Distance, walkMin))
				} else {
					sb.WriteString(fmt.Sprintf("  🚶 Byte (%d min)\n", walkMin))
				}
			} else {
				emoji := GetCategoryEmoji(leg.Category)
				lineStr := leg.Line
				if lineStr == "" {
					lineStr = GetCategoryName(leg.Category)
				}

				depTime := leg.DepartureTime.Format("15:04")
				arrTime := leg.ArrivalTime.Format("15:04")

				trackInfo := ""
				if leg.OriginTrack != "" {
					trackInfo = fmt.Sprintf(" [spår %s]", leg.OriginTrack)
				}

				sb.WriteString(fmt.Sprintf("  %s %s %s%s\n", emoji, lineStr, leg.Origin, trackInfo))
				sb.WriteString(fmt.Sprintf("     %s → %s mot %s\n", depTime, arrTime, leg.Direction))

				if leg.Operator != "" && leg.Operator != leg.Line {
					sb.WriteString(fmt.Sprintf("     (%s)\n", leg.Operator))
				}
			}
		}

		// Google Maps link
		if mapsURL := GenerateTripMapsURL(trip); mapsURL != "" {
			sb.WriteString(fmt.Sprintf("\n  🗺️  %s\n", mapsURL))
		}
	}

	sb.WriteString("\n")
	sb.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	return sb.String()
}

// GenerateTripMapsURL creates a Google Maps URL for a transit trip
func GenerateTripMapsURL(trip ParsedTrip) string {
	if len(trip.Legs) == 0 {
		return ""
	}

	// Get origin from first leg and destination from last leg
	origin := trip.Legs[0].Origin
	dest := trip.Legs[len(trip.Legs)-1].Destination

	params := url.Values{}
	params.Set("api", "1")
	params.Set("origin", origin)
	params.Set("destination", dest)
	params.Set("travelmode", "transit")

	// Add waypoints for transfers
	if len(trip.Legs) > 1 {
		var waypoints []string
		for i := 0; i < len(trip.Legs)-1; i++ {
			waypoints = append(waypoints, trip.Legs[i].Destination)
		}
		if len(waypoints) > 0 {
			params.Set("waypoints", strings.Join(waypoints, "|"))
		}
	}

	return "https://www.google.com/maps/dir/?" + params.Encode()
}
