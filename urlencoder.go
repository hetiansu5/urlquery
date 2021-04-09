package urlquery

import "net/url"

var (
	// global URL-Encoder, priority: global > default
	gUrlEncoder UrlEncoder
	// default URL-Encoder
	cUrlEncoder DefaultUrlEncoder
)

// SetGlobalUrlEncoder set global URL-Encode
func SetGlobalUrlEncoder(u UrlEncoder) {
	gUrlEncoder = u
}

// get URL-Encoder
func getUrlEncoder() UrlEncoder {
	if gUrlEncoder == nil {
		return cUrlEncoder
	}
	return gUrlEncoder
}

// A UrlEncoder is a interface implementing Escape and UnEscape method
type UrlEncoder interface {
	Escape(s string) string
	UnEscape(s string) (string, error)
}

// A DefaultUrlEncoder is a default URL-Encoder
type DefaultUrlEncoder struct{}

// Escape text
func (u DefaultUrlEncoder) Escape(s string) string {
	return url.QueryEscape(s)
}

// UnEscape text
func (u DefaultUrlEncoder) UnEscape(s string) (string, error) {
	return url.QueryUnescape(s)
}
