package display

import (
	"fmt"
	"strings"
	"time"

	"transport/internal/api"
)

const (
	lineWidth = 70
)

// Formatter handles formatting of journey results
type Formatter struct {
	Language string
}

// NewFormatter creates a new formatter
func NewFormatter(lang string) *Formatter {
	if lang == "" {
		lang = "sv"
	}
	return &Formatter{Language: lang}
}

// FormatJourneys formats multiple journeys for display
func (f *Formatter) FormatJourneys(origin, dest string, journeys []api.Journey) string {
	var sb strings.Builder

	// Header
	sb.WriteString(strings.Repeat("━", lineWidth) + "\n")
	sb.WriteString(fmt.Sprintf(" %s → %s", origin, dest))
	sb.WriteString(strings.Repeat(" ", max(0, lineWidth-len(origin)-len(dest)-6-len(time.Now().Format("Mon 2 Jan")))))
	sb.WriteString(time.Now().Format("Mon 2 Jan") + "\n")
	sb.WriteString(strings.Repeat("━", lineWidth) + "\n\n")

	// Each journey
	for i, journey := range journeys {
		sb.WriteString(f.FormatJourney(i+1, journey))
		sb.WriteString("\n")
	}

	sb.WriteString(strings.Repeat("━", lineWidth) + "\n")

	return sb.String()
}

// FormatJourney formats a single journey
func (f *Formatter) FormatJourney(num int, journey api.Journey) string {
	var sb strings.Builder

	// Journey header
	duration := formatDuration(journey.TripDuration)
	changes := f.formatChanges(journey.Interchanges)
	header := fmt.Sprintf(" Resa %d", num)
	stats := fmt.Sprintf("%s │ %s", duration, changes)
	padding := max(0, lineWidth-len(header)-len(stats)-2)
	sb.WriteString(fmt.Sprintf("%s%s%s\n", header, strings.Repeat(" ", padding), stats))
	sb.WriteString(strings.Repeat("─", lineWidth) + "\n")

	// Each leg
	for i, leg := range journey.Legs {
		sb.WriteString(f.FormatLeg(leg, i == 0, i == len(journey.Legs)-1))
	}

	return sb.String()
}

// FormatLeg formats a single leg of a journey
func (f *Formatter) FormatLeg(leg api.Leg, isFirst, isLast bool) string {
	var sb strings.Builder

	// Get times
	depTime := parseTime(leg.Origin.DepartureTimePlanned)
	arrTime := parseTime(leg.Destination.ArrivalTimePlanned)

	// Origin
	originName := leg.Origin.GetStopName()
	platform := formatPlatform(leg.Origin.GetPlatform())

	sb.WriteString(fmt.Sprintf("  %s  %-30s  %s\n",
		depTime,
		originName,
		platform))

	// Transport line or walking
	if leg.Transportation != nil && !leg.Transportation.IsWalking() {
		icon := TransportIcon(leg.Transportation)
		lineName := leg.Transportation.GetLineName()
		direction := ""
		if d := leg.Transportation.GetDirection(); d != "" {
			direction = " → " + d
		}
		sb.WriteString(fmt.Sprintf("    │    %s %s%s\n", icon, lineName, direction))
	} else {
		// Walking
		walkMin := leg.Duration / 60
		if walkMin < 1 {
			walkMin = 1
		}
		sb.WriteString(fmt.Sprintf("    │    🚶 Gång %d min\n", walkMin))
	}

	// Destination (only show if last leg or if there's a transfer)
	if isLast {
		destName := leg.Destination.GetStopName()
		destPlatform := formatPlatform(leg.Destination.GetPlatform())
		sb.WriteString(fmt.Sprintf("  %s  %-30s  %s\n",
			arrTime,
			destName,
			destPlatform))
	}

	return sb.String()
}

// formatDuration formats seconds to "X min" or "X h Y min"
func formatDuration(seconds int) string {
	mins := seconds / 60
	if mins < 60 {
		return fmt.Sprintf("%d min", mins)
	}
	hours := mins / 60
	remainMins := mins % 60
	if remainMins == 0 {
		return fmt.Sprintf("%d h", hours)
	}
	return fmt.Sprintf("%d h %d min", hours, remainMins)
}

// formatChanges formats number of transfers
func (f *Formatter) formatChanges(n int) string {
	if f.Language == "en" {
		if n == 0 {
			return "direct"
		} else if n == 1 {
			return "1 change"
		}
		return fmt.Sprintf("%d changes", n)
	}

	// Swedish
	if n == 0 {
		return "direkt"
	} else if n == 1 {
		return "1 byte"
	}
	return fmt.Sprintf("%d byten", n)
}

// parseTime parses API time format and returns HH:MM in local time
func parseTime(timeStr string) string {
	if timeStr == "" {
		return "     "
	}

	// Try parsing ISO format with Z (UTC): 2025-01-31T08:15:00Z
	t, err := time.Parse("2006-01-02T15:04:05Z", timeStr)
	if err != nil {
		// Try without Z
		t, err = time.Parse("2006-01-02T15:04:05", timeStr)
		if err != nil {
			// Try RFC3339
			t, err = time.Parse(time.RFC3339, timeStr)
			if err != nil {
				// Return as-is if parsing fails
				if len(timeStr) >= 16 {
					return timeStr[11:16] // Extract HH:MM from ISO string
				}
				return timeStr
			}
		}
	}

	// Convert to local time
	return t.Local().Format("15:04")
}

// shortName returns the disassembled name if available, otherwise the full name
func shortName(name, disassembled string) string {
	if disassembled != "" {
		return disassembled
	}
	// Remove "Stockholm, " prefix if present
	if strings.HasPrefix(name, "Stockholm, ") {
		return name[11:]
	}
	return name
}

// formatPlatform formats platform information
func formatPlatform(platform string) string {
	if platform == "" {
		return ""
	}
	return fmt.Sprintf("Spår %s", platform)
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
