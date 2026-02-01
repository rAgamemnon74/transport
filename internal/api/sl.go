package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL          = "https://journeyplanner.integration.sl.se/v2"
	transportBaseURL = "https://transport.integration.sl.se/v1"
	defaultTimeout   = 10 * time.Second
)

// Client handles communication with the SL Journey Planner API
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new SL API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: baseURL,
	}
}

// SearchStops finds stops matching the given query
func (c *Client) SearchStops(query string) ([]Location, error) {
	params := url.Values{}
	params.Set("type_sf", "any")
	params.Set("name_sf", query)
	params.Set("any_obj_filter_sf", "2") // Only stops (not addresses/POI)

	reqURL := fmt.Sprintf("%s/stop-finder?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search stops: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var result StopFinderResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for API errors only if no locations found
	if len(result.Locations) == 0 {
		for _, msg := range result.SystemMessages {
			if msg.Type == "error" && msg.Text != "" {
				return nil, fmt.Errorf("API error: %s", msg.Text)
			}
		}
	}

	return result.Locations, nil
}

// TripOptions contains options for trip planning
type TripOptions struct {
	Time          time.Time // Departure or arrival time
	ArriveBy      bool      // If true, time is arrival time
	MaxChanges    int       // Maximum number of transfers (-1 for unlimited)
	NumResults    int       // Number of results to return (1-6)
	Language      string    // "sv" or "en"
}

// DefaultTripOptions returns default options for trip planning
func DefaultTripOptions() TripOptions {
	return TripOptions{
		Time:       time.Now(),
		ArriveBy:   false,
		MaxChanges: -1,
		NumResults: 3,
		Language:   "sv",
	}
}

// PlanTrip finds journeys between origin and destination
func (c *Client) PlanTrip(originID, destID string, opts TripOptions) ([]Journey, error) {
	params := url.Values{}
	params.Set("type_origin", "any")
	params.Set("name_origin", originID)
	params.Set("type_destination", "any")
	params.Set("name_destination", destID)
	params.Set("calc_number_of_trips", fmt.Sprintf("%d", opts.NumResults))

	if opts.Language != "" {
		params.Set("language", opts.Language)
	}

	if !opts.Time.IsZero() {
		params.Set("itd_date", opts.Time.Format("20060102"))
		params.Set("itd_time", opts.Time.Format("1504"))
	}

	if opts.ArriveBy {
		params.Set("itd_trip_date_time_dep_arr", "arr")
	} else {
		params.Set("itd_trip_date_time_dep_arr", "dep")
	}

	if opts.MaxChanges >= 0 {
		params.Set("maxChanges", fmt.Sprintf("%d", opts.MaxChanges))
	}

	reqURL := fmt.Sprintf("%s/trips?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to plan trip: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var result TripsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for API errors (but not if we have journeys)
	if len(result.Journeys) == 0 {
		for _, msg := range result.SystemMessages {
			if msg.Type == "error" {
				return nil, fmt.Errorf("API error: %s", msg.Text)
			}
		}
	}

	return result.Journeys, nil
}

// PlanTripByName is a convenience method that resolves stop names first
func (c *Client) PlanTripByName(origin, dest string, opts TripOptions) ([]Journey, error) {
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

	// Use the best match for each
	return c.PlanTrip(originStops[0].ID, destStops[0].ID, opts)
}

// SearchSites searches for sites by name using the Transport API
func (c *Client) SearchSites(query string) ([]Site, error) {
	reqURL := fmt.Sprintf("%s/sites", transportBaseURL)

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sites: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var allSites []Site
	if err := json.NewDecoder(resp.Body).Decode(&allSites); err != nil {
		return nil, fmt.Errorf("failed to decode sites: %w", err)
	}

	// Filter by query (matches name, aliases, or abbreviation)
	var matches []Site
	for _, site := range allSites {
		if site.MatchesQuery(query) {
			matches = append(matches, site)
		}
	}

	// Sort: exact matches first, then by name length (shorter = better)
	sortSiteMatches(matches, query)

	return matches, nil
}

// sortSiteMatches sorts sites with better matches first
func sortSiteMatches(sites []Site, query string) {
	queryLower := toLower(query)

	// Simple bubble sort (sites list is small after filtering)
	for i := 0; i < len(sites); i++ {
		for j := i + 1; j < len(sites); j++ {
			// Prefer exact name match
			iExact := toLower(sites[i].Name) == queryLower
			jExact := toLower(sites[j].Name) == queryLower

			if jExact && !iExact {
				sites[i], sites[j] = sites[j], sites[i]
				continue
			}

			// Prefer exact alias match
			for _, alias := range sites[j].Alias {
				if toLower(alias) == queryLower {
					if !iExact {
						sites[i], sites[j] = sites[j], sites[i]
					}
					break
				}
			}
		}
	}
}

// GetDepartures fetches departures from a site
func (c *Client) GetDepartures(siteID int) ([]Departure, error) {
	reqURL := fmt.Sprintf("%s/sites/%d/departures", transportBaseURL, siteID)

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch departures: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var result DeparturesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode departures: %w", err)
	}

	return result.Departures, nil
}

// GetNextDepartures gets filtered departures by mode and destination
func (c *Client) GetNextDepartures(location string, mode string, towards string, count int) ([]Departure, *Site, error) {
	// Find the site
	sites, err := c.SearchSites(location)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search sites: %w", err)
	}
	if len(sites) == 0 {
		// Try a fuzzy search to suggest alternatives
		suggestions := c.findSimilarSites(location)
		if len(suggestions) > 0 {
			return nil, nil, fmt.Errorf("no sites found for '%s'. Did you mean: %s", location, formatSuggestions(suggestions))
		}
		return nil, nil, fmt.Errorf("no sites found for '%s'", location)
	}

	site := sites[0]

	// Get departures
	departures, err := c.GetDepartures(site.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get departures: %w", err)
	}

	// Normalize mode for comparison
	modeUpper := toUpper(mode)
	towardsLower := toLower(towards)

	// Filter departures
	var filtered []Departure
	for _, dep := range departures {
		// Filter by transport mode
		if modeUpper != "" && dep.Line.TransportMode != modeUpper {
			continue
		}

		// Filter by destination (substring match)
		if towardsLower != "" && !containsIgnoreCase(dep.Destination, towardsLower) {
			continue
		}

		filtered = append(filtered, dep)
		if len(filtered) >= count {
			break
		}
	}

	return filtered, &site, nil
}

