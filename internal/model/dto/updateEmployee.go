package dto

type UpdateEmployee struct {
	ID          string   `json:"id"`
	Name        *string  `json:"name"`
	Surname     *string  `json:"surname"`
	BirthYear   *uint16  `json:"birthyear"`
	Departments []string `json:"departments_ids"`
}
