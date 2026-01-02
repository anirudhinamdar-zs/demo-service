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
		desc    string
		input   *department.Department
		mock    func(dep *store.MockDepartment)
		result  *department.Department
		mockErr error
	}{
		{
			desc: "invalid department code",
			input: &department.Department{
				Code: "BIO",
				Name: "Biology",
			},
			mock:    func(dep *store.MockDepartment) {},
			result:  nil,
			mockErr: errors.InvalidParam{Param: []string{"code"}},
		},
		{
			desc: "department already exists",
			input: &department.Department{
				Code: "IT",
				Name: "Information Technology",
			},
			mock: func(dep *store.MockDepartment) {
				dep.EXPECT().
					ExistsByName(gomock.Any(), "Information Technology", nil).
					Return(true, nil)
			},
			result:  nil,
			mockErr: errors.EntityAlreadyExists{},
		},
		{
			desc: "exists check error",
			input: &department.Department{
				Code: "IT",
				Name: "Information Technology",
			},
			mock: func(dep *store.MockDepartment) {
				dep.EXPECT().
					ExistsByName(gomock.Any(), "Information Technology", nil).
					Return(false, errors.DB{Err: errors.DB{}})
			},
			result:  nil,
			mockErr: errors.DB{Err: errors.DB{}},
		},
		{
			desc: "success",
			input: &department.Department{
				Code: "IT",
				Name: "Information Technology",
			},
			mock: func(dep *store.MockDepartment) {
				dep.EXPECT().
					ExistsByName(gomock.Any(), "Information Technology", nil).
					Return(false, nil)

				dep.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&department.Department{Code: "IT"}, nil)
			},
			result:  &department.Department{Code: "IT"},
			mockErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockDep := store.NewMockDepartment(ctrl)
			mockEmp := store.NewMockEmployee(ctrl)

			tt.mock(mockDep)

			svc := New(mockDep, mockEmp)

			resp, err := svc.Create(ctx, tt.input)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.result, resp)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, tt.mockErr, resp)
		})
	}
}

func TestDepartmentService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	tests := []struct {
		desc    string
		result  []*department.Department
		mockErr error
		mock    func(dep *store.MockDepartment)
	}{
		{
			desc:   "success",
			result: []*department.Department{{Code: "IT"}},
			mock: func(dep *store.MockDepartment) {
				dep.EXPECT().
					Get(gomock.Any()).
					Return([]*department.Department{{Code: "IT"}}, nil)
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockDep := store.NewMockDepartment(ctrl)
			mockEmp := store.NewMockEmployee(ctrl)

			tt.mock(mockDep)

			svc := New(mockDep, mockEmp)

			resp, err := svc.Get(ctx)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, resp, tt.result)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, err, tt.mockErr)
		})
	}
}

func TestDepartmentService_GetByCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	tests := []struct {
		desc    string
		code    string
		result  *department.Department
		mockErr error
		mock    func(dep *store.MockDepartment)
	}{
		{
			desc:   "success",
			code:   "IT",
			result: &department.Department{Code: "IT"},
			mock: func(dep *store.MockDepartment) {
				dep.EXPECT().
					GetByCode(gomock.Any(), "IT").
					Return(&department.Department{Code: "IT"}, nil)
			},
		},
		{
			desc:    "not found",
			code:    "IT",
			result:  nil,
			mockErr: errors.EntityNotFound{Entity: "IT"},
			mock: func(dep *store.MockDepartment) {
				dep.EXPECT().
					GetByCode(gomock.Any(), "IT").
					Return(nil, errors.EntityNotFound{Entity: "IT"})
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockDep := store.NewMockDepartment(ctrl)
			mockEmp := store.NewMockEmployee(ctrl)

			tt.mock(mockDep)

			svc := New(mockDep, mockEmp)

			resp, err := svc.GetByCode(ctx, tt.code)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, resp, tt.result)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, err, tt.mockErr)
		})
	}
}

func TestDepartmentService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	tests := []struct {
		desc    string
		code    string
		input   *department.NewDepartment
		result  *department.Department
		mockErr error
		mock    func(dep *store.MockDepartment)
	}{
		{
			desc:   "success",
			code:   "IT",
			input:  &department.NewDepartment{Name: "Updated"},
			result: &department.Department{Code: "IT"},
			mock: func(dep *store.MockDepartment) {
				dep.EXPECT().
					Update(gomock.Any(), "IT", gomock.Any()).
					Return(&department.Department{Code: "IT"}, nil)
			},
		},
		{
			desc:    "update error",
			code:    "IT",
			input:   &department.NewDepartment{Name: "Updated"},
			result:  nil,
			mockErr: errors.DB{Err: errors.DB{}},
			mock: func(dep *store.MockDepartment) {
				dep.EXPECT().
					Update(gomock.Any(), "IT", gomock.Any()).
					Return(nil, errors.DB{Err: errors.DB{}})
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockDep := store.NewMockDepartment(ctrl)
			mockEmp := store.NewMockEmployee(ctrl)

			tt.mock(mockDep)

			svc := New(mockDep, mockEmp)

			resp, err := svc.Update(ctx, tt.code, tt.input)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, resp, tt.result)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, err, tt.mockErr)
		})
	}
}

func TestDepartmentService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := &gofr.Context{}

	tests := []struct {
		desc    string
		code    string
		result  string
		mockErr error
		mock    func(dep *store.MockDepartment, emp *store.MockEmployee)
	}{
		{
			desc:    "employees mapped",
			code:    "IT",
			result:  "",
			mockErr: errors.Error("department has employees mapped"),
			mock: func(dep *store.MockDepartment, emp *store.MockEmployee) {
				emp.EXPECT().
					CountByDepartment(gomock.Any(), "IT").
					Return(2, nil)
			},
		},
		{
			desc:    "count error",
			code:    "IT",
			result:  "",
			mockErr: errors.DB{Err: errors.DB{}},
			mock: func(dep *store.MockDepartment, emp *store.MockEmployee) {
				emp.EXPECT().
					CountByDepartment(gomock.Any(), "IT").
					Return(0, errors.DB{Err: errors.DB{}})
			},
		},
		{
			desc:   "success",
			code:   "IT",
			result: "Department deleted successfully",
			mock: func(dep *store.MockDepartment, emp *store.MockEmployee) {
				emp.EXPECT().
					CountByDepartment(gomock.Any(), "IT").
					Return(0, nil)

				dep.EXPECT().
					Delete(gomock.Any(), "IT").
					Return("Department deleted successfully", nil)
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			mockDep := store.NewMockDepartment(ctrl)
			mockEmp := store.NewMockEmployee(ctrl)

			tt.mock(mockDep, mockEmp)

			svc := New(mockDep, mockEmp)

			resp, err := svc.Delete(ctx, tt.code)

			assert.Equalf(t, tt.result, resp, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, resp, tt.result)
			assert.Equalf(t, tt.mockErr, err, "Failed [%v]:%v \t Got: %v \t Expected: %v", tt.desc, i, err, tt.mockErr)
		})
	}
}
