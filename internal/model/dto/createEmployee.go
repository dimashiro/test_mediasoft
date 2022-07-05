package dto

type CreateEmployee struct {
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	BirthYear   int      `json:"birthyear"`
	Departments []string `json:"departments_ids"`
}
