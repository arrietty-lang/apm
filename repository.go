package apm

type Repository struct {
	Host    string
	Author  string
	Name    string
	Version string
	Deps    []*Dependencies
}
