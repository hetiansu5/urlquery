package query

type builderOptions struct {
	u UrlEncoder
}

type BuilderOption interface {
	apply(*builderOptions)
}

type urlEncoderOption struct {
	u UrlEncoder
}

func (u urlEncoderOption) apply(opts *builderOptions) {
	opts.u = u.u
}

//support customized urlEncoder
func WithUrlEncoder(u UrlEncoder) BuilderOption {
	return urlEncoderOption{u: u}
}
