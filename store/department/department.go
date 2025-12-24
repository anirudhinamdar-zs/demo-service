package department

import (
	"database/sql"
	"demo-service/models/department"
	"errors"

	"gofr.dev/pkg/gofr"
)

type Department struct {
	db *sql.DB
}

func Init(db *sql.DB) *Department {
	return &Department{
		db: db,
	}
}

func (d *Department) Create(ctx *gofr.Context, dep *department.Department) (*department.Department, error) {
	query := `
		INSERT INTO departments (code, name, floor, description)
		VALUES (?, ?, ?, ?)
	`

	// Code is same as name enum (Option A)
	result, err := d.db.ExecContext(
		ctx,
		query,
		dep.Code,
		dep.Name,
		dep.Floor,
		dep.Description,
	)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return nil, sql.ErrNoRows
	}

	return &department.Department{
		Code:        dep.Code,
		Name:        dep.Name,
		Floor:       dep.Floor,
		Description: dep.Description,
	}, nil
}

func (d *Department) Get(ctx *gofr.Context) ([]*department.Department, error) {
	query := `
		SELECT code, name, floor, description
		FROM departments
	`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []*department.Department

	for rows.Next() {
		var dep department.Department

		err := rows.Scan(
			&dep.Code,
			&dep.Name,
			&dep.Floor,
			&dep.Description,
		)
		if err != nil {
			return nil, err
		}

		departments = append(departments, &dep)
	}

	return departments, nil
}

func (d *Department) GetByCode(ctx *gofr.Context, code string) (*department.Department, error) {
	query := `
		SELECT code, name, floor, description
		FROM departments
		WHERE code = ?
	`

	var dep department.Department

	err := d.db.QueryRowContext(ctx, query, code).Scan(
		&dep.Code,
		&dep.Name,
		&dep.Floor,
		&dep.Description,
	)
	if err != nil {
		return nil, err
	}

	return &dep, nil
}

func (d *Department) Update(
	ctx *gofr.Context,
	code string,
	dep *department.NewDepartment,
) (*department.Department, error) {

	query := `
		UPDATE departments
		SET
			name = ?,
			floor = ?,
			description = ?
		WHERE code = ?
	`

	result, err := d.db.ExecContext(
		ctx,
		query,
		dep.Name,
		dep.Floor,
		dep.Description,
		code,
	)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, sql.ErrNoRows
	}

	return &department.Department{
		Code:        code,
		Name:        dep.Name,
		Floor:       dep.Floor,
		Description: dep.Description,
	}, nil
}

func (d *Department) Delete(ctx *gofr.Context, code string) (string, error) {
	checkQuery := `
		SELECT COUNT(1)
		FROM employees
		WHERE department = ?
	`

	var count int
	err := d.db.QueryRowContext(ctx, checkQuery, code).Scan(&count)
	if err != nil {
		return "", err
	}

	if count > 0 {
		return "", errors.New("department has employees mapped and cannot be deleted")
	}

	deleteQuery := `
		DELETE FROM departments
		WHERE code = ?
	`

	result, err := d.db.ExecContext(ctx, deleteQuery, code)
	if err != nil {
		return "", err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rows == 0 {
		return "", sql.ErrNoRows
	}

	return "Department deleted successfully", nil
}

func (d *Department) ExistsByName(
	ctx *gofr.Context,
	name string,
	excludeCode *string,
) (bool, error) {

	query := `SELECT COUNT(1) FROM departments WHERE name = ?`
	args := []interface{}{name}

	if excludeCode != nil {
		query += ` AND code != ?`
		args = append(args, *excludeCode)
	}

	var count int
	err := d.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
