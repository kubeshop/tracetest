package resources_parameters

type List struct {
	Take          int32
	Skip          int32
	SortBy        string
	SortDirection string
	All           bool
}

func NewList() *List {
	return &List{}
}
