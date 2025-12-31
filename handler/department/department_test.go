package department

import (
	"bytes"
	"context"
	"demo-service/models/department"
	"demo-service/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
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
		desc        string
		body        []byte
		mock        func()
		expectErr   error
		expectedRes interface{}
	}{
		{
			desc: "success",
			body: []byte(`{"code":"IT","name":"Information Technology"}`),
			mock: func() {
				mockService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&department.Department{Code: "IT"}, nil)
			},
			expectedRes: &department.Department{Code: "IT"},
		},
		{
			desc:      "bind error",
			body:      []byte(`invalid-json`),
			expectErr: errors.Error("invalid character 'i' looking for beginning of value"),
		},
		{
			desc: "service error",
			body: []byte(`{"code":"IT","name":"Information Technology"}`),
			mock: func() {
				mockService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.Error("service error"))
			},
			expectErr: errors.Error("service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodPost, "/department", bytes.NewReader(tt.body))
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()

			resp, err := handler.Create(ctx)

			if tt.expectErr != nil {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.expectedRes, resp)
		})
	}
}

func TestGet(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc        string
		mock        func()
		expectErr   error
		expectedRes interface{}
	}{
		{
			desc: "success",
			mock: func() {
				mockService.EXPECT().
					Get(gomock.Any()).
					Return([]*department.Department{
						{Code: "IT"},
					}, nil)
			},
			expectErr:   nil,
			expectedRes: []*department.Department{{Code: "IT"}},
		},
		{
			desc: "service error",
			mock: func() {
				mockService.EXPECT().
					Get(gomock.Any()).
					Return(nil, errors.Error("service error"))
			},
			expectErr: errors.Error("service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodGet, "/department", nil)
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()

			resp, err := handler.Get(ctx)

			if tt.expectErr != nil {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.expectedRes, resp)
		})
	}
}

func TestGetByCode(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc        string
		code        string
		mock        func()
		expectErr   error
		expectedRes interface{}
	}{
		{
			desc: "success",
			code: "IT",
			mock: func() {
				mockService.EXPECT().
					GetByCode(gomock.Any(), "IT").
					Return(&department.Department{Code: "IT"}, nil)
			},
			expectErr:   nil,
			expectedRes: &department.Department{Code: "IT"},
		},
		{
			desc: "service error",
			code: "IT",
			mock: func() {
				mockService.EXPECT().
					GetByCode(gomock.Any(), "IT").
					Return(nil, errors.Error("service error"))
			},
			expectErr: errors.Error("service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodGet, "/department/IT", nil)
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()
			ctx.SetPathParams(map[string]string{"code": tt.code})

			resp, err := handler.GetByCode(ctx)

			if tt.expectErr != nil {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.expectedRes, resp)
		})
	}
}

func TestUpdate(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc        string
		code        string
		body        []byte
		mock        func()
		expectErr   error
		expectedRes interface{}
	}{
		{
			desc: "success",
			code: "IT",
			body: []byte(`{"name":"Tech"}`),
			mock: func() {
				mockService.EXPECT().
					Update(gomock.Any(), "IT", gomock.Any()).
					Return(&department.Department{Code: "IT"}, nil)
			},
			expectedRes: &department.Department{Code: "IT"},
		},
		{
			desc:      "bind error",
			code:      "IT",
			body:      []byte(`invalid`),
			expectErr: errors.Error("invalid character 'i' looking for beginning of value"),
		},
		{
			desc: "service error",
			code: "IT",
			body: []byte(`{"name":"Tech"}`),
			mock: func() {
				mockService.EXPECT().
					Update(gomock.Any(), "IT", gomock.Any()).
					Return(nil, errors.Error("service error"))
			},
			expectErr: errors.Error("service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodPut, "/department/IT", bytes.NewReader(tt.body))
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()
			ctx.SetPathParams(map[string]string{"code": tt.code})

			resp, err := handler.Update(ctx)

			if tt.expectErr != nil {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.expectedRes, resp)
		})
	}
}

func TestDelete(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc        string
		code        string
		mock        func()
		expectErr   error
		expectedRes interface{}
	}{
		{
			desc: "success",
			code: "IT",
			mock: func() {
				mockService.EXPECT().
					Delete(gomock.Any(), "IT").
					Return("deleted", nil)
			},
			expectErr:   nil,
			expectedRes: "deleted",
		},
		{
			desc: "service error",
			code: "IT",
			mock: func() {
				mockService.EXPECT().
					Delete(gomock.Any(), "IT").
					Return("", errors.Error("service error"))
			},
			expectErr: errors.Error("service error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodDelete, "/department/IT", nil)
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()
			ctx.SetPathParams(map[string]string{"code": tt.code})

			resp, err := handler.Delete(ctx)

			if tt.expectErr != nil {
				assert.Error(t, err)
				return
			}

			assert.Equal(t, tt.expectedRes, resp)
		})
	}
}
