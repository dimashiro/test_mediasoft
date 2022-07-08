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

func (e Employee) CreateEmployee(ctx context.Context, dto *dto.CreateEmployee) (model.Employee, error) {
	empl, err := e.rEmpl.Create(ctx, dto)
	if err != nil {
		return model.Employee{}, err
	}
	return empl, nil
}

func (e Employee) GetAllEmployees(ctx context.Context) ([]model.Employee, error) {
	return e.rEmpl.GetAll(ctx)
}

func (e Employee) UpdateEmployee(ctx context.Context, dto *dto.UpdateEmployee) error {
	return e.rEmpl.Update(ctx, dto)
}

func (e Employee) DeleteEmployee(ctx context.Context, dto *dto.DeleteEmployee) error {
	return e.rEmpl.Delete(ctx, dto)
}
