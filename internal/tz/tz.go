// Package tz provides the Swedish timezone for the transport service.
package tz

import "time"

// Stockholm is the Europe/Stockholm timezone used by Swedish transport APIs.
var Stockholm *time.Location

func init() {
	var err error
	Stockholm, err = time.LoadLocation("Europe/Stockholm")
	if err != nil {
		// Fallback: if tzdata is missing, use fixed CET offset.
		Stockholm = time.FixedZone("CET", 3600)
	}
}

// Now returns the current time in Stockholm.
func Now() time.Time {
	return time.Now().In(Stockholm)
}

// ParseStockholm parses a time string as Stockholm local time.
func ParseStockholm(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, Stockholm)
}
