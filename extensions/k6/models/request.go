package models

type Metadata map[string]string

type Request struct {
	Method   string
	URL      string
	ID       string
	Metadata Metadata
}
