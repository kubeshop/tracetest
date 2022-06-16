package testfixtures

func WithCacheDisabled() Option {
	return func(opt *FixtureOptions) {
		opt.DisableCache = true
	}
}

func WithArguments(args interface{}) Option {
	return func(opt *FixtureOptions) {
		opt.Arguments = args
	}
}
