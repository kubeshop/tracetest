package types

type ResourceList[T any] struct {
	Count int `json:"count"`
	Items []T `json:"items"`
}
