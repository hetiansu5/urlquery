package urlquery

import "net/url"

//support user-defined global urlEncoder and have default.

var (
	gUrlEncoder UrlEncoder
	cUrlEncoder commonUrlEncoder
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
}

type commonUrlEncoder struct{}

func (u commonUrlEncoder) Escape(s string) string {
	return url.QueryEscape(s)
}
