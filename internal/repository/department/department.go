package department

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/dimashiro/test_mediasoft/internal/model"
	"github.com/dimashiro/test_mediasoft/internal/model/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	// tables
	departmentTable = "departments"
	itemsTable      = "items"
	orderItemsTable = "order_items"
)

type DepartmentRepo interface {
	GetByID(ctx context.Context, departmentID string) (model.Department, error)
	Create(ctx context.Context, dto *dto.CreateDepartment) (model.Department, error)
	Update(ctx context.Context, dto *dto.UpdateDepartment) error
	// GetHierarchy()
	// Update()
	// Delete()
}

type Repository struct {
	db *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

func (r *Repository) GetByID(ctx context.Context, departmentID string) (model.Department, error) {
	dp := model.Department{}
	query, args, err := sq.
		Select("department_id", "department_name", "department_path").
		From(departmentTable).
		Where(sq.Eq{"department_id": departmentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return dp, fmt.Errorf("can't build query: %s", err.Error())
	}
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return dp, fmt.Errorf("can't select orders: %s", err.Error())
	}
	if rows.Next() {
		err := rows.Scan(&dp.ID, &dp.Name, &dp.Path)
		if err != nil {
			return dp, fmt.Errorf("can't scan order: %s", err.Error())
		}
	}
	defer rows.Close()
	return dp, nil
}

func (r *Repository) Create(ctx context.Context, dto *dto.CreateDepartment) (model.Department, error) {
	dpParent := model.Department{}
	//TODO move to validation later
	if dto.ParentID != "" {
		if _, err := uuid.Parse(dto.ParentID); err == nil {
			dpParent, err = r.GetByID(ctx, dto.ParentID)
			if err != nil {
				return model.Department{}, fmt.Errorf("parent not found: %s", err.Error())
			}
		} else {
			return model.Department{}, fmt.Errorf("wrond parent id: %s", err.Error())
		}
	}
	uuid := GenerateID()
	path := dpParent.Path
	if path != "" {
		path = path + "."
	}

	path = path + strings.ReplaceAll(uuid, "-", "_")
	query, args, err := sq.
		Insert(departmentTable).
		Columns("department_id", "department_name", "department_path").
		Values(
			uuid,
			dto.Name,
			path,
		).Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return model.Department{}, fmt.Errorf("can't build sql: %s", err.Error())
	}

	// insert into departments table.
	var newDepartment model.Department
	err = r.db.QueryRow(ctx, query, args...).
		Scan(&newDepartment.ID, &newDepartment.Name, &newDepartment.Path)
	if err != nil {
		return model.Department{}, fmt.Errorf("can't scan department: %w", err)
	}

	return newDepartment, nil
}

func (r *Repository) Update(ctx context.Context, dto *dto.UpdateDepartment) error {
	//TODO move to validation later
	dp := model.Department{}
	if _, err := uuid.Parse(dto.ID); err == nil {
		dp, err = r.GetByID(ctx, *dto.ParentID)
		if err != nil {
			return fmt.Errorf("department not found: %s", err.Error())
		}
	}

	dpParent := model.Department{}
	//TODO move to validation later
	if dto.ParentID != nil {
		if _, err := uuid.Parse(*dto.ParentID); err == nil {
			dpParent, err = r.GetByID(ctx, *dto.ParentID)
			if err != nil {
				return fmt.Errorf("parent not found: %s", err.Error())
			}
		} else {
			return fmt.Errorf("wrond parent id: %s", err.Error())
		}
		dp.Parent = &dpParent
		dp.Path = dpParent.Path + "." + strings.ReplaceAll(dp.ID, "-", "_")
	}

	if dto.Name != nil {
		dp.Name = *dto.Name
	}

	query, args, err := sq.
		Update(departmentTable).
		Set("department_name", dp.Name).
		Set("department_path", dp.Path).
		Where(sq.Eq{"department_id": dp.ID}).PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("can't build sql: %s", err.Error())
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("sql exec err: %w", err)
	}

	return nil
}

func GenerateID() string {
	return uuid.NewString()
}