package seamlink

import "time"

// SeamlinkClick represents a tracked outbound click
type SeamlinkClick struct {
	URL       string    `json:"url"`
	Referrer  string    `json:"referrer"`
	UserAgent string    `json:"userAgent"`
	Timestamp time.Time `json:"timestamp"`
}

// SeamlinkConfig middleware configuration
type SeamlinkConfig struct {
	StoreLinkClick func(click SeamlinkClick) error
	StorePageVisit func(visit PageVisit) error // New callback for page visits
	ExcludeDomains []string
}

// DefaultConfig returns default configuration
func DefaultConfig() SeamlinkConfig {
	return SeamlinkConfig{
		StoreLinkClick: func(click SeamlinkClick) error { return nil },
		StorePageVisit: func(visit PageVisit) error { return nil },
		ExcludeDomains: []string{},
	}
}
