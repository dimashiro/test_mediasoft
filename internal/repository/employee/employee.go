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
	departmentTable         = "departments"
)

type EmployeeRepo interface {
	GetByID(ctx context.Context, employeeID string) (model.Employee, error)
	Create(ctx context.Context, dto *dto.CreateEmployee) (model.Employee, error)
	GetAll(ctx context.Context) ([]model.Employee, error)
	Delete(ctx context.Context, dto *dto.DeleteEmployee) error
	Update(ctx context.Context, dto *dto.UpdateEmployee) error
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
		err := rows.Scan(&employee.ID, &employee.Name, &employee.Surname, &employee.BirthYear)
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
	err = tx.QueryRow(ctx, query, args...).
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

func (r *Repository) GetAll(ctx context.Context) ([]model.Employee, error) {
	var empls []model.Employee
	emplMap := make(map[string]model.Employee)
	query, args, err := sq.
		Select("employee_id", "employee_name", "employee_surname",
			"employee_birthyear", "d.department_id", "d.department_name",
			"d.department_path").
		From(employeeTable).
		Join(employeeDepartmentTable + " USING (employee_id)").
		Join(departmentTable + " AS d USING (department_id)").
		ToSql()
	if err != nil {
		return empls, fmt.Errorf("can't build query: %s", err.Error())
	}
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return empls, fmt.Errorf("can't select employees: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		empl := model.Employee{}
		dptm := model.Department{}
		err := rows.Scan(&empl.ID, &empl.Name, &empl.Surname,
			&empl.BirthYear, &dptm.ID, &dptm.Name, &dptm.Path)
		if err != nil {
			return empls, fmt.Errorf("can't scan employee: %s", err.Error())
		}
		if e, ok := emplMap[empl.ID]; ok {
			e.Departments = append(e.Departments, dptm)
			emplMap[empl.ID] = e
		} else {
			empl.Departments = append(empl.Departments, dptm)
			emplMap[empl.ID] = empl
		}
	}

	for _, e := range emplMap {
		empls = append(empls, e)
	}
	return empls, nil
}

func (r *Repository) Update(ctx context.Context, dto *dto.UpdateEmployee) error {
	employee := model.Employee{}
	if _, err := uuid.Parse(dto.ID); err == nil {
		employee, err = r.GetByID(ctx, dto.ID)
		if err != nil {
			return fmt.Errorf("employee not found: %s", err.Error())
		}
	}
	if dto.Name != nil {
		employee.Name = *dto.Name
	}
	if dto.Surname != nil {
		employee.Surname = *dto.Surname
	}
	if dto.BirthYear != nil {
		employee.BirthYear = *dto.BirthYear
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("can't create tx: %s", err.Error())
	}

	//update employee
	query, args, err := sq.
		Update(employeeTable).
		Set("employee_name", employee.Name).
		Set("employee_surname", employee.Surname).
		Set("employee_birthyear", employee.BirthYear).
		Where(sq.Eq{"employee_id": employee.ID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("can't build sql: %s", err.Error())
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			return fmt.Errorf("rollback err: %s, err: %s", rollbackErr.Error(), err.Error())
		}
		return fmt.Errorf("can't update employee: %w", err)
	}

	//get old departments
	query, args, err = sq.Select("department_id").
		From(employeeDepartmentTable).
		Where(sq.Eq{"employee_id": dto.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("can't build sql: %s", err.Error())
	}
	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("can't get old departments: %s", err.Error())
	}
	defer rows.Close()
	oldDepartments := make(map[string]bool)
	for rows.Next() {
		var dpID string
		err := rows.Scan(&dpID)
		if err != nil {
			return fmt.Errorf("can't scan department id: %s", err.Error())
		}
		oldDepartments[dpID] = false
	}

	//insert into connection table
	//TODO check for not existing departments in dto
	qBuilder := sq.
		Insert(employeeDepartmentTable).
		Columns("employee_id", "department_id")
	needInsert := false
	for _, dpID := range dto.Departments {
		if _, err := uuid.Parse(dpID); err != nil {
			return fmt.Errorf("wrond department id: %s", err.Error())
		}
		_, ok := oldDepartments[dpID]
		if !ok {
			qBuilder = qBuilder.Values(dto.ID, dpID)
			needInsert = true
		} else {
			oldDepartments[dpID] = true
		}
	}

	if needInsert {
		query, args, err = qBuilder.PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return fmt.Errorf("can't build sql: %s", err.Error())
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				return fmt.Errorf("rollback err: %s, err: %s", rollbackErr.Error(), err.Error())
			}

			return err
		}
	}

	//delete from connection table
	var forDelete []string

	for dID, inNew := range oldDepartments {
		if !inNew {
			forDelete = append(forDelete, dID)
		}
	}
	if len(forDelete) > 0 {
		query, args, err = sq.
			Delete(employeeDepartmentTable).
			Where(sq.Eq{"department_id": forDelete, "employee_id": employee.ID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return fmt.Errorf("can't build query: %s", err.Error())
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				return fmt.Errorf("rollback err: %s, err: %s", rollbackErr.Error(), err.Error())
			}

			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("can't commit tx: %s", err.Error())
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, dto *dto.DeleteEmployee) error {
	if _, err := uuid.Parse(dto.ID); err == nil {
		_, err := r.GetByID(ctx, dto.ID)
		if err != nil {
			return fmt.Errorf("employee not found: %s", err.Error())
		}
	} else {
		return fmt.Errorf("wrong id: %s", err.Error())
	}

	query, args, err := sq.
		Delete(employeeTable).
		Where(sq.Eq{"employee_id": dto.ID}).
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
