package department

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"demo-service/models/department"
)

func Initialize(t *testing.T) (*gomock.Controller, *gofr.Context, sqlmock.Sqlmock, error) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	return ctrl, ctx, mock, err
}

func TestCreate(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()

	query := `INSERT INTO departments (code, name, floor, description) VALUES (?, ?, ?, ?)`

	input := &department.Department{
		Code:        "IT",
		Name:        "IT",
		Floor:       1,
		Description: "desc",
	}

	resp := &department.Department{
		Code:        "IT",
		Name:        "IT",
		Floor:       1,
		Description: "desc",
	}

	tests := []struct {
		desc    string
		result  *department.Department
		mock    *sqlmock.ExpectedExec
		mockErr error
	}{
		{
			desc:    "db error",
			result:  nil,
			mockErr: errors.DB{Err: err},
			mock: mock.ExpectExec(query).
				WithArgs(
					input.Code,
					input.Name,
					input.Floor,
					input.Description,
				).
				WillReturnResult(nil).
				WillReturnError(errors.DB{Err: err}),
		},
		{
			desc:    "no rows affected",
			result:  nil,
			mockErr: sql.ErrNoRows,
			mock: mock.ExpectExec(query).
				WithArgs(
					input.Code,
					input.Name,
					input.Floor,
					input.Description,
				).
				WillReturnResult(sqlmock.NewResult(0, 0)),
		},
		{
			desc:    "success",
			result:  resp,
			mockErr: nil,
			mock: mock.ExpectExec(query).
				WithArgs(
					input.Code,
					input.Name,
					input.Floor,
					input.Description,
				).
				WillReturnResult(sqlmock.NewResult(1, 1)),
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			resp, err := store.Create(ctx, input)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.result, resp)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.mockErr, resp)
		})
	}
}

func TestGet(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()

	query := `
		SELECT code, name, floor, description
		FROM departments
	`

	successResp := []*department.Department{
		{
			Code:        "IT",
			Name:        "IT",
			Floor:       1,
			Description: "desc",
		},
	}

	tests := []struct {
		desc    string
		result  []*department.Department
		mock    func()
		mockErr error
	}{
		{
			desc:    "query error",
			result:  nil,
			mockErr: sql.ErrConnDone,
			mock: func() {
				mock.ExpectQuery(query).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			desc:    "scan error",
			result:  nil,
			mockErr: fmt.Errorf("%w", errors.Error(`sql: Scan error on column index 2, name \"floor\": converting driver.Value type string (\"INVALID_INT\") to a int: invalid syntax`)),
			mock: func() {
				_ = sqlmock.NewRows(
					[]string{"code", "name", "floor", "description"},
				).AddRow("IT", "IT", "INVALID_INT", "desc")

				mock.ExpectQuery(query).
					WillReturnError(fmt.Errorf("%w", errors.Error(`sql: Scan error on column index 2, name \"floor\": converting driver.Value type string (\"INVALID_INT\") to a int: invalid syntax`)))
			},
		},
		{
			desc:    "rows error",
			result:  nil,
			mockErr: sql.ErrConnDone,
			mock: func() {
				rows := sqlmock.NewRows(
					[]string{"code", "name", "floor", "description"},
				).
					AddRow("IT", "IT", 1, "desc").
					RowError(0, sql.ErrConnDone)

				mock.ExpectQuery(query).
					WillReturnRows(rows)
			},
		},
		{
			desc:    "success",
			result:  successResp,
			mockErr: nil,
			mock: func() {
				rows := sqlmock.NewRows(
					[]string{"code", "name", "floor", "description"},
				).AddRow("IT", "IT", 1, "desc")

				mock.ExpectQuery(query).
					WillReturnRows(rows)
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mock()

			resp, err := store.Get(ctx)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.result, resp)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.mockErr, resp)
		})
	}
}

func TestGetByCode(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()

	query := `
		SELECT code, name, floor, description
		FROM departments
		WHERE code = ?
	`

	successResp := &department.Department{
		Code:        "IT",
		Name:        "IT",
		Floor:       1,
		Description: "d",
	}

	tests := []struct {
		desc    string
		code    string
		result  *department.Department
		mockErr error
		mock    func()
	}{
		{
			desc:    "success",
			code:    "IT",
			result:  successResp,
			mockErr: nil,
			mock: func() {
				rows := sqlmock.NewRows(
					[]string{"code", "name", "floor", "description"},
				).AddRow("IT", "IT", 1, "d")

				mock.ExpectQuery(query).
					WithArgs("IT").
					WillReturnRows(rows)
			},
		},
		{
			desc:    "not found",
			code:    "IT",
			result:  nil,
			mockErr: sql.ErrNoRows,
			mock: func() {
				mock.ExpectQuery(query).
					WithArgs("IT").
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mock()

			resp, err := store.GetByCode(ctx, tt.code)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.result, resp)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.mockErr, resp)
		})
	}
}

