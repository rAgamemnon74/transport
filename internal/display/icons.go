package display

import "transport/internal/api"

// TransportIcon returns the appropriate icon for a transport type
func TransportIcon(t *api.Transportation) string {
	if t == nil {
		return "🚶" // Default to walking
	}

	// Check product class first
	if t.Product != nil {
		switch t.Product.Class {
		case api.ProductClassMetro:
			return "🚇"
		case api.ProductClassTrain:
			return "🚂"
		case api.ProductClassTram:
			return "🚊"
		case api.ProductClassBus:
			return "🚌"
		case api.ProductClassFerry:
			return "⛴️"
		}
	}

	// Fallback based on name patterns
	name := t.Name
	if t.Product != nil && t.Product.Name != "" {
		name = t.Product.Name + " " + name
	}

	switch {
	case containsAny(name, "tunnelbana", "metro", "t-bana"):
		return "🚇"
	case containsAny(name, "pendeltåg", "commuter", "train", "tåg"):
		return "🚂"
	case containsAny(name, "spårvagn", "tram", "tvärbanan", "lidingöbanan", "nockebybanan", "saltsjöbanan"):
		return "🚊"
	case containsAny(name, "buss", "bus", "ersättnings"):
		return "🚌"
	case containsAny(name, "båt", "ferry", "färja", "waxholm", "sjövägen"):
		return "⛴️"
	case containsAny(name, "gång", "walk"):
		return "🚶"
	default:
		return "🚍"
	}
}

// TransportTypeName returns a Swedish name for the transport type
func TransportTypeName(t *api.Transportation) string {
	if t == nil {
		return "Gång"
	}

	if t.Product != nil {
		switch t.Product.Class {
		case api.ProductClassMetro:
			return "Tunnelbana"
		case api.ProductClassTrain:
			return "Pendeltåg"
		case api.ProductClassTram:
			return "Spårvagn"
		case api.ProductClassBus:
			return "Buss"
		case api.ProductClassFerry:
			return "Båt"
		}
	}

	return "Transport"
}

// containsAny checks if s contains any of the substrings (case-insensitive)
func containsAny(s string, substrs ...string) bool {
	lower := toLower(s)
	for _, sub := range substrs {
		if contains(lower, toLower(sub)) {
			return true
		}
	}
	return false
}

// toLower converts string to lowercase (simple ASCII)
func toLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}

// contains checks if s contains substr
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || indexString(s, substr) >= 0)
}

// indexString returns index of substr in s
func indexString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
