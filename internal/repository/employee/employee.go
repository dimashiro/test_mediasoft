package employee

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/dimashiro/test_mediasoft/internal/model"
	"github.com/dimashiro/test_mediasoft/internal/model/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	// tables
	employeeTable           = "employees"
	employeeDepartmentTable = "employee_department"
)

type EmployeeRepo interface {
	GetByID(ctx context.Context, employeeID string) (model.Employee, error)
	Create(ctx context.Context, dto *dto.CreateEmployee) (model.Employee, error)
	// Update(ctx context.Context, dto *dto.UpdateDepartment) error
	// GetAll(ctx context.Context) ([]model.Department, error)
	// Delete(ctx context.Context, dto *dto.DeleteDepartment) error
}

type Repository struct {
	db *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

func (r *Repository) GetByID(ctx context.Context, employeeID string) (model.Employee, error) {
	employee := model.Employee{}
	query, args, err := sq.
		Select("*").
		From(employeeTable).
		Where(sq.Eq{"employee_id": employeeID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return employee, fmt.Errorf("can't build query: %s", err.Error())
	}
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return employee, fmt.Errorf("can't select employees: %s", err.Error())
	}
	if rows.Next() {
		err := rows.Scan(&employee)
		if err != nil {
			return employee, fmt.Errorf("can't scan employee: %s", err.Error())
		}
	}
	defer rows.Close()
	return employee, nil
}

func (r *Repository) Create(ctx context.Context, dto *dto.CreateEmployee) (model.Employee, error) {
	employee := model.Employee{}
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return employee, fmt.Errorf("can't create tx: %s", err.Error())
	}

	uuidEmployee := uuid.NewString()
	query, args, err := sq.
		Insert(employeeTable).
		Columns("employee_id", "employee_name", "employee_surname", "employee_birthyear").
		Values(
			uuidEmployee,
			dto.Name,
			dto.Surname,
			dto.BirthYear,
		).Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return employee, fmt.Errorf("can't build sql: %s", err.Error())
	}
	// insert into empoyee table.
	err = r.db.QueryRow(ctx, query, args...).
		Scan(&employee.ID, &employee.Name, &employee.Surname, &employee.BirthYear)
	if err != nil {
		return model.Employee{}, fmt.Errorf("can't scan Employee: %w", err)
	}

	//insert into connection table
	qBuilder := sq.
		Insert(employeeDepartmentTable).
		Columns("employee_id", "department_id")
	for _, dpID := range dto.Departments {
		if _, err := uuid.Parse(dpID); err != nil {
			return employee, fmt.Errorf("wrond department id: %s", err.Error())
		}
		qBuilder = qBuilder.Values(uuidEmployee, dpID)
	}
	query, args, err = qBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return employee, fmt.Errorf("can't build sql: %s", err.Error())
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			return employee, fmt.Errorf("rollback err: %s, err: %s", rollbackErr.Error(), err.Error())
		}

		return employee, err
	}

	if err := tx.Commit(ctx); err != nil {
		return employee, fmt.Errorf("can't commit tx: %s", err.Error())
	}

	return employee, nil
}
