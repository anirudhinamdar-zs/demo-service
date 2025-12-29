package employee

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"demo-service/models/employee"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	store := Init(db)

	query := `
		INSERT INTO employees
			(name, email, phone_number, dob, major, city, department)
		VALUES
			(?, ?, ?, ?, ?, ?, ?)
	`

	dob := time.Now().String()

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
			desc:      "last insert id error",
			result:    sqlmock.NewErrorResult(sql.ErrConnDone),
			expectErr: true,
		},
		{
			desc:   "success",
			result: sqlmock.NewResult(1, 1),
		},
	}

	input := &employee.NewEmployee{
		Name:        "John",
		Email:       "john@test.com",
		PhoneNumber: "999",
		DOB:         dob,
		Major:       "CS",
		City:        "NY",
		Department:  "IT",
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			exp := mock.ExpectExec(query).
				WithArgs(
					input.Name,
					input.Email,
					input.PhoneNumber,
					input.DOB,
					input.Major,
					input.City,
					input.Department,
				)

			if tt.mockErr != nil {
				exp.WillReturnError(tt.mockErr)
			} else {
				exp.WillReturnResult(tt.result)
			}

			_, err := store.Create(ctx, input)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	store := Init(db)

	baseQuery := `
		SELECT
			id, name, email, phone_number, dob, major, city, department
		FROM employees
		WHERE 1 = 1
	`

	tests := []struct {
		desc      string
		filter    employee.Filter
		mockSetup func()
		expectErr bool
	}{
		{
			desc: "query error",
			mockSetup: func() {
				mock.ExpectQuery(baseQuery).
					WillReturnError(sql.ErrConnDone)
			},
			expectErr: true,
		},
		{
			desc: "scan error",
			mockSetup: func() {
				rows := sqlmock.NewRows(
					[]string{"id", "name", "email", "phone_number", "dob", "major", "city", "department"},
				).AddRow("BAD", "x", "x", "x", time.Now(), "x", "x", "x")

				mock.ExpectQuery(baseQuery).
					WillReturnRows(rows)
			},
			expectErr: true,
		},
		{
			desc: "rows error",
			mockSetup: func() {
				rows := sqlmock.NewRows(
					[]string{"id", "name", "email", "phone_number", "dob", "major", "city", "department"},
				).AddRow(1, "x", "x", "x", time.Now(), "x", "x", "x").
					RowError(0, sql.ErrConnDone)

				mock.ExpectQuery(baseQuery).
					WillReturnRows(rows)
			},
			expectErr: true,
		},
		{
			desc: "success",
			mockSetup: func() {
				rows := sqlmock.NewRows(
					[]string{"id", "name", "email", "phone_number", "dob", "major", "city", "department"},
				).AddRow(1, "x", "x", "x", time.Now(), "x", "x", "IT")

				mock.ExpectQuery(baseQuery).
					WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mockSetup()
			_, err := store.Get(ctx, employee.Filter{})
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetById(t *testing.T) {
	ctx := context.Background()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	store := Init(db)

	query := `
		SELECT id, name, email, phone_number, dob, major, city, department
		FROM employees
		WHERE id = ?
	`

	tests := []struct {
		desc      string
		mockSetup func()
		expectErr bool
	}{
		{
			desc: "query error",
			mockSetup: func() {
				mock.ExpectQuery(query).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			expectErr: true,
		},
		{
			desc: "success",
			mockSetup: func() {
				rows := sqlmock.NewRows(
					[]string{"id", "name", "email", "phone_number", "dob", "major", "city", "department"},
				).AddRow(1, "x", "x", "x", time.Now(), "x", "x", "IT")

				mock.ExpectQuery(query).
					WithArgs(1).
					WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			tt.mockSetup()
			_, err := store.GetById(ctx, 1)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	store := Init(db)

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
				WithArgs(
					input.Name,
					input.Email,
					input.PhoneNumber,
					input.DOB,
					input.Major,
					input.City,
					input.Department,
					1,
				)

			if tt.mockErr != nil {
				exp.WillReturnError(tt.mockErr)
			} else {
				exp.WillReturnResult(tt.result)
			}

			_, err := store.Update(ctx, 1, input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	store := Init(db)

	query := `DELETE FROM employees WHERE id = ?`

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
			desc:      "no rows deleted",
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
			exp := mock.ExpectExec(query).WithArgs(1)
			if tt.mockErr != nil {
				exp.WillReturnError(tt.mockErr)
			} else {
				exp.WillReturnResult(tt.result)
			}

			_, err := store.Delete(ctx, 1)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExistsByEmail(t *testing.T) {
	ctx := context.Background()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	store := Init(db)

	tests := []struct {
		desc      string
		email     string
		excludeID *int
		count     int
		expect    bool
		expectErr bool
	}{
		{
			desc:   "exists without exclude",
			email:  "a@a.com",
			count:  1,
			expect: true,
		},
		{
			desc:      "exists with exclude",
			email:     "a@a.com",
			excludeID: intPtr(1),
			count:     1,
			expect:    true,
		},
		{
			desc:      "query error",
			email:     "a@a.com",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			query := `SELECT COUNT(1) FROM employees WHERE email = ?`
			args := []driver.Value{tt.email}

			if tt.excludeID != nil {
				query += ` AND id != ?`
				args = append(args, *tt.excludeID)
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

			res, err := store.ExistsByEmail(ctx, tt.email, tt.excludeID)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, res)
			}
		})
	}
}

func TestCountByDepartment(t *testing.T) {
	ctx := context.Background()
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	store := Init(db)

	query := `
		SELECT COUNT(1)
		FROM employees
		WHERE department = ?
	`

	tests := []struct {
		desc      string
		count     int
		expectErr bool
	}{
		{
			desc:      "query error",
			expectErr: true,
		},
		{
			desc:  "success",
			count: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.expectErr {
				mock.ExpectQuery(query).
					WithArgs("IT").
					WillReturnError(sql.ErrConnDone)
			} else {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(tt.count)
				mock.ExpectQuery(query).
					WithArgs("IT").
					WillReturnRows(rows)
			}

			res, err := store.CountByDepartment(ctx, "IT")

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.count, res)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