// Helper functions for case-insensitive matching
func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		} else {
			b[i] = c
		}
	}
	return string(b)
}

func toUpper(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		} else {
			b[i] = c
		}
	}
	return string(b)
}

func containsIgnoreCase(s, substr string) bool {
	sLower := toLower(s)
	return len(sLower) >= len(substr) && (sLower == substr || indexString(sLower, substr) >= 0)
}

func indexString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// findSimilarSites finds sites that partially match the query
func (c *Client) findSimilarSites(query string) []string {
	reqURL := fmt.Sprintf("%s/sites", transportBaseURL)

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var allSites []Site
	if err := json.NewDecoder(resp.Body).Decode(&allSites); err != nil {
		return nil
	}

	// Extract first word of query for fuzzy matching
	queryLower := toLower(query)
	words := splitWords(queryLower)
	if len(words) == 0 {
		return nil
	}
	firstWord := words[0]

	type match struct {
		name  string
		score int // lower is better
	}
	var matches []match

	for _, site := range allSites {
		nameLower := toLower(site.Name)
		// Check if first word matches start of site name
		if len(nameLower) >= len(firstWord) && nameLower[:len(firstWord)] == firstWord {
			// Prefer exact first-word matches (e.g., "Spånga" over "Spångavägen")
			score := len(site.Name)
			if nameLower == firstWord {
				score = 0 // Exact match gets best score
			}
			matches = append(matches, match{site.Name, score})
		}
	}

	// Sort by score (shorter/exact matches first)
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[j].score < matches[i].score {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}

	// Return top 3
	var suggestions []string
	for i := 0; i < len(matches) && i < 3; i++ {
		suggestions = append(suggestions, matches[i].name)
	}

	return suggestions
}

// splitWords splits a string into words
func splitWords(s string) []string {
	var words []string
	var current []byte
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\t' {
			if len(current) > 0 {
				words = append(words, string(current))
				current = nil
			}
		} else {
			current = append(current, s[i])
		}
	}
	if len(current) > 0 {
		words = append(words, string(current))
	}
	return words
}

// formatSuggestions formats a list of suggestions for display
func formatSuggestions(suggestions []string) string {
	if len(suggestions) == 0 {
		return ""
	}
	if len(suggestions) == 1 {
		return suggestions[0]
	}
	result := suggestions[0]
	for i := 1; i < len(suggestions)-1; i++ {
		result += ", " + suggestions[i]
	}
	result += " or " + suggestions[len(suggestions)-1]
	return result
}
