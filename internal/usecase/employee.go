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
