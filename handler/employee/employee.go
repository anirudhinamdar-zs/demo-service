package employee

import (
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"demo-service/handler"
	"demo-service/models/employee"
	"demo-service/service"
)

type Handler struct {
	service service.Employee
}

func New(service service.Employee) handler.Employee {
	return &Handler{service: service}
}

func (h *Handler) Create(ctx *gofr.Context) (interface{}, error) {
	var emp *employee.NewEmployee

	if err := ctx.Bind(&emp); err != nil {
		return nil, err
	}

	resp, err := h.service.Create(ctx, emp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) Get(ctx *gofr.Context) (interface{}, error) {
	var filter employee.Filter

	if id := ctx.Param("id"); id != "" {
		parsed, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		filter.ID = &parsed
	}

	if name := ctx.Param("name"); name != "" {
		filter.Name = &name
	}

	if dept := ctx.Param("department"); dept != "" {
		filter.Department = &dept
	}

	resp, err := h.service.Get(ctx, filter)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"data": map[string]interface{}{
			"employees": resp,
		},
	}, nil
}

func (h *Handler) GetById(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	empID, err := strconv.Atoi(id)

	if err != nil {
		return nil, err
	}

	resp, err := h.service.GetById(ctx, empID)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) Update(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	empID, err := strconv.Atoi(id)

	if err != nil {
		return nil, err
	}

	var emp *employee.NewEmployee

	if err := ctx.Bind(&emp); err != nil {
		return nil, err
	}

	resp, err := h.service.Update(ctx, empID, emp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *Handler) Delete(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	empID, err := strconv.Atoi(id)

	if err != nil {
		return nil, err
	}

	resp, err := h.service.Delete(ctx, empID)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
