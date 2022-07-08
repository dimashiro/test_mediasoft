package usecase

import (
	"context"
	"strings"

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

func (d Department) CreateDepartment(ctx context.Context, dto *dto.CreateDepartment) (model.Department, error) {
	dp, err := d.rDptm.Create(ctx, dto)
	if err != nil {
		return model.Department{}, err
	}
	return dp, nil
}

func (d Department) UpdateDepartment(ctx context.Context, dto *dto.UpdateDepartment) error {
	err := d.rDptm.Update(context.Background(), dto)
	if err != nil {
		return err
	}
	return nil
}

func (d Department) HierarchyDepartment(ctx context.Context) ([]*model.Department, error) {
	m, err := d.rDptm.Hierarchy(ctx)
	if err != nil {
		return []*model.Department{}, err
	}
	dps := []*model.Department{}
	for _, v := range m {
		if strings.ReplaceAll(v.ID, "-", "_") == v.Path {
			dps = append(dps, v)
		}
	}
	return dps, nil
}

func (d Department) GetAllDepartments(ctx context.Context) ([]dto.ViewAllDepartments, error) {
	return d.rDptm.GetAll(ctx)
}

func (d Department) DeleteDepartment(ctx context.Context, dto *dto.DeleteDepartment) error {
	return d.rDptm.Delete(ctx, dto)
}

func (d Department) GetEmployeesByDepartment(ctx context.Context, departmentID string) ([]model.Employee, error) {
	return d.rEmpl.GetByDepartment(ctx, departmentID)
}

func (d Department) GetEmployeesInDepartmentHierarchy(ctx context.Context, departmentID string) ([]model.Employee, error) {
	dp, err := d.rDptm.GetByID(ctx, departmentID)
	if err != nil {
		return []model.Employee{}, err
	}
	return d.rEmpl.GetInDepartmentHierarchy(ctx, dp)
}
