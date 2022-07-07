package usecase

import (
	"context"

	"github.com/dimashiro/test_mediasoft/internal/model"
	"github.com/dimashiro/test_mediasoft/internal/model/dto"
	"github.com/dimashiro/test_mediasoft/internal/repository/department"
	"github.com/dimashiro/test_mediasoft/internal/repository/employee"
	"go.uber.org/zap"
)

// Usecase responsible for saving request.
type Department struct {
	log   *zap.SugaredLogger
	rEmpl employee.EmployeeRepo
	rDptm department.DepartmentRepo
}

func NewDepartment(log *zap.SugaredLogger, rDptm department.DepartmentRepo, rEmpl employee.EmployeeRepo) *Department {
	return &Department{log: log, rDptm: rDptm, rEmpl: rEmpl}
}

func (d Department) CreateDepartment(dto *dto.CreateDepartment) (model.Department, error) {
	dp, err := d.rDptm.Create(context.Background(), dto)
	if err != nil {
		return model.Department{}, err
	}
	return dp, nil
}

func (d Department) UpdateDepartment(dto *dto.UpdateDepartment) error {
	err := d.rDptm.Update(context.Background(), dto)
	if err != nil {
		return err
	}
	return nil
}

func (d Department) HierarchyDepartment() ([]model.Department, error) {
	return d.rDptm.Hierarchy(context.Background())
}

func (d Department) DeleteDepartment(dto *dto.DeleteDepartment) error {
	return d.rDptm.Delete(context.Background(), dto)
}

func (d Department) GetEmployeesByDepartment(departmentID string) ([]model.Employee, error) {
	return d.rEmpl.GetByDepartment(context.Background(), departmentID)
}
