package display

import (
	"fmt"
	"strings"

	"transport/internal/api"
	"transport/internal/tz"
)

// FormatDepartures formats departures for display
func (f *Formatter) FormatDepartures(site *api.Site, mode, towards string, departures []api.Departure) string {
	var sb strings.Builder

	// Header
	sb.WriteString(strings.Repeat("━", lineWidth) + "\n")

	modeIcon := getModeIcon(mode)
	header := fmt.Sprintf(" %s Nästa %s från %s", modeIcon, getModeName(mode), site.Name)
	if towards != "" {
		header += fmt.Sprintf(" mot %s", towards)
	}

	sb.WriteString(header + "\n")
	sb.WriteString(strings.Repeat("━", lineWidth) + "\n\n")

	if len(departures) == 0 {
		sb.WriteString("  Inga avgångar hittades.\n\n")
		sb.WriteString(strings.Repeat("━", lineWidth) + "\n")
		return sb.String()
	}

	// Each departure
	for _, dep := range departures {
		sb.WriteString(f.FormatDeparture(dep))
	}

	sb.WriteString("\n" + strings.Repeat("━", lineWidth) + "\n")

	return sb.String()
}

// FormatDeparture formats a single departure
func (f *Formatter) FormatDeparture(dep api.Departure) string {
	var sb strings.Builder

	icon := getDepartureIcon(dep.Line.TransportMode)
	line := dep.Line.Designation
	destination := dep.Destination

	// Calculate time until departure
	inMinutes := f.calculateMinutesUntil(dep.Expected)
	timeDisplay := formatTimeUntil(inMinutes)

	// Get actual time
	actualTime := parseTimeOnly(dep.Expected)

	// Format: 🚌 117  Brommaplan           om 7 min (11:09)
	sb.WriteString(fmt.Sprintf("  %s %-4s %-25s %s (%s)\n",
		icon,
		line,
		truncate(destination, 25),
		timeDisplay,
		actualTime))

	// Show platform/stop point if available
	if dep.StopPoint.Designation != "" {
		sb.WriteString(fmt.Sprintf("         Läge %s\n", dep.StopPoint.Designation))
	}

	return sb.String()
}

// calculateMinutesUntil calculates minutes until the given time
func (f *Formatter) calculateMinutesUntil(timeStr string) int {
	if timeStr == "" {
		return 0
	}

	// API returns Swedish local time without timezone indicator
	t, err := tz.ParseStockholm("2006-01-02T15:04:05", timeStr)
	if err != nil {
		return 0
	}

	diff := tz.Now().Sub(t)
	minutes := -int(diff.Minutes())

	if minutes < 0 {
		return 0
	}
	return minutes
}

// formatTimeUntil formats the time until departure
func formatTimeUntil(minutes int) string {
	if minutes <= 0 {
		return "   Nu    "
	} else if minutes == 1 {
		return " om 1 min"
	} else if minutes < 60 {
		return fmt.Sprintf("om %2d min", minutes)
	} else {
		hours := minutes / 60
		mins := minutes % 60
		if mins == 0 {
			return fmt.Sprintf("om %d h   ", hours)
		}
		return fmt.Sprintf("om %dh%02dm ", hours, mins)
	}
}

// parseTimeOnly extracts HH:MM from a time string
func parseTimeOnly(timeStr string) string {
	if timeStr == "" || len(timeStr) < 16 {
		return "     "
	}
	return timeStr[11:16]
}

// getModeIcon returns the icon for a transport mode
func getModeIcon(mode string) string {
	switch strings.ToUpper(mode) {
	case "BUS":
		return "🚌"
	case "METRO":
		return "🚇"
	case "TRAIN":
		return "🚂"
	case "TRAM":
		return "🚊"
	case "SHIP":
		return "⛴️"
	default:
		return "🚍"
	}
}

// getDepartureIcon returns the icon based on API transport mode
func getDepartureIcon(transportMode string) string {
	switch transportMode {
	case api.TransportModeBus:
		return "🚌"
	case api.TransportModeMetro:
		return "🚇"
	case api.TransportModeTrain:
		return "🚂"
	case api.TransportModeTram:
		return "🚊"
	case api.TransportModeShip:
		return "⛴️"
	default:
		return "🚍"
	}
}

// getModeName returns the Swedish name for a transport mode
func getModeName(mode string) string {
	switch strings.ToUpper(mode) {
	case "BUS":
		return "buss"
	case "METRO":
		return "tunnelbana"
	case "TRAIN":
		return "tåg"
	case "TRAM":
		return "spårvagn"
	case "SHIP":
		return "båt"
	default:
		return "avgång"
	}
}

// truncate truncates a string to max length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}
