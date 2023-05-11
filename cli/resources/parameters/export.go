package resources_parameters

type Export struct {
	ResourceID   string
	ResourceFile string
}

func NewExport() *Export {
	return &Export{}
}
