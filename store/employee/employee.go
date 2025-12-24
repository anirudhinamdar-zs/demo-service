package employee

import (
	"demo-service/models/employee"

	"database/sql"

	"gofr.dev/pkg/gofr"
)

type Employee struct {
	db *sql.DB
}

func Init(db *sql.DB) *Employee {
	return &Employee{
		db: db,
	}
}

func (e *Employee) Create(ctx *gofr.Context, emp *employee.NewEmployee) (*employee.Employee, error) {
	query := `
		INSERT INTO employees
			(name, email, phone_number, dob, major, city, department)
		VALUES
			(?, ?, ?, ?, ?, ?, ?)
	`

	result, err := e.db.ExecContext(
		ctx,
		query,
		emp.Name,
		emp.Email,
		emp.PhoneNumber,
		emp.DOB,
		emp.Major,
		emp.City,
		emp.Department,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &employee.Employee{
		ID:          int(id),
		Name:        emp.Name,
		Email:       emp.Email,
		PhoneNumber: emp.PhoneNumber,
		DOB:         emp.DOB,
		Major:       emp.Major,
		City:        emp.City,
		Department:  emp.Department,
	}, nil
}

func (e *Employee) Get(
	ctx *gofr.Context,
	filter employee.Filter,
) ([]*employee.Employee, error) {

	baseQuery := `
		SELECT
			id, name, email, phone_number, dob, major, city, department
		FROM employees
		WHERE 1 = 1
	`

	args := []interface{}{}

	if filter.ID != nil {
		baseQuery += ` AND id = ?`
		args = append(args, *filter.ID)
	}

	if filter.Name != nil {
		baseQuery += ` AND name LIKE ?`
		args = append(args, "%"+*filter.Name+"%")
	}

	if filter.Department != nil {
		baseQuery += ` AND department = ?`
		args = append(args, *filter.Department)
	}

	rows, err := e.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []*employee.Employee

	for rows.Next() {
		var emp employee.Employee

		err := rows.Scan(
			&emp.ID,
			&emp.Name,
			&emp.Email,
			&emp.PhoneNumber,
			&emp.DOB,
			&emp.Major,
			&emp.City,
			&emp.Department,
		)
		if err != nil {
			return nil, err
		}

		employees = append(employees, &emp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

func (e *Employee) GetById(ctx *gofr.Context, employeeId int) (*employee.Employee, error) {
	query := `
		SELECT id, name, email, phone_number, dob, major, city, department
		FROM employees
		WHERE id = ?
	`

	var emp employee.Employee

	err := e.db.QueryRowContext(ctx, query, employeeId).Scan(
		&emp.ID,
		&emp.Name,
		&emp.Email,
		&emp.PhoneNumber,
		&emp.DOB,
		&emp.Major,
		&emp.City,
		&emp.Department,
	)
	if err != nil {
		return nil, err
	}

	return &emp, nil
}

func (e *Employee) Update(ctx *gofr.Context, employeeId int, emp *employee.NewEmployee) (*employee.Employee, error) {

	query := `
		UPDATE employees
		SET
			name = ?,
			email = ?,
			phone_number = ?,
			dob = ?,
			major = ?,
			city = ?,
			department = ?
		WHERE id = ?
	`

	result, err := e.db.ExecContext(
		ctx,
		query,
		emp.Name,
		emp.Email,
		emp.PhoneNumber,
		emp.DOB,
		emp.Major,
		emp.City,
		emp.Department,
		employeeId,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return &employee.Employee{
		ID:          employeeId,
		Name:        emp.Name,
		Email:       emp.Email,
		PhoneNumber: emp.PhoneNumber,
		DOB:         emp.DOB,
		Major:       emp.Major,
		City:        emp.City,
		Department:  emp.Department,
	}, nil
}

func (e *Employee) Delete(ctx *gofr.Context, employeeId int) (string, error) {
	query := `DELETE FROM employees WHERE id = ?`

	result, err := e.db.ExecContext(ctx, query, employeeId)
	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowsAffected == 0 {
		return "", sql.ErrNoRows
	}

	return "Employee deleted successfully", nil
}

func (e *Employee) ExistsByEmail(
	ctx *gofr.Context,
	email string,
	excludeID *int,
) (bool, error) {

	query := `SELECT COUNT(1) FROM employees WHERE email = ?`
	args := []interface{}{email}

	if excludeID != nil {
		query += ` AND id != ?`
		args = append(args, *excludeID)
	}

	var count int
	err := e.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (e *Employee) CountByDepartment(
	ctx *gofr.Context,
	deptCode string,
) (int, error) {

	query := `
		SELECT COUNT(1)
		FROM employees
		WHERE department = ?
	`

	var count int
	err := e.db.QueryRowContext(ctx, query, deptCode).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
