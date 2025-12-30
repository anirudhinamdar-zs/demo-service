package employee

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"

	"demo-service/models/employee"
	"demo-service/service"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockEmployee(ctrl)
	handler := Handler{service: mockService}

	tests := []struct {
		desc        string
		body        []byte
		mock        func()
		expectErr   bool
		expectedRes interface{}
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
		},
		{
			desc:      "bind error",
			body:      []byte(`invalid-json`),
			expectErr: true,
		},
		{
			desc: "service error",
			body: []byte(`{"name":"John","email":"john@test.com","department":"IT"}`),
			mock: func() {
				mockService.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("create failed"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req, _ := http.NewRequest(http.MethodPost, "/employee", bytes.NewReader(tt.body))
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()

			resp, err := handler.Create(ctx)

			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedRes, resp)
		})
	}
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockEmployee(ctrl)
	handler := Handler{service: mockService}

	tests := []struct {
		desc        string
		query       string
		mock        func()
		expectErr   bool
		expectedRes interface{}
	}{
		{
			desc:  "success with filters",
			query: "?department=IT",
			mock: func() {
				mockService.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return([]*employee.Employee{
						{ID: 1},
					}, nil)
			},
			expectedRes: map[string]interface{}{
				"data": map[string]interface{}{
					"employees": []*employee.Employee{{ID: 1}},
				},
			},
		},
		{
			desc: "service error",
			mock: func() {
				mockService.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("fetch failed"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req, _ := http.NewRequest(http.MethodGet, "/employee"+tt.query, nil)
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()

			resp, err := handler.Get(ctx)

			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedRes, resp)
		})
	}
}

func TestGetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockEmployee(ctrl)
	handler := Handler{service: mockService}

	tests := []struct {
		desc        string
		id          string
		mock        func()
		expectErr   bool
		expectedRes interface{}
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
		},
		{
			desc:      "invalid id",
			id:        "abc",
			expectErr: true,
		},
		{
			desc: "service error",
			id:   "1",
			mock: func() {
				mockService.EXPECT().
					GetById(gomock.Any(), 1).
					Return(nil, errors.New("not found"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req, _ := http.NewRequest(http.MethodGet, "/employee/1", nil)
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()
			ctx.SetPathParams(map[string]string{"id": tt.id})

			resp, err := handler.GetById(ctx)

			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedRes, resp)
		})
	}
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockEmployee(ctrl)
	handler := Handler{service: mockService}

	tests := []struct {
		desc        string
		id          string
		body        []byte
		mock        func()
		expectErr   bool
		expectedRes interface{}
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
		},
		{
			desc:      "invalid id",
			id:        "abc",
			expectErr: true,
		},
		{
			desc:      "bind error",
			id:        "1",
			body:      []byte(`invalid`),
			expectErr: true,
		},
		{
			desc: "service error",
			id:   "1",
			body: []byte(`{"name":"John"}`),
			mock: func() {
				mockService.EXPECT().
					Update(gomock.Any(), 1, gomock.Any()).
					Return(nil, errors.New("update failed"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req, _ := http.NewRequest(http.MethodPut, "/employee/1", bytes.NewReader(tt.body))
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()
			ctx.SetPathParams(map[string]string{"id": tt.id})

			resp, err := handler.Update(ctx)

			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedRes, resp)
		})
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockEmployee(ctrl)
	handler := Handler{service: mockService}

	tests := []struct {
		desc        string
		id          string
		mock        func()
		expectErr   bool
		expectedRes interface{}
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
		},
		{
			desc:      "invalid id",
			id:        "abc",
			expectErr: true,
		},
		{
			desc: "service error",
			id:   "1",
			mock: func() {
				mockService.EXPECT().
					Delete(gomock.Any(), 1).
					Return("", errors.New("delete failed"))
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock()
			}

			req, _ := http.NewRequest(http.MethodDelete, "/employee/1", nil)
			ctx := gofr.NewContext(nil, request.NewHTTPRequest(req), nil)
			ctx.Context = context.Background()
			ctx.SetPathParams(map[string]string{"id": tt.id})

			resp, err := handler.Delete(ctx)

			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedRes, resp)
		})
	}
}
