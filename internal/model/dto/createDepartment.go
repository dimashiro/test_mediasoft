package dto

type CreateDepartment struct {
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
}
