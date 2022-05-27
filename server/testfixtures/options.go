package testfixtures

func WithCacheDisabled() Option {
	return func(opt *FixtureOptions) {
		opt.DisableCache = true
	}
}
