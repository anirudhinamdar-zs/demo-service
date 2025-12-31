package department

import (
	"database/sql"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"demo-service/models/department"
	"demo-service/store"
)

type Department struct {
}

func Init() store.Department {
	return &Department{}
}

func (d *Department) Create(ctx *gofr.Context, dep *department.Department) (*department.Department, error) {
	createQuery := `INSERT INTO departments (code, name, floor, description) VALUES (?, ?, ?, ?)`
	result, err := ctx.DB().ExecContext(
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

func (d *Department) Get(ctx *gofr.Context) ([]*department.Department, error) {
	query := `
		SELECT code, name, floor, description
		FROM departments
	`

	rows, err := ctx.DB().QueryContext(ctx, query)
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

func (d *Department) GetByCode(ctx *gofr.Context, code string) (*department.Department, error) {
	query := `
		SELECT code, name, floor, description
		FROM departments
		WHERE code = ?
	`

	var dep department.Department

	err := ctx.DB().QueryRowContext(ctx, query, code).Scan(
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
		SET name = ?, floor = ?, description = ?
		WHERE code = ?
	`

	result, err := ctx.DB().ExecContext(
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

func (d *Department) Delete(ctx *gofr.Context, code string) (string, error) {
	checkQuery := `
		SELECT COUNT(1)
		FROM employees
		WHERE department = ?
	`

	var count int
	if err := ctx.DB().QueryRowContext(ctx, checkQuery, code).Scan(&count); err != nil {
		return "", err
	}

	if count > 0 {
		return "", errors.Error("department has employees mapped")
	}

	deleteQuery := `
		DELETE FROM departments
		WHERE code = ?
	`

	result, err := ctx.DB().ExecContext(ctx, deleteQuery, code)
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
	if err := ctx.DB().QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil
}
