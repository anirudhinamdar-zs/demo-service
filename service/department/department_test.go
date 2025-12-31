package department

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"demo-service/models/department"
	"demo-service/store"
)

func TestDepartmentService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	tests := []struct {
		desc        string
		input       *department.Department
		setupMocks  func(dep *store.MockDepartment)
		expectError bool
	}{
		{
			desc: "invalid department code",
			input: &department.Department{
				Code: "BIO",
				Name: "Biology",
			},
			expectError: true,
		},
		{
			desc: "department already exists",
			input: &department.Department{
				Code: "IT",
				Name: "Information Technology",
			},
			setupMocks: func(dep *store.MockDepartment) {
				dep.EXPECT().
					ExistsByName(gomock.Any(), "Information Technology", nil).
					Return(true, nil)
			},
			expectError: true,
		},
		{
			desc: "exists check error",
			input: &department.Department{
				Code: "IT",
				Name: "Information Technology",
			},
			setupMocks: func(dep *store.MockDepartment) {
				dep.EXPECT().
					ExistsByName(gomock.Any(), "Information Technology", nil).
					Return(false, errors.DB{Err: errors.DB{}})
			},
			expectError: true,
		},
		{
			desc: "success",
			input: &department.Department{
				Code: "IT",
				Name: "Information Technology",
			},
			setupMocks: func(dep *store.MockDepartment) {
				dep.EXPECT().
					ExistsByName(gomock.Any(), "Information Technology", nil).
					Return(false, nil)

				dep.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&department.Department{Code: "IT"}, nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockDep := store.NewMockDepartment(ctrl)
			mockEmp := store.NewMockEmployee(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockDep)
			}

			svc := New(mockDep, mockEmp)

			_, err := svc.Create(ctx, tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDepartmentService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	mockDep := store.NewMockDepartment(ctrl)
	mockEmp := store.NewMockEmployee(ctrl)

	mockDep.EXPECT().
		Get(gomock.Any()).
		Return([]*department.Department{
			{Code: "IT"},
		}, nil)

	svc := New(mockDep, mockEmp)

	res, err := svc.Get(ctx)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestDepartmentService_GetByCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	tests := []struct {
		desc        string
		code        string
		setupMocks  func(dep *store.MockDepartment)
		expectError bool
	}{
		{
			desc: "success",
			code: "IT",
			setupMocks: func(dep *store.MockDepartment) {
				dep.EXPECT().
					GetByCode(gomock.Any(), "IT").
					Return(&department.Department{Code: "IT"}, nil)
			},
		},
		{
			desc: "not found",
			code: "IT",
			setupMocks: func(dep *store.MockDepartment) {
				dep.EXPECT().
					GetByCode(gomock.Any(), "IT").
					Return(nil, errors.EntityNotFound{Entity: "IT"})
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockDep := store.NewMockDepartment(ctrl)
			mockEmp := store.NewMockEmployee(ctrl)

			tt.setupMocks(mockDep)

			svc := New(mockDep, mockEmp)

			_, err := svc.GetByCode(ctx, tt.code)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDepartmentService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	tests := []struct {
		desc        string
		code        string
		input       *department.NewDepartment
		setupMocks  func(dep *store.MockDepartment)
		expectError bool
	}{
		{
			desc:  "success",
			code:  "IT",
			input: &department.NewDepartment{Name: "Updated"},
			setupMocks: func(dep *store.MockDepartment) {
				dep.EXPECT().
					Update(gomock.Any(), "IT", gomock.Any()).
					Return(&department.Department{Code: "IT"}, nil)
			},
		},
		{
			desc:  "update error",
			code:  "IT",
			input: &department.NewDepartment{Name: "Updated"},
			setupMocks: func(dep *store.MockDepartment) {
				dep.EXPECT().
					Update(gomock.Any(), "IT", gomock.Any()).
					Return(nil, errors.DB{Err: errors.DB{}})
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockDep := store.NewMockDepartment(ctrl)
			mockEmp := store.NewMockEmployee(ctrl)

			tt.setupMocks(mockDep)

			svc := New(mockDep, mockEmp)

			_, err := svc.Update(ctx, tt.code, tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDepartmentService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	tests := []struct {
		desc        string
		code        string
		setupMocks  func(dep *store.MockDepartment, emp *store.MockEmployee)
		expectError bool
	}{
		{
			desc: "employees mapped",
			code: "IT",
			setupMocks: func(dep *store.MockDepartment, emp *store.MockEmployee) {
				emp.EXPECT().
					CountByDepartment(gomock.Any(), "IT").
					Return(2, nil)
			},
			expectError: true,
		},
		{
			desc: "count error",
			code: "IT",
			setupMocks: func(dep *store.MockDepartment, emp *store.MockEmployee) {
				emp.EXPECT().
					CountByDepartment(gomock.Any(), "IT").
					Return(0, errors.DB{Err: errors.DB{}})
			},
			expectError: true,
		},
		{
			desc: "success",
			code: "IT",
			setupMocks: func(dep *store.MockDepartment, emp *store.MockEmployee) {
				emp.EXPECT().
					CountByDepartment(gomock.Any(), "IT").
					Return(0, nil)

				dep.EXPECT().
					Delete(gomock.Any(), "IT").
					Return("Department deleted successfully", nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockDep := store.NewMockDepartment(ctrl)
			mockEmp := store.NewMockEmployee(ctrl)

			tt.setupMocks(mockDep, mockEmp)

			svc := New(mockDep, mockEmp)

			_, err := svc.Delete(ctx, tt.code)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
