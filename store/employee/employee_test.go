package employee

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"demo-service/models/employee"
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

	query := `INSERT INTO employees (name, email, phone_number, dob, major, city, department, deleted_at) VALUES (?, ?, ?, ?, ?, ?, ?, NULL)`

	dob := time.Now().String()

	input := &employee.NewEmployee{
		Name:        "John",
		Email:       "john@test.com",
		PhoneNumber: "999",
		DOB:         dob,
		Major:       "CS",
		City:        "NY",
		Department:  "IT",
	}

	resp := &employee.Employee{
		ID:          1,
		Name:        "John",
		Email:       "john@test.com",
		PhoneNumber: "999",
		DOB:         dob,
		Major:       "CS",
		City:        "NY",
		Department:  "IT",
	}

	tests := []struct {
		desc    string
		result  *employee.Employee
		mock    *sqlmock.ExpectedExec
		mockErr error
	}{
		{
			desc:    "exec error",
			mockErr: sql.ErrConnDone,
			mock: mock.ExpectExec(query).
				WithArgs(
					input.Name,
					input.Email,
					input.PhoneNumber,
					input.DOB,
					input.Major,
					input.City,
					input.Department,
				).WillReturnResult(nil).WillReturnError(sql.ErrConnDone),
		},
		{
			desc:    "last insert id error",
			mockErr: sql.ErrConnDone,
			mock: mock.ExpectExec(query).
				WithArgs(
					input.Name,
					input.Email,
					input.PhoneNumber,
					input.DOB,
					input.Major,
					input.City,
					input.Department,
				).WillReturnResult(nil).WillReturnError(sql.ErrConnDone),
		},
		{
			desc:    "success",
			result:  resp,
			mockErr: nil,
			mock: mock.ExpectExec(query).
				WithArgs(
					input.Name,
					input.Email,
					input.PhoneNumber,
					input.DOB,
					input.Major,
					input.City,
					input.Department,
				).WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(nil),
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

	baseQuery := `
		SELECT
			id,
			name,
			email,
			phone_number,
			dob,
			major,
			city,
			department,
			deleted_at
		FROM employees
		WHERE deleted_at IS NULL
	`

	now := time.Now().String()

	successResp := []*employee.Employee{
		{
			ID:          1,
			Name:        "x",
			Email:       "x",
			PhoneNumber: "x",
			DOB:         now,
			Major:       "x",
			City:        "x",
			Department:  "IT",
			DeletedAt:   nil,
		},
	}

	tests := []struct {
		desc    string
		filter  employee.Filter
		mock    func()
		result  []*employee.Employee
		mockErr error
	}{
		{
			desc: "query error",
			mock: func() {
				mock.ExpectQuery(baseQuery).
					WillReturnError(sql.ErrConnDone)
			},
			mockErr: sql.ErrConnDone,
		},
		{
			desc: "scan error",
			mock: func() {
				_ = sqlmock.NewRows(
					[]string{"id", "name", "email", "phone_number", "dob", "major", "city", "department", "deleted_at"},
				).AddRow("BAD", "x", "x", "x", now, "x", "x", "x", "x")

				mock.ExpectQuery(baseQuery).
					WillReturnError(fmt.Errorf("%w", errors.Error(`sql: Scan error on column index 0, name "id": converting driver.Value type string ("BAD") to a int: invalid syntax`)))
			},
			mockErr: fmt.Errorf("%w", errors.Error(`sql: Scan error on column index 0, name "id": converting driver.Value type string ("BAD") to a int: invalid syntax`)),
		},
		{
			desc: "rows error",
			mock: func() {
				rows := sqlmock.NewRows(
					[]string{"id", "name", "email", "phone_number", "dob", "major", "city", "department"},
				).
					AddRow(1, "x", "x", "x", now, "x", "x", "IT").
					RowError(0, sql.ErrConnDone)

				mock.ExpectQuery(baseQuery).
					WillReturnRows(rows)
			},
			mockErr: sql.ErrConnDone,
		},
		{
			desc: "success",
			mock: func() {
				rows := sqlmock.NewRows(
					[]string{"id", "name", "email", "phone_number", "dob", "major", "city", "department", "deleted_at"},
				).
					AddRow(1, "x", "x", "x", now, "x", "x", "IT", nil)

				mock.ExpectQuery(baseQuery).
					WillReturnRows(rows)
			},
			result:  successResp,
			mockErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mock()

			resp, err := store.Get(ctx, tt.filter)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, resp, tt.result)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, err, tt.mockErr)
		})
	}
}

