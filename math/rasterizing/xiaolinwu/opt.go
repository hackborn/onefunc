package xiaolinwu

type Option func(*rasterizer)

// WithDebug adds some print debugging.
func WithDebug() Option {
	return func(r *rasterizer) {
		r.debug = true
	}
}
