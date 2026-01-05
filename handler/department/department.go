package department

import (
	"strings"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"demo-service/models/department"
	"demo-service/service"
)

type Handler struct {
	service service.Department
}

func New(service service.Department) *Handler {
	return &Handler{service: service}
}

func validateDepartment(dep *department.Department) error {
	if dep == nil {
		return errors.MissingParam{}
	}

	if strings.TrimSpace(dep.Code) == "" {
		return errors.MissingParam{}
	}

	if strings.TrimSpace(dep.Name) == "" {
		return errors.MissingParam{}
	}

	if dep.Floor <= 0 {
		return errors.InvalidParam{}
	}

	return nil
}

func (h *Handler) Create(ctx *gofr.Context) (interface{}, error) {
	var dep *department.Department

	if err := ctx.Bind(&dep); err != nil {
		return nil, errors.Error("Binding failed")
	}

	if err := validateDepartment(dep); err != nil {
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
		return nil, errors.Error("Binding failed")
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
