package urlquery

// options for encoder or parser
type options struct {
	queryEncoder   QueryEncoder
	needEmptyValue bool
}

// An Option is a func type for applying diff options
type Option func(*options)

// WithQueryEncoder is supposed customized query encoder option
func WithQueryEncoder(u QueryEncoder) Option {
	return func(ops *options) {
		ops.queryEncoder = u
	}
}

// WithNeedEmptyValue is supposed to control whether to ignore zero value.
// It just happen to the element directly in structure, not including map slice array
// default:false, meaning ignore zero-value
func WithNeedEmptyValue(c bool) Option {
	return func(ops *options) {
		ops.needEmptyValue = c
	}
}
