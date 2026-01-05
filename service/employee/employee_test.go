package employee

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"demo-service/models/department"
	"demo-service/models/employee"
	"demo-service/store"
)

func TestEmployeeService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	tests := []struct {
		desc    string
		input   *employee.NewEmployee
		mock    func(emp *store.MockEmployee, dep *store.MockDepartment)
		result  *employee.Employee
		mockErr error
	}{
		{
			desc: "invalid department",
			input: &employee.NewEmployee{
				Department: "BIO",
			},
			mock: func(emp *store.MockEmployee, dep *store.MockDepartment) {
				// no calls expected
			},
			result:  nil,
			mockErr: errors.InvalidParam{Param: []string{"department"}},
		},
		{
			desc: "department does not exist",
			input: &employee.NewEmployee{
				Department: "IT",
				Email:      "a@a.com",
			},
			mock: func(emp *store.MockEmployee, dep *store.MockDepartment) {
				dep.EXPECT().
					GetByCode(gomock.Any(), "IT").
					Return(nil, errors.EntityNotFound{Entity: "IT"})
			},
			result:  nil,
			mockErr: errors.EntityNotFound{Entity: "IT"},
		},
		{
			desc: "email already exists",
			input: &employee.NewEmployee{
				Department: "IT",
				Email:      "a@a.com",
			},
			mock: func(emp *store.MockEmployee, dep *store.MockDepartment) {
				dep.EXPECT().
					GetByCode(gomock.Any(), "IT").
					Return(&department.Department{}, nil)

				emp.EXPECT().
					ExistsByEmail(gomock.Any(), "a@a.com", nil).
					Return(true, nil)
			},
			result:  nil,
			mockErr: errors.EntityAlreadyExists{},
		},
		{
			desc: "success",
			input: &employee.NewEmployee{
				Department: "IT",
				Email:      "a@a.com",
			},
			mock: func(emp *store.MockEmployee, dep *store.MockDepartment) {
				dep.EXPECT().
					GetByCode(gomock.Any(), "IT").
					Return(&department.Department{}, nil)

				emp.EXPECT().
					ExistsByEmail(gomock.Any(), "a@a.com", nil).
					Return(false, nil)

				emp.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&employee.Employee{ID: 1}, nil)
			},
			result:  &employee.Employee{ID: 1},
			mockErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockEmp := store.NewMockEmployee(ctrl)
			mockDep := store.NewMockDepartment(ctrl)

			tt.mock(mockEmp, mockDep)

			svc := New(mockEmp, mockDep)

			resp, err := svc.Create(ctx, tt.input)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.result, resp)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.mockErr, resp)
		})
	}
}

func TestEmployeeService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	validDept := "IT"
	invalidDept := "BIO"

	tests := []struct {
		desc    string
		filter  employee.Filter
		mock    func(emp *store.MockEmployee)
		result  []*employee.Employee
		mockErr error
	}{
		{
			desc: "invalid department filter",
			filter: employee.Filter{
				Department: &invalidDept,
			},
			mock: func(emp *store.MockEmployee) {
				// no store call expected
			},
			result:  nil,
			mockErr: errors.EntityNotFound{Entity: "BIO"},
		},
		{
			desc:   "success without filters",
			filter: employee.Filter{},
			mock: func(emp *store.MockEmployee) {
				emp.EXPECT().
					Get(gomock.Any(), employee.Filter{}).
					Return([]*employee.Employee{}, nil)
			},
			result:  []*employee.Employee{},
			mockErr: nil,
		},
		{
			desc: "success with department filter",
			filter: employee.Filter{
				Department: &validDept,
			},
			mock: func(emp *store.MockEmployee) {
				emp.EXPECT().
					Get(gomock.Any(), employee.Filter{
						Department: &validDept,
					}).
					Return([]*employee.Employee{
						{ID: 1, Department: "IT"},
					}, nil)
			},
			result: []*employee.Employee{
				{ID: 1, Department: "IT"},
			},
			mockErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockEmp := store.NewMockEmployee(ctrl)
			mockDep := store.NewMockDepartment(ctrl)

			tt.mock(mockEmp)

			svc := New(mockEmp, mockDep)

			resp, err := svc.Get(ctx, tt.filter)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.result, resp)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.mockErr, resp)
		})
	}
}

func TestEmployeeService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	tests := []struct {
		desc    string
		id      int
		mock    func(emp *store.MockEmployee)
		result  *employee.Employee
		mockErr error
	}{
		{
			desc: "success",
			id:   1,
			mock: func(emp *store.MockEmployee) {
				emp.EXPECT().
					GetById(gomock.Any(), 1).
					Return(&employee.Employee{ID: 1}, nil)
			},
			result:  &employee.Employee{ID: 1},
			mockErr: nil,
		},
		{
			desc: "not found",
			id:   99,
			mock: func(emp *store.MockEmployee) {
				emp.EXPECT().
					GetById(gomock.Any(), 99).
					Return(nil, errors.EntityNotFound{})
			},
			result:  nil,
			mockErr: errors.EntityNotFound{},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockEmp := store.NewMockEmployee(ctrl)
			mockDep := store.NewMockDepartment(ctrl)

			tt.mock(mockEmp)

			svc := New(mockEmp, mockDep)

			resp, err := svc.GetById(ctx, tt.id)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.result, resp)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.mockErr, resp)
		})
	}
}

