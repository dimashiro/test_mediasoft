package usecase

import (
	"context"

	"github.com/dimashiro/test_mediasoft/internal/model"
	"github.com/dimashiro/test_mediasoft/internal/model/dto"
	"github.com/dimashiro/test_mediasoft/internal/repository/department"
	"go.uber.org/zap"
)

// Usecase responsible for saving request.
type Department struct {
	log  *zap.SugaredLogger
	repo department.DepartmentRepo
}

func NewDepartment(log *zap.SugaredLogger, repo department.DepartmentRepo) *Department {
	return &Department{log: log, repo: repo}
}

func (d Department) CreateDepartment(dto *dto.CreateDepartment) (model.Department, error) {
	dp, err := d.repo.Create(context.Background(), dto)
	if err != nil {
		return model.Department{}, err
	}
	return dp, nil
}

func (d Department) UpdateDepartment(dto *dto.UpdateDepartment) error {
	err := d.repo.Update(context.Background(), dto)
	if err != nil {
		return err
	}
	return nil
}
