package flows_parameters

type Configure struct {
	AnalyticsEnabled bool
	Endpoint         string
	Global           bool
}

func NewConfigure() *Configure {
	return &Configure{}
}
