package query

import "net/url"

//support user-defined global urlEncoder and have default.

var (
	globalUrlEncoder UrlEncoder
	commonUrlEncoder CommonUrlEncoder
)

func SetGlobalUrlEncoder(u UrlEncoder) {
	globalUrlEncoder = u
}

func GetUrlEncoder() UrlEncoder {
	if globalUrlEncoder == nil {
		return commonUrlEncoder
	}
	return globalUrlEncoder
}

type UrlEncoder interface {
	Escape(s string) string
}

type CommonUrlEncoder struct{}

func (u CommonUrlEncoder) Escape(s string) string {
	return url.QueryEscape(s)
}
