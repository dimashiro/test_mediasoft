package department

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/dimashiro/test_mediasoft/internal/model"
	"github.com/dimashiro/test_mediasoft/internal/model/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	// tables
	departmentTable = "departments"
)

type DepartmentRepo interface {
	GetByID(ctx context.Context, departmentID string) (model.Department, error)
	Create(ctx context.Context, dto *dto.CreateDepartment) (model.Department, error)
	Update(ctx context.Context, dto *dto.UpdateDepartment) error
	Hierarchy(ctx context.Context) ([]model.Department, error)
	GetAll(ctx context.Context) ([]dto.ViewAllDepartments, error)
	Delete(ctx context.Context, dto *dto.DeleteDepartment) error
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
		return dp, fmt.Errorf("can't select departments: %s", err.Error())
	}
	if rows.Next() {
		err := rows.Scan(&dp.ID, &dp.Name, &dp.Path)
		if err != nil {
			return dp, fmt.Errorf("can't scan department: %s", err.Error())
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
		dp, err = r.GetByID(ctx, dto.ID)
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

func (r *Repository) Hierarchy(ctx context.Context) ([]model.Department, error) {
	var dps []model.Department
	query, args, err := sq.
		Select("department_id", "department_name", "department_path").
		From(departmentTable).
		ToSql()
	if err != nil {
		return dps, fmt.Errorf("can't build query: %s", err.Error())
	}
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return dps, fmt.Errorf("can't select departments: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		dp := model.Department{}
		err := rows.Scan(&dp.ID, &dp.Name, &dp.Path)
		if err != nil {
			return dps, fmt.Errorf("can't scan department: %s", err.Error())
		}
		dps = append(dps, dp)
	}
	return dps, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]dto.ViewAllDepartments, error) {
	var dps []dto.ViewAllDepartments

	sql := `select
	d.department_id,
    d.department_name,
	(select count(*) from employee_department ed 
	    where ed.department_id=d.department_id) as count_empl,
	(select count(*) from employee_department ed 
	    where ed.department_id in (select d1.department_id from departments d1 where d1.department_path <@ d.department_path)) as count_with_child_empl
	from departments d
	order by d.department_name`

	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return dps, fmt.Errorf("can't select departments: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		dp := dto.ViewAllDepartments{}
		err := rows.Scan(&dp.ID, &dp.Name, &dp.EmployeesAmount, &dp.EmployeesAmountInHierarchy)
		if err != nil {
			return dps, fmt.Errorf("can't scan department: %s", err.Error())
		}
		dps = append(dps, dp)
	}
	return dps, nil
}

func (r *Repository) Delete(ctx context.Context, dto *dto.DeleteDepartment) error {
	dp := model.Department{}
	if _, err := uuid.Parse(dto.ID); err == nil {
		dp, err = r.GetByID(ctx, dto.ID)
		if err != nil {
			return fmt.Errorf("department not found: %s", err.Error())
		}
	} else {
		return fmt.Errorf("wrong id: %s", err.Error())
	}

	sql := "SELECT department_path FROM departments WHERE department_path <@ $1 AND department_path != $1 LIMIT 1"
	var dpPath string
	err := r.db.QueryRow(ctx, sql, dp.Path).Scan(&dpPath)
	if err == nil {
		return errors.New("Cannot delete department with descendants.")
	} else {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("query error: %s", err.Error())
		}
	}

	query, args, err := sq.
		Delete(departmentTable).
		Where(sq.Eq{"department_id": dto.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("can't build query: %s", err.Error())
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
