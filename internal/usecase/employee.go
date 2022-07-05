package usecase

import (
	"context"

	"github.com/dimashiro/test_mediasoft/internal/model"
	"github.com/dimashiro/test_mediasoft/internal/model/dto"
	"github.com/dimashiro/test_mediasoft/internal/repository/department"
	"github.com/dimashiro/test_mediasoft/internal/repository/employee"
	"go.uber.org/zap"
)

type Employee struct {
	log   *zap.SugaredLogger
	rEmpl employee.EmployeeRepo
	rDptm department.DepartmentRepo
}

func NewEmployee(log *zap.SugaredLogger, rEmpl employee.EmployeeRepo, rDptm department.DepartmentRepo) *Employee {
	return &Employee{log: log, rEmpl: rEmpl, rDptm: rDptm}
}

func (e Employee) CreateEmployee(dto *dto.CreateEmployee) (model.Employee, error) {
	empl, err := e.rEmpl.Create(context.Background(), dto)
	if err != nil {
		return model.Employee{}, err
	}
	return empl, nil
}

func (e Employee) GetAllEmployees() ([]model.Employee, error) {
	return e.rEmpl.GetAll(context.Background())
}

func (e Employee) UpdateEmployee(dto *dto.UpdateEmployee) error {
	return e.rEmpl.Update(context.Background(), dto)
}

func (e Employee) DeleteEmployee(dto *dto.DeleteEmployee) error {
	return e.rEmpl.Delete(context.Background(), dto)
}
