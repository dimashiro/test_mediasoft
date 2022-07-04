package model

type Department struct {
	ID      string
	Name    string
	Path    string
	Parent  *Department
	Persons []*Employee
}
