package employee

import (
	"testing"

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
		desc        string
		input       *employee.NewEmployee
		setupMocks  func(emp *store.MockEmployee, dep *store.MockDepartment)
		expectError bool
	}{
		{
			desc: "invalid department",
			input: &employee.NewEmployee{
				Department: "BIO",
			},
			expectError: true,
		},
		{
			desc: "department does not exist",
			input: &employee.NewEmployee{
				Department: "IT",
				Email:      "a@a.com",
			},
			setupMocks: func(emp *store.MockEmployee, dep *store.MockDepartment) {
				dep.EXPECT().
					GetByCode(gomock.Any(), "IT").
					Return(nil, errors.EntityNotFound{Entity: "IT"})
			},
			expectError: true,
		},
		{
			desc: "email already exists",
			input: &employee.NewEmployee{
				Department: "IT",
				Email:      "a@a.com",
			},
			setupMocks: func(emp *store.MockEmployee, dep *store.MockDepartment) {
				dep.EXPECT().
					GetByCode(gomock.Any(), "IT").
					Return(&department.Department{}, nil)

				emp.EXPECT().
					ExistsByEmail(gomock.Any(), "a@a.com", nil).
					Return(true, nil)
			},
			expectError: true,
		},
		{
			desc: "success",
			input: &employee.NewEmployee{
				Department: "IT",
				Email:      "a@a.com",
			},
			setupMocks: func(emp *store.MockEmployee, dep *store.MockDepartment) {
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockEmp := store.NewMockEmployee(ctrl)
			mockDep := store.NewMockDepartment(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockEmp, mockDep)
			}

			svc := New(mockEmp, mockDep)

			_, err := svc.Create(ctx, tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
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
		desc        string
		filter      employee.Filter
		setupMocks  func(emp *store.MockEmployee)
		expectError bool
	}{
		{
			desc: "invalid department filter",
			filter: employee.Filter{
				Department: &invalidDept,
			},
			expectError: true,
		},
		{
			desc:   "success without filters",
			filter: employee.Filter{},
			setupMocks: func(emp *store.MockEmployee) {
				emp.EXPECT().
					Get(gomock.Any(), employee.Filter{}).
					Return([]*employee.Employee{}, nil)
			},
			expectError: false,
		},
		{
			desc: "success with department filter",
			filter: employee.Filter{
				Department: &validDept,
			},
			setupMocks: func(emp *store.MockEmployee) {
				emp.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return([]*employee.Employee{
						{ID: 1, Department: "IT"},
					}, nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockEmp := store.NewMockEmployee(ctrl)
			mockDep := store.NewMockDepartment(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockEmp)
			}

			svc := New(mockEmp, mockDep)

			_, err := svc.Get(ctx, tt.filter)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmployeeService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	tests := []struct {
		desc        string
		id          int
		setupMocks  func(emp *store.MockEmployee)
		expectError bool
	}{
		{
			desc: "success",
			id:   1,
			setupMocks: func(emp *store.MockEmployee) {
				emp.EXPECT().
					GetById(gomock.Any(), 1).
					Return(&employee.Employee{ID: 1}, nil)
			},
		},
		{
			desc: "not found",
			id:   99,
			setupMocks: func(emp *store.MockEmployee) {
				emp.EXPECT().
					GetById(gomock.Any(), 99).
					Return(nil, errors.EntityNotFound{})
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockEmp := store.NewMockEmployee(ctrl)
			mockDep := store.NewMockDepartment(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockEmp)
			}

			svc := New(mockEmp, mockDep)

			_, err := svc.GetById(ctx, tt.id)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
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
		desc        string
		input       *employee.NewEmployee
		setupMocks  func(emp *store.MockEmployee, dep *store.MockDepartment)
		expectError bool
	}{
		{
			desc: "invalid department",
			input: &employee.NewEmployee{
				Department: invalidDept,
			},
			expectError: true,
		},
		{
			desc: "department does not exist",
			input: &employee.NewEmployee{
				Department: validDept,
				Email:      "a@a.com",
			},
			setupMocks: func(emp *store.MockEmployee, dep *store.MockDepartment) {
				dep.EXPECT().
					GetByCode(gomock.Any(), validDept).
					Return(nil, errors.EntityNotFound{})
			},
			expectError: true,
		},
		{
			desc: "email already exists",
			input: &employee.NewEmployee{
				Department: validDept,
				Email:      "a@a.com",
			},
			setupMocks: func(emp *store.MockEmployee, dep *store.MockDepartment) {
				dep.EXPECT().
					GetByCode(gomock.Any(), validDept).
					Return(&department.Department{}, nil)

				emp.EXPECT().
					ExistsByEmail(gomock.Any(), "a@a.com", &id).
					Return(true, nil)
			},
			expectError: true,
		},
		{
			desc: "success",
			input: &employee.NewEmployee{
				Department: validDept,
				Email:      "a@a.com",
			},
			setupMocks: func(emp *store.MockEmployee, dep *store.MockDepartment) {
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockEmp := store.NewMockEmployee(ctrl)
			mockDep := store.NewMockDepartment(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockEmp, mockDep)
			}

			svc := New(mockEmp, mockDep)

			_, err := svc.Update(ctx, id, tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmployeeService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	tests := []struct {
		desc        string
		id          int
		setupMocks  func(emp *store.MockEmployee)
		expectError bool
	}{
		{
			desc: "success",
			id:   1,
			setupMocks: func(emp *store.MockEmployee) {
				emp.EXPECT().
					Delete(gomock.Any(), 1).
					Return("Employee deleted successfully", nil)
			},
		},
		{
			desc: "not found",
			id:   99,
			setupMocks: func(emp *store.MockEmployee) {
				emp.EXPECT().
					Delete(gomock.Any(), 99).
					Return("", errors.EntityNotFound{})
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockEmp := store.NewMockEmployee(ctrl)
			mockDep := store.NewMockDepartment(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockEmp)
			}

			svc := New(mockEmp, mockDep)

			_, err := svc.Delete(ctx, tt.id)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
