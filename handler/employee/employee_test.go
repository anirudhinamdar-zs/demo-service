package employee

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"

	"demo-service/models/employee"
	"demo-service/service"
)

func Initialize(t *testing.T) (*service.MockEmployee, Handler) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockEmployee(ctrl)
	handler := Handler{service: mockService}
	return mockService, handler
}

func TestCreate(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc        string
		body        []byte
		mock        func()
		expectedRes interface{}
		expectedErr error
	}{
		{
			desc: "success",
			body: []byte(`{"name":"John","email":"john@test.com","department":"IT"}`),
			mock: func() {
				mockService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(&employee.Employee{ID: 1}, nil)
			},
			expectedRes: &employee.Employee{ID: 1},
			expectedErr: nil,
		},
		{
			desc:        "bind error",
			body:        []byte(`invalid-json`),
			expectedRes: nil,
			expectedErr: errors.Error("Binding failed"),
		},
		{
			desc: "service error",
			body: []byte(`{"name":"John","email":"john@test.com","department":"IT"}`),
			mock: func() {
				mockService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.Error("service error"))
			},
			expectedRes: nil,
			expectedErr: errors.Error("service error"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodPost, "/employee", bytes.NewReader(tt.body))
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()

			resp, err := handler.Create(ctx)

			assert.Equalf(t, tt.expectedRes, resp, "case %d response mismatch", i)
			assert.Equalf(t, tt.expectedErr, err, "case %d error mismatch", i)
		})
	}
}

func TestGet(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc        string
		query       string
		mock        func()
		expectedRes interface{}
		expectedErr error
	}{
		{
			desc:  "success",
			query: "?department=IT",
			mock: func() {
				mockService.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return([]*employee.Employee{{ID: 1}}, nil)
			},
			expectedRes: map[string]interface{}{
				"data": map[string]interface{}{
					"employees": []*employee.Employee{{ID: 1}},
				},
			},
			expectedErr: nil,
		},
		{
			desc: "service error",
			mock: func() {
				mockService.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, errors.Error("service error"))
			},
			expectedRes: nil,
			expectedErr: errors.Error("service error"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodGet, "/employee"+tt.query, nil)
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()

			resp, err := handler.Get(ctx)

			assert.Equalf(t, tt.expectedRes, resp, "case %d response mismatch", i)
			assert.Equalf(t, tt.expectedErr, err, "case %d error mismatch", i)
		})
	}
}

func TestGetById(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc        string
		id          string
		mock        func()
		expectedRes interface{}
		expectedErr error
	}{
		{
			desc: "success",
			id:   "1",
			mock: func() {
				mockService.EXPECT().
					GetById(gomock.Any(), 1).
					Return(&employee.Employee{ID: 1}, nil)
			},
			expectedRes: &employee.Employee{ID: 1},
			expectedErr: nil,
		},
		{
			desc:        "invalid id",
			id:          "abc",
			expectedRes: nil,
			expectedErr: &strconv.NumError{Func: "Atoi", Num: "abc", Err: strconv.ErrSyntax},
		},
		{
			desc: "service error",
			id:   "1",
			mock: func() {
				mockService.EXPECT().
					GetById(gomock.Any(), 1).
					Return(nil, errors.Error("service error"))
			},
			expectedRes: nil,
			expectedErr: errors.Error("service error"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodGet, "/employee/1", nil)
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()
			ctx.SetPathParams(map[string]string{"id": tt.id})

			resp, err := handler.GetById(ctx)

			assert.Equalf(t, tt.expectedRes, resp, "case %d response mismatch", i)
			assert.Equalf(t, tt.expectedErr, err, "case %d error mismatch", i)
		})
	}
}

func TestUpdate(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc        string
		id          string
		body        []byte
		mock        func()
		expectedRes interface{}
		expectedErr error
	}{
		{
			desc: "success",
			id:   "1",
			body: []byte(`{"name":"John"}`),
			mock: func() {
				mockService.EXPECT().
					Update(gomock.Any(), 1, gomock.Any()).
					Return(&employee.Employee{ID: 1}, nil)
			},
			expectedRes: &employee.Employee{ID: 1},
			expectedErr: nil,
		},
		{
			desc:        "invalid id",
			id:          "abc",
			expectedRes: nil,
			expectedErr: &strconv.NumError{Func: "Atoi", Num: "abc", Err: strconv.ErrSyntax},
		},
		{
			desc:        "bind error",
			id:          "1",
			body:        []byte(`invalid`),
			expectedRes: nil,
			expectedErr: errors.Error("Binding failed"),
		},
		{
			desc: "service error",
			id:   "1",
			body: []byte(`{"name":"John"}`),
			mock: func() {
				mockService.EXPECT().
					Update(gomock.Any(), 1, gomock.Any()).
					Return(nil, errors.Error("service error"))
			},
			expectedRes: nil,
			expectedErr: errors.Error("service error"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodPut, "/employee/1", bytes.NewReader(tt.body))
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()
			ctx.SetPathParams(map[string]string{"id": tt.id})

			resp, err := handler.Update(ctx)

			assert.Equalf(t, tt.expectedRes, resp, "case %d response mismatch", i)
			assert.Equalf(t, tt.expectedErr, err, "case %d error mismatch", i)
		})
	}
}

func TestDelete(t *testing.T) {
	mockService, handler := Initialize(t)

	tests := []struct {
		desc        string
		id          string
		mock        func()
		expectedRes interface{}
		expectedErr error
	}{
		{
			desc: "success",
			id:   "1",
			mock: func() {
				mockService.EXPECT().
					Delete(gomock.Any(), 1).
					Return("deleted", nil)
			},
			expectedRes: "deleted",
			expectedErr: nil,
		},
		{
			desc:        "invalid id",
			id:          "abc",
			expectedRes: nil,
			expectedErr: &strconv.NumError{"Atoi", "abc", strconv.ErrSyntax},
		},
		{
			desc: "service error",
			id:   "1",
			mock: func() {
				mockService.EXPECT().
					Delete(gomock.Any(), 1).
					Return("", errors.Error("service error"))
			},
			expectedRes: nil,
			expectedErr: errors.Error("service error"),
		},
	}

	for i, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req := httptest.NewRequest(http.MethodDelete, "/employee/1", nil)
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()
			ctx.SetPathParams(map[string]string{"id": tt.id})

			resp, err := handler.Delete(ctx)

			assert.Equalf(t, tt.expectedRes, resp, "case %d response mismatch", i)
			assert.Equalf(t, tt.expectedErr, err, "case %d error mismatch", i)
		})
	}
}