func TestUpdate(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()

	query := `
		UPDATE departments
		SET name = ?, floor = ?, description = ?
		WHERE code = ?
	`

	input := &department.NewDepartment{
		Name:        "HR",
		Floor:       2,
		Description: "updated",
	}

	successResp := &department.Department{
		Code:        "IT",
		Name:        "HR",
		Floor:       2,
		Description: "updated",
	}

	tests := []struct {
		desc    string
		result  *department.Department
		mock    *sqlmock.ExpectedExec
		mockErr error
	}{
		{
			desc:    "exec error",
			result:  nil,
			mockErr: sql.ErrConnDone,
			mock: mock.ExpectExec(query).
				WithArgs(input.Name, input.Floor, input.Description, "IT").
				WillReturnError(sql.ErrConnDone),
		},
		{
			desc:    "no rows affected",
			result:  nil,
			mockErr: sql.ErrNoRows,
			mock: mock.ExpectExec(query).
				WithArgs(input.Name, input.Floor, input.Description, "IT").
				WillReturnResult(sqlmock.NewResult(0, 0)),
		},
		{
			desc:    "success",
			result:  successResp,
			mockErr: nil,
			mock: mock.ExpectExec(query).
				WithArgs(input.Name, input.Floor, input.Description, "IT").
				WillReturnResult(sqlmock.NewResult(1, 1)),
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			resp, err := store.Update(ctx, "IT", input)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.result, resp)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.mockErr, resp)
		})
	}
}

func TestDelete(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()

	checkQuery := `
		SELECT COUNT(1)
		FROM employees
		WHERE department = ?
	`

	deleteQuery := `
		DELETE FROM departments
		WHERE code = ?
	`

	tests := []struct {
		desc    string
		result  string
		mockErr error
		mock    func()
	}{
		{
			desc:    "count query error",
			result:  "",
			mockErr: sql.ErrConnDone,
			mock: func() {
				mock.ExpectQuery(checkQuery).
					WithArgs("IT").
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			desc:    "employees exist",
			result:  "",
			mockErr: errors.Error("department has employees mapped"),
			mock: func() {
				_ = sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(checkQuery).
					WithArgs("IT").
					WillReturnError(errors.Error("department has employees mapped"))
			},
		},
		{
			desc:    "delete exec error",
			result:  "",
			mockErr: sql.ErrConnDone,
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(checkQuery).
					WithArgs("IT").
					WillReturnRows(rows)

				mock.ExpectExec(deleteQuery).
					WithArgs("IT").
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			desc:    "no rows deleted",
			result:  "",
			mockErr: sql.ErrNoRows,
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(checkQuery).
					WithArgs("IT").
					WillReturnRows(rows)

				mock.ExpectExec(deleteQuery).
					WithArgs("IT").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			desc:    "success",
			result:  "Department deleted successfully",
			mockErr: nil,
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(checkQuery).
					WithArgs("IT").
					WillReturnRows(rows)

				mock.ExpectExec(deleteQuery).
					WithArgs("IT").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mock()

			resp, err := store.Delete(ctx, "IT")

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.result, resp)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.mockErr, resp)
		})
	}
}

func TestExistsByName(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()

	tests := []struct {
		desc        string
		name        string
		excludeCode *string
		result      bool
		mockErr     error
		mock        func()
	}{
		{
			desc:    "exists without exclude",
			name:    "IT",
			result:  true,
			mockErr: nil,
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(
					`SELECT COUNT(1) FROM departments WHERE name = ?`,
				).
					WithArgs("IT").
					WillReturnRows(rows)
			},
		},
		{
			desc:        "exists with exclude",
			name:        "IT",
			excludeCode: strPtr("IT"),
			result:      true,
			mockErr:     nil,
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(
					`SELECT COUNT(1) FROM departments WHERE name = ? AND code != ?`,
				).
					WithArgs("IT", "IT").
					WillReturnRows(rows)
			},
		},
		{
			desc:    "query error",
			name:    "IT",
			result:  false,
			mockErr: sql.ErrConnDone,
			mock: func() {
				mock.ExpectQuery(
					`SELECT COUNT(1) FROM departments WHERE name = ? AND code != ?`,
				).
					WithArgs("IT").
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mock()

			resp, err := store.ExistsByName(ctx, tt.name, tt.excludeCode)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, resp, tt.result)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, err, tt.mockErr)
		})
	}
}

func strPtr(s string) *string {
	return &s
}
