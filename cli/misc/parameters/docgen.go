package misc_parameters

type DocGen struct {
	DocsOutputDir    string
	DocusaurusFolder string
}

func NewDocGen() *DocGen {
	return &DocGen{}
}
