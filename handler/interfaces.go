package handler

import "gofr.dev/pkg/gofr"

type Employee interface {
	Create(ctx *gofr.Context) (interface{}, error)

	Get(ctx *gofr.Context) (interface{}, error)

	GetById(ctx *gofr.Context) (interface{}, error)

	Update(ctx *gofr.Context) (interface{}, error)

	Delete(ctx *gofr.Context) (interface{}, error)
}

type Department interface {
	Create(ctx *gofr.Context) (interface{}, error)

	Get(ctx *gofr.Context) (interface{}, error)

	GetByCode(ctx *gofr.Context) (interface{}, error)

	Update(ctx *gofr.Context) (interface{}, error)

	Delete(ctx *gofr.Context) (interface{}, error)
}