func TestEmployeeService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	validDept := "IT"
	invalidDept := "BIO"
	id := 1

	tests := []struct {
		desc    string
		input   *employee.NewEmployee
		mock    func(emp *store.MockEmployee, dep *store.MockDepartment)
		result  *employee.Employee
		mockErr error
	}{
		{
			desc: "invalid department",
			input: &employee.NewEmployee{
				Department: invalidDept,
			},
			mock: func(emp *store.MockEmployee, dep *store.MockDepartment) {
				emp.EXPECT().
					GetById(gomock.Any(), id).
					Return(&employee.Employee{ID: id}, nil)
			},
			result:  nil,
			mockErr: errors.EntityNotFound{Entity: invalidDept},
		},
		{
			desc: "department does not exist",
			input: &employee.NewEmployee{
				Department: validDept,
				Email:      "a@a.com",
			},
			mock: func(emp *store.MockEmployee, dep *store.MockDepartment) {
				emp.EXPECT().
					GetById(gomock.Any(), id).
					Return(&employee.Employee{ID: id}, nil)

				dep.EXPECT().
					GetByCode(gomock.Any(), validDept).
					Return(nil, errors.EntityNotFound{Entity: validDept})
			},
			result:  nil,
			mockErr: errors.EntityNotFound{Entity: validDept},
		},
		{
			desc: "email already exists",
			input: &employee.NewEmployee{
				Department: validDept,
				Email:      "a@a.com",
			},
			mock: func(emp *store.MockEmployee, dep *store.MockDepartment) {
				emp.EXPECT().
					GetById(gomock.Any(), id).
					Return(&employee.Employee{ID: id}, nil)

				dep.EXPECT().
					GetByCode(gomock.Any(), validDept).
					Return(&department.Department{}, nil)

				emp.EXPECT().
					ExistsByEmail(gomock.Any(), "a@a.com", &id).
					Return(true, nil)
			},
			result:  nil,
			mockErr: errors.EntityAlreadyExists{},
		},
		{
			desc: "success",
			input: &employee.NewEmployee{
				Department: validDept,
				Email:      "a@a.com",
			},
			mock: func(emp *store.MockEmployee, dep *store.MockDepartment) {
				emp.EXPECT().
					GetById(gomock.Any(), id).
					Return(&employee.Employee{ID: id}, nil)

				dep.EXPECT().
					GetByCode(gomock.Any(), validDept).
					Return(&department.Department{}, nil)

				emp.EXPECT().
					ExistsByEmail(gomock.Any(), "a@a.com", &id).
					Return(false, nil)

				emp.EXPECT().
					Update(gomock.Any(), id, gomock.Any()).
					Return(&employee.Employee{ID: id}, nil)
			},
			result:  &employee.Employee{ID: id},
			mockErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockEmp := store.NewMockEmployee(ctrl)
			mockDep := store.NewMockDepartment(ctrl)

			tt.mock(mockEmp, mockDep)

			svc := New(mockEmp, mockDep)

			resp, err := svc.Update(ctx, id, tt.input)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.result, resp)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.mockErr, resp)
		})
	}
}

func TestEmployeeService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	deletedDate := date.Date{}

	tests := []struct {
		desc    string
		id      int
		mock    func(emp *store.MockEmployee)
		result  string
		mockErr error
	}{
		{
			desc: "success",
			id:   1,
			mock: func(emp *store.MockEmployee) {
				emp.EXPECT().
					GetById(gomock.Any(), 1).
					Return(&employee.Employee{ID: 1}, nil)

				emp.EXPECT().
					Delete(gomock.Any(), 1).
					Return("Employee deleted successfully", nil)
			},
			result:  "Employee deleted successfully",
			mockErr: nil,
		},
		{
			desc: "already deleted",
			id:   1,
			mock: func(emp *store.MockEmployee) {
				emp.EXPECT().
					GetById(gomock.Any(), 1).
					Return(&employee.Employee{
						ID:        1,
						DeletedAt: &deletedDate,
					}, nil)
			},
			result:  "",
			mockErr: errors.EntityNotFound{Entity: "employee"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockEmp := store.NewMockEmployee(ctrl)
			mockDep := store.NewMockDepartment(ctrl)

			tt.mock(mockEmp)

			svc := New(mockEmp, mockDep)

			resp, err := svc.Delete(ctx, tt.id)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.result, resp)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.mockErr, resp)
		})
	}
}
