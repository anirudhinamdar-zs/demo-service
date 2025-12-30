package department

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"demo-service/models/department"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/golang/mock/gomock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Initialize(t *testing.T) (*gomock.Controller, *gofr.Context, sqlmock.Sqlmock, error) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	return ctrl, ctx, mock, err
}

func TestCreate(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()
	tests := []struct {
		desc           string
		input          *department.Department
		mockResult     driver.Result
		mockError      error
		expectedOutput *department.Department
		expectedErr    error
	}{
		{
			desc: "success",
			input: &department.Department{
				Code: "IT", Name: "IT", Floor: 1, Description: "desc",
			},
			mockResult:     sqlmock.NewResult(1, 1),
			expectedOutput: &department.Department{Code: "IT", Name: "IT", Floor: 1, Description: "desc"},
		},
		{
			desc:        "no rows affected",
			input:       &department.Department{},
			mockResult:  sqlmock.NewResult(0, 0),
			expectedErr: sql.ErrNoRows,
		},
		{
			desc:        "db error",
			input:       &department.Department{},
			mockError:   errors.New("db error"),
			expectedErr: errors.New("db error"),
		},
	}

	createQuery := `INSERT INTO departments (code, name, floor, description) VALUES (?, ?, ?, ?)`

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mock.ExpectExec(createQuery).
				WillReturnResult(tt.mockResult).
				WillReturnError(tt.mockError)

			res, err := store.Create(ctx, tt.input)

			assert.Equal(t, tt.expectedErr != nil, err != nil)
			assert.Equal(t, tt.expectedOutput, res)
		})
	}
}

func TestGet(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()

	query := `
		SELECT code, name, floor, description
		FROM departments
	`

	tests := []struct {
		desc       string
		mockSetup  func()
		expectErr  bool
		expectSize int
	}{
		{
			desc: "query error",
			mockSetup: func() {
				mock.ExpectQuery(query).
					WillReturnError(sql.ErrConnDone)
			},
			expectErr: true,
		},
		{
			desc: "scan error",
			mockSetup: func() {
				rows := sqlmock.NewRows(
					[]string{"code", "name", "floor", "description"},
				).AddRow("IT", "IT", "INVALID_INT", "desc")

				mock.ExpectQuery(query).
					WillReturnRows(rows)
			},
			expectErr: true,
		},
		{
			desc: "rows error",
			mockSetup: func() {
				rows := sqlmock.NewRows(
					[]string{"code", "name", "floor", "description"},
				).AddRow("IT", "IT", 1, "desc").
					RowError(0, sql.ErrConnDone)

				mock.ExpectQuery(query).
					WillReturnRows(rows)
			},
			expectErr: true,
		},
		{
			desc: "success",
			mockSetup: func() {
				rows := sqlmock.NewRows(
					[]string{"code", "name", "floor", "description"},
				).AddRow("IT", "IT", 1, "desc")

				mock.ExpectQuery(query).
					WillReturnRows(rows)
			},
			expectSize: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mockSetup()

			res, err := store.Get(ctx)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, res, tt.expectSize)
			}
		})
	}
}

func TestGetByCode(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()

	query := `
		SELECT code, name, floor, description
		FROM departments
		WHERE code = ?
	`

	tests := []struct {
		desc        string
		code        string
		rows        *sqlmock.Rows
		mockError   error
		expectError bool
	}{
		{
			desc: "success",
			code: "IT",
			rows: sqlmock.NewRows([]string{"code", "name", "floor", "description"}).
				AddRow("IT", "IT", 1, "d"),
		},
		{
			desc:        "not found",
			code:        "IT",
			mockError:   sql.ErrNoRows,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			exp := mock.ExpectQuery(query).WithArgs(tt.code)
			if tt.mockError != nil {
				exp.WillReturnError(tt.mockError)
			} else {
				exp.WillReturnRows(tt.rows)
			}

			_, err := store.GetByCode(ctx, tt.code)
			assert.Equal(t, tt.expectError, err != nil)
		})
	}
}

func TestUpdate(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
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

	tests := []struct {
		desc      string
		result    driver.Result
		mockErr   error
		expectErr bool
	}{
		{
			desc:      "exec error",
			mockErr:   sql.ErrConnDone,
			expectErr: true,
		},
		{
			desc:      "no rows affected",
			result:    sqlmock.NewResult(0, 0),
			expectErr: true,
		},
		{
			desc:   "success",
			result: sqlmock.NewResult(1, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			exp := mock.ExpectExec(query).
				WithArgs(input.Name, input.Floor, input.Description, "IT")

			if tt.mockErr != nil {
				exp.WillReturnError(tt.mockErr)
			} else {
				exp.WillReturnResult(tt.result)
			}

			_, err := store.Update(ctx, "IT", input)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
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
		desc      string
		mockSetup func()
		expectErr bool
	}{
		{
			desc: "count query error",
			mockSetup: func() {
				mock.ExpectQuery(checkQuery).
					WithArgs("IT").
					WillReturnError(sql.ErrConnDone)
			},
			expectErr: true,
		},
		{
			desc: "employees exist",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(checkQuery).
					WithArgs("IT").
					WillReturnRows(rows)
			},
			expectErr: true,
		},
		{
			desc: "delete exec error",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(checkQuery).
					WithArgs("IT").
					WillReturnRows(rows)

				mock.ExpectExec(deleteQuery).
					WithArgs("IT").
					WillReturnError(sql.ErrConnDone)
			},
			expectErr: true,
		},
		{
			desc: "no rows deleted",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(checkQuery).
					WithArgs("IT").
					WillReturnRows(rows)

				mock.ExpectExec(deleteQuery).
					WithArgs("IT").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectErr: true,
		},
		{
			desc: "success",
			mockSetup: func() {
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

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mockSetup()

			_, err := store.Delete(ctx, "IT")

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExistsByName(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()

	tests := []struct {
		desc        string
		name        string
		excludeCode *string
		count       int
		expect      bool
		expectErr   bool
	}{
		{
			desc:   "exists without exclude",
			name:   "IT",
			count:  1,
			expect: true,
		},
		{
			desc:        "exists with exclude",
			name:        "IT",
			excludeCode: strPtr("IT"),
			count:       1,
			expect:      true,
		},
		{
			desc:      "query error",
			name:      "IT",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			query := `SELECT COUNT(1) FROM departments WHERE name = ?`
			args := []driver.Value{tt.name}

			if tt.excludeCode != nil {
				query += ` AND code != ?`
				args = append(args, *tt.excludeCode)
			}

			if tt.expectErr {
				mock.ExpectQuery(query).
					WithArgs(args...).
					WillReturnError(sql.ErrConnDone)
			} else {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(tt.count)
				mock.ExpectQuery(query).
					WithArgs(args...).
					WillReturnRows(rows)
			}

			res, err := store.ExistsByName(ctx, tt.name, tt.excludeCode)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, res)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
