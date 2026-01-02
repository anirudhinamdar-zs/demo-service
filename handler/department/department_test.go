package department

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"

	"demo-service/models/department"
	"demo-service/service"
)

func Initialize(t *testing.T) (*service.MockDepartment, Handler) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockDepartment(ctrl)
	handler := Handler{service: mockService}
	return mockService, handler
}

func TestCreate(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc    string
		body    []byte
		result  interface{}
		mockErr error
		mock    func()
	}{
		{
			desc: "success",
			body: []byte(`{"code":"IT","name":"Information Technology","floor":1,"description":"New IT Department"}`),
			result: &department.Department{
				Code: "IT", Name: "Information Technology", Floor: 1, Description: "New IT Department",
			},
			mock: func() {
				mockService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&department.Department{
						Code: "IT", Name: "Information Technology", Floor: 1, Description: "New IT Department",
					}, nil)
			},
		},
		{
			desc:    "bind error",
			body:    []byte(`invalid-json`),
			result:  nil,
			mockErr: errors.Error("Binding failed"),
		},
		{
			desc:    "Invalid param error",
			body:    []byte(`{"code":"BIO","name":"Biology Department"}`),
			result:  nil,
			mockErr: errors.InvalidParam(errors.InvalidParam{Param: []string(nil)}),
			mock: func() {
				mockService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.Error("service error"))
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodPost, "/department", bytes.NewReader(tt.body))
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()

			resp, err := handler.Create(ctx)

			assert.Equalf(t, tt.result, resp,
				"Failed [%v]:%v \t Got: %v \t Expected: %v",
				tt.desc, i, resp, tt.result)

			assert.Equalf(t, tt.mockErr, err,
				"Failed [%v]:%v \t Got: %v \t Expected: %v",
				tt.desc, i, err, tt.mockErr)
		})
	}
}

func TestGet(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc    string
		result  interface{}
		mockErr error
		mock    func()
	}{
		{
			desc:   "success",
			result: []*department.Department{{Code: "IT"}},
			mock: func() {
				mockService.EXPECT().
					Get(gomock.Any()).
					Return([]*department.Department{{Code: "IT"}}, nil)
			},
		},
		{
			desc:    "service error",
			result:  nil,
			mockErr: errors.Error("service error"),
			mock: func() {
				mockService.EXPECT().
					Get(gomock.Any()).
					Return(nil, errors.Error("service error"))
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodGet, "/department", nil)
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()

			resp, err := handler.Get(ctx)

			assert.Equalf(t, tt.result, resp,
				"Failed [%v]:%v \t Got: %v \t Expected: %v",
				tt.desc, i, resp, tt.result)

			assert.Equalf(t, tt.mockErr, err,
				"Failed [%v]:%v \t Got: %v \t Expected: %v",
				tt.desc, i, err, tt.mockErr)
		})
	}
}

func TestGetByCode(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc    string
		code    string
		result  interface{}
		mockErr error
		mock    func()
	}{
		{
			desc:   "success",
			code:   "IT",
			result: &department.Department{Code: "IT"},
			mock: func() {
				mockService.EXPECT().
					GetByCode(gomock.Any(), "IT").
					Return(&department.Department{Code: "IT"}, nil)
			},
		},
		{
			desc:    "service error",
			code:    "IT",
			result:  nil,
			mockErr: errors.Error("service error"),
			mock: func() {
				mockService.EXPECT().
					GetByCode(gomock.Any(), "IT").
					Return(nil, errors.Error("service error"))
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodGet, "/department/IT", nil)
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()
			ctx.SetPathParams(map[string]string{"code": tt.code})

			resp, err := handler.GetByCode(ctx)

			assert.Equalf(t, tt.result, resp,
				"Failed [%v]:%v \t Got: %v \t Expected: %v",
				tt.desc, i, resp, tt.result)

			assert.Equalf(t, tt.mockErr, err,
				"Failed [%v]:%v \t Got: %v \t Expected: %v",
				tt.desc, i, err, tt.mockErr)
		})
	}
}

func TestUpdate(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc    string
		code    string
		body    []byte
		result  interface{}
		mockErr error
		mock    func()
	}{
		{
			desc:   "success",
			code:   "IT",
			body:   []byte(`{"name":"Tech"}`),
			result: &department.Department{Code: "IT"},
			mock: func() {
				mockService.EXPECT().
					Update(gomock.Any(), "IT", gomock.Any()).
					Return(&department.Department{Code: "IT"}, nil)
			},
		},
		{
			desc:    "bind error",
			code:    "IT",
			body:    []byte(`invalid`),
			result:  nil,
			mockErr: errors.Error("Binding failed"),
		},
		{
			desc:    "service error",
			code:    "IT",
			body:    []byte(`{"name":"Tech"}`),
			result:  nil,
			mockErr: errors.Error("service error"),
			mock: func() {
				mockService.EXPECT().
					Update(gomock.Any(), "IT", gomock.Any()).
					Return(nil, errors.Error("service error"))
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodPut, "/department/IT", bytes.NewReader(tt.body))
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()
			ctx.SetPathParams(map[string]string{"code": tt.code})

			resp, err := handler.Update(ctx)

			assert.Equalf(t, tt.result, resp,
				"Failed [%v]:%v \t Got: %v \t Expected: %v",
				tt.desc, i, resp, tt.result)

			assert.Equalf(t, tt.mockErr, err,
				"Failed [%v]:%v \t Got: %v \t Expected: %v",
				tt.desc, i, err, tt.mockErr)
		})
	}
}

func TestDelete(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc    string
		code    string
		result  interface{}
		mockErr error
		mock    func()
	}{
		{
			desc:   "success",
			code:   "IT",
			result: "deleted",
			mock: func() {
				mockService.EXPECT().
					Delete(gomock.Any(), "IT").
					Return("deleted", nil)
			},
		},
		{
			desc:    "service error",
			code:    "IT",
			result:  nil,
			mockErr: errors.Error("service error"),
			mock: func() {
				mockService.EXPECT().
					Delete(gomock.Any(), "IT").
					Return("", errors.Error("service error"))
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodDelete, "/department/IT", nil)
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()
			ctx.SetPathParams(map[string]string{"code": tt.code})

			resp, err := handler.Delete(ctx)

			assert.Equalf(t, tt.result, resp,
				"Failed [%v]:%v \t Got: %v \t Expected: %v",
				tt.desc, i, resp, tt.result)

			assert.Equalf(t, tt.mockErr, err,
				"Failed [%v]:%v \t Got: %v \t Expected: %v",
				tt.desc, i, err, tt.mockErr)
		})
	}
}
