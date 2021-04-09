package urlquery

import "net/url"

var (
	// global query encoder, priority: global > default
	gQueryEncoder QueryEncoder
	// default query encoder
	defaultQueryEncoder DefaultQueryEncoder
)

// SetGlobalQueryEncoder set global query encoder
func SetGlobalQueryEncoder(u QueryEncoder) {
	gQueryEncoder = u
}

// get query encoder
func getQueryEncoder() QueryEncoder {
	if gQueryEncoder == nil {
		return defaultQueryEncoder
	}
	return gQueryEncoder
}

// A QueryEncoder is a interface implementing Escape and UnEscape method
type QueryEncoder interface {
	Escape(s string) string
	UnEscape(s string) (string, error)
}

// A DefaultQueryEncoder is a default URL-Encoder
type DefaultQueryEncoder struct{}

// Escape text
func (u DefaultQueryEncoder) Escape(s string) string {
	return url.QueryEscape(s)
}

// UnEscape text
func (u DefaultQueryEncoder) UnEscape(s string) (string, error) {
	return url.QueryUnescape(s)
}
