package urlquery

import "net/url"

//support user-defined global urlEncoder and have default.

var (
	gUrlEncoder UrlEncoder
	cUrlEncoder DefaultUrlEncoder
)

func SetGlobalUrlEncoder(u UrlEncoder) {
	gUrlEncoder = u
}

func getUrlEncoder() UrlEncoder {
	if gUrlEncoder == nil {
		return cUrlEncoder
	}
	return gUrlEncoder
}

type UrlEncoder interface {
	Escape(s string) string
	UnEscape(s string) (string, error)
}

//default url encoder
type DefaultUrlEncoder struct{}

func (u DefaultUrlEncoder) Escape(s string) string {
	return url.QueryEscape(s)
}

func (u DefaultUrlEncoder) UnEscape(s string) (string, error) {
	return url.QueryUnescape(s)
}
