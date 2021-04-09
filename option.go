package urlquery

// options for encoder or parser
type options struct {
	urlEncoder     UrlEncoder
	needEmptyValue bool
}

// An Option is a func type for applying diff options
type Option func(*options)

// WithUrlEncoder is supposed customized URL-Encoder option
func WithUrlEncoder(u UrlEncoder) Option {
	return func(ops *options) {
		ops.urlEncoder = u
	}
}

// WithNeedEmptyValue is supposed to control whether to ignore zero-value.
// It just happen to the element directly in structure, not including map slice array
// default:false, meaning ignore zero-value
func WithNeedEmptyValue(c bool) Option {
	return func(ops *options) {
		ops.needEmptyValue = c
	}
}
