package urlquery

import "net/url"

// support user-defined global urlEncoder and have default.
var (
	gUrlEncoder UrlEncoder
	cUrlEncoder DefaultUrlEncoder
)

// set global URL-Encode
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

// define URL-Encoder interface
type UrlEncoder interface {
	Escape(s string) string
	UnEscape(s string) (string, error)
}

// default url encoder
type DefaultUrlEncoder struct{}

// escape text
func (u DefaultUrlEncoder) Escape(s string) string {
	return url.QueryEscape(s)
}

// unescape text
func (u DefaultUrlEncoder) UnEscape(s string) (string, error) {
	return url.QueryUnescape(s)
}
