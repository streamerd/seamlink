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
	// Function to store link clicks (e.g., to database)
	StoreLinkClick func(click SeamlinkClick) error
	// Optional: domains to exclude from tracking
	ExcludeDomains []string
}

// DefaultConfig returns the default Seamlink configuration
func DefaultConfig() SeamlinkConfig {
	return SeamlinkConfig{
		StoreLinkClick: func(click SeamlinkClick) error {
			// Default implementation: log to console
			return nil
		},
		ExcludeDomains: []string{},
	}
}
