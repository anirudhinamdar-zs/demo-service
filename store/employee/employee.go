package employee

import (
	"database/sql"

	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"demo-service/models/employee"
	"demo-service/store"
)

type Employee struct {
}

func Init() store.Employee {
	return &Employee{}
}

func (e *Employee) Create(ctx *gofr.Context, emp *employee.NewEmployee) (*employee.Employee, error) {
	query := `INSERT INTO employees (name, email, phone_number, dob, major, city, department) VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := ctx.DB().ExecContext(
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
	args := []interface{}{}

	if filter.ID != nil {
		baseGetEmployeesQuery += ` AND id = ?`
		args = append(args, *filter.ID)
	}

	if filter.Name != nil {
		baseGetEmployeesQuery += ` AND name LIKE ?`
		args = append(args, "%"+*filter.Name+"%")
	}

	if filter.Department != nil {
		baseGetEmployeesQuery += ` AND department = ?`
		args = append(args, *filter.Department)
	}

	rows, err := ctx.DB().QueryContext(ctx, baseGetEmployeesQuery, args...)
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
	var emp employee.Employee

	err := ctx.DB().QueryRowContext(ctx, queryByID, employeeId).Scan(
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
	result, err := ctx.DB().ExecContext(
		ctx,
		updateQueryById,
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
	result, err := ctx.DB().ExecContext(ctx, deleteQueryById, employeeId)
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
	args := []interface{}{email}

	if excludeID != nil {
		selectQueryForEmail += ` AND id != ?`
		args = append(args, *excludeID)
	}

	var count int
	err := ctx.DB().QueryRowContext(ctx, selectQueryForEmail, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (e *Employee) CountByDepartment(
	ctx *gofr.Context,
	deptCode string,
) (int, error) {
	var count int
	err := ctx.DB().QueryRowContext(ctx, selectQueryCountByDepartment, deptCode).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
