package model

type Department struct {
	ID       string
	Name     string
	Path     string
	Children []*Department
}
