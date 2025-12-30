package department

import (
	"demo-service/handler"
	"demo-service/models/department"
	"demo-service/service"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Handler struct {
	service service.Department
}

func New(service service.Department) handler.Department {
	return &Handler{service: service}
}

func (h *Handler) Create(ctx *gofr.Context) (interface{}, error) {
	var dep *department.Department

	if err := ctx.Bind(&dep); err != nil {
		return nil, err
	}

	resp, err := h.service.Create(ctx, dep)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) Get(ctx *gofr.Context) (interface{}, error) {
	resp, err := h.service.Get(ctx)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) GetByCode(ctx *gofr.Context) (interface{}, error) {
	code := ctx.PathParam("code")

	resp, err := h.service.GetByCode(ctx, code)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) Update(ctx *gofr.Context) (interface{}, error) {
	code := ctx.PathParam("code")

	var dep *department.NewDepartment

	if err := ctx.Bind(&dep); err != nil {
		return nil, err
	}

	resp, err := h.service.Update(ctx, code, dep)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) Delete(ctx *gofr.Context) (interface{}, error) {
	code := ctx.PathParam("code")

	resp, err := h.service.Delete(ctx, code)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
