package department

import (
	"context"
	"database/sql"
	"errors"

	"demo-service/models/department"
)

type Department struct {
	db *sql.DB
}

func Init(db *sql.DB) *Department {
	return &Department{db: db}
}

const createQuery = `
	INSERT INTO departments (code, name, floor, description)
	VALUES (?, ?, ?, ?)
`

func (d *Department) Create(ctx context.Context, dep *department.Department) (*department.Department, error) {
	result, err := d.db.ExecContext(
		ctx,
		createQuery,
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

	return dep, nil
}

func (d *Department) Get(ctx context.Context) ([]*department.Department, error) {
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

		if err := rows.Scan(
			&dep.Code,
			&dep.Name,
			&dep.Floor,
			&dep.Description,
		); err != nil {
			return nil, err
		}

		departments = append(departments, &dep)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return departments, nil
}

func (d *Department) GetByCode(ctx context.Context, code string) (*department.Department, error) {
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
	ctx context.Context,
	code string,
	dep *department.NewDepartment,
) (*department.Department, error) {

	query := `
		UPDATE departments
		SET name = ?, floor = ?, description = ?
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
	if err != nil || rows == 0 {
		return nil, sql.ErrNoRows
	}

	return &department.Department{
		Code:        code,
		Name:        dep.Name,
		Floor:       dep.Floor,
		Description: dep.Description,
	}, nil
}

func (d *Department) Delete(ctx context.Context, code string) (string, error) {
	checkQuery := `
		SELECT COUNT(1)
		FROM employees
		WHERE department = ?
	`

	var count int
	if err := d.db.QueryRowContext(ctx, checkQuery, code).Scan(&count); err != nil {
		return "", err
	}

	if count > 0 {
		return "", errors.New("department has employees mapped")
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
	if err != nil || rows == 0 {
		return "", sql.ErrNoRows
	}

	return "Department deleted successfully", nil
}

func (d *Department) ExistsByName(
	ctx context.Context,
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
	if err := d.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}
