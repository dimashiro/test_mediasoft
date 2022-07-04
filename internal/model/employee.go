package model

type Employee struct {
	ID        uint64
	Name      string
	Surname   string
	BirthYear uint16
	Groups    []Department
}