func TestGetById(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()

	query := `
		SELECT id, name, email, phone_number, dob, major, city, department, deleted_at
		FROM employees
		WHERE id = ? AND deleted_at IS NULL
	`

	successResp := &employee.Employee{
		ID:          1,
		Name:        "x",
		Email:       "x",
		PhoneNumber: "x",
		DOB:         time.Now().String(),
		Major:       "x",
		City:        "x",
		Department:  "IT",
		DeletedAt:   nil,
	}

	tests := []struct {
		desc    string
		result  *employee.Employee
		mock    func()
		mockErr error
	}{
		{
			desc:    "query error",
			mockErr: sql.ErrNoRows,
			mock: func() {
				mock.ExpectQuery(query).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			desc:    "success",
			result:  successResp,
			mockErr: nil,
			mock: func() {
				rows := sqlmock.NewRows(
					[]string{"id", "name", "email", "phone_number", "dob", "major", "city", "department", "deleted_at"},
				).AddRow(
					1,
					"x",
					"x",
					"x",
					successResp.DOB,
					"x",
					"x",
					"IT",
					nil,
				)

				mock.ExpectQuery(query).
					WithArgs(1).
					WillReturnRows(rows)
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mock()

			resp, err := store.GetById(ctx, 1)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, resp, tt.result)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, err, tt.mockErr)
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
		UPDATE employees
		SET
			name = ?,
			email = ?,
			phone_number = ?,
			dob = ?,
			major = ?,
			city = ?,
			department = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	getByIDQuery := `
		SELECT id, name, email, phone_number, dob, major, city, department, deleted_at
		FROM employees
		WHERE id = ? AND deleted_at IS NULL
	`

	dob := time.Now().String()

	input := &employee.NewEmployee{
		Name:        "A",
		Email:       "a@a.com",
		PhoneNumber: "1",
		DOB:         dob,
		Major:       "CS",
		City:        "NY",
		Department:  "IT",
	}

	successResp := &employee.Employee{
		ID:          1,
		Name:        "A",
		Email:       "a@a.com",
		PhoneNumber: "1",
		DOB:         dob,
		Major:       "CS",
		City:        "NY",
		Department:  "IT",
	}

	tests := []struct {
		desc    string
		result  *employee.Employee
		mock    func()
		mockErr error
	}{
		{
			desc:    "exec error",
			result:  nil,
			mockErr: sql.ErrConnDone,
			mock: func() {
				mock.ExpectExec(query).
					WithArgs(
						input.Name,
						input.Email,
						input.PhoneNumber,
						input.DOB,
						input.Major,
						input.City,
						input.Department,
						1,
					).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			desc:    "no rows affected",
			result:  nil,
			mockErr: sql.ErrNoRows,
			mock: func() {
				mock.ExpectExec(query).
					WithArgs(
						input.Name,
						input.Email,
						input.PhoneNumber,
						input.DOB,
						input.Major,
						input.City,
						input.Department,
						1,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			desc:    "success",
			result:  successResp,
			mockErr: nil,
			mock: func() {
				mock.ExpectExec(query).
					WithArgs(
						input.Name,
						input.Email,
						input.PhoneNumber,
						input.DOB,
						input.Major,
						input.City,
						input.Department,
						1,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				rows := sqlmock.NewRows(
					[]string{"id", "name", "email", "phone_number", "dob", "major", "city", "department", "deleted_at"},
				).AddRow(
					1,
					input.Name,
					input.Email,
					input.PhoneNumber,
					input.DOB,
					input.Major,
					input.City,
					input.Department,
					nil,
				)

				mock.ExpectQuery(getByIDQuery).
					WithArgs(1).
					WillReturnRows(rows)
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mock()

			resp, err := store.Update(ctx, 1, input)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, resp, tt.result)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, err, tt.mockErr)
		})
	}
}

func TestDelete(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()

	query := `UPDATE employees SET deleted_at = CURRENT_DATE WHERE id = ? AND deleted_at IS NULL`

	tests := []struct {
		desc    string
		result  string
		mock    func()
		mockErr error
	}{
		{
			desc:    "exec error",
			result:  "",
			mockErr: sql.ErrConnDone,
			mock: func() {
				mock.ExpectExec(query).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			desc:    "no rows deleted",
			result:  "",
			mockErr: sql.ErrNoRows,
			mock: func() {
				mock.ExpectExec(query).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			desc:    "success",
			result:  "Employee deleted successfully",
			mockErr: nil,
			mock: func() {
				mock.ExpectExec(query).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mock()

			resp, err := store.Delete(ctx, 1)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, resp, tt.result)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, err, tt.mockErr)
		})
	}
}

func TestExistsByEmail(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()

	tests := []struct {
		desc      string
		email     string
		excludeID *int
		result    bool
		mockErr   error
		mock      func()
	}{
		{
			desc:    "exists without exclude",
			email:   "a@a.com",
			result:  true,
			mockErr: nil,
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(
					`SELECT COUNT(1) FROM employees WHERE email = ? AND deleted_at IS NULL`,
				).
					WithArgs("a@a.com").
					WillReturnRows(rows)
			},
		},
		{
			desc:      "exists with exclude",
			email:     "a@a.com",
			excludeID: intPtr(1),
			result:    true,
			mockErr:   nil,
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(
					`SELECT COUNT(1) FROM employees WHERE email = ? AND deleted_at IS NULL AND id != ?`,
				).
					WithArgs("a@a.com", 1).
					WillReturnRows(rows)
			},
		},
		{
			desc:    "query error",
			email:   "a@a.com",
			result:  false,
			mockErr: sql.ErrConnDone,
			mock: func() {
				mock.ExpectQuery(
					`SELECT COUNT(1) FROM employees WHERE email = ? AND deleted_at IS NULL AND id != ?`,
				).
					WithArgs("a@a.com").
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mock()

			resp, err := store.ExistsByEmail(ctx, tt.email, tt.excludeID)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, resp, tt.result)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, err, tt.mockErr)
		})
	}
}

func TestCountByDepartment(t *testing.T) {
	_, ctx, mock, err := Initialize(t)
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	store := Init()

	query := `
		SELECT COUNT(1)
		FROM employees
		WHERE department = ? AND deleted_at IS NULL
	`

	tests := []struct {
		desc    string
		result  int
		mockErr error
		mock    func()
	}{
		{
			desc:    "query error",
			result:  0,
			mockErr: sql.ErrConnDone,
			mock: func() {
				mock.ExpectQuery(query).
					WithArgs("IT").
					WillReturnError(sql.ErrConnDone)
			},
		},
		{
			desc:    "success",
			result:  3,
			mockErr: nil,
			mock: func() {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(3)
				mock.ExpectQuery(query).
					WithArgs("IT").
					WillReturnRows(rows)
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mock()

			resp, err := store.CountByDepartment(ctx, "IT")

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, resp, tt.result)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, err, tt.mockErr)
		})
	}
}

func intPtr(i int) *int {
	return &i
}
