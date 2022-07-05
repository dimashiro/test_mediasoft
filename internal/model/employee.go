package model

type Employee struct {
	ID        string
	Name      string
	Surname   string
	BirthYear uint16
	Groups    []Department
}
