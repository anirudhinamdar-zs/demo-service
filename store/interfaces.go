package store

import (
	"demo-service/models/department"
	"demo-service/models/employee"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Employee interface {
	Create(ctx *gofr.Context, emp *employee.NewEmployee) (*employee.Employee, error)

	Get(ctx *gofr.Context, filter employee.Filter) ([]*employee.Employee, error)

	GetById(ctx *gofr.Context, id int) (*employee.Employee, error)

	Update(ctx *gofr.Context, id int, emp *employee.NewEmployee) (*employee.Employee, error)

	Delete(ctx *gofr.Context, id int) (string, error)

	ExistsByEmail(ctx *gofr.Context, email string, excludeID *int) (bool, error)

	CountByDepartment(ctx *gofr.Context, deptCode string) (int, error)
}

type Department interface {
	Create(ctx *gofr.Context, dep *department.Department) (*department.Department, error)

	Get(ctx *gofr.Context) ([]*department.Department, error)

	GetByCode(ctx *gofr.Context, code string) (*department.Department, error)

	Update(ctx *gofr.Context, code string, dep *department.NewDepartment) (*department.Department, error)

	Delete(ctx *gofr.Context, code string) (string, error)

	ExistsByName(ctx *gofr.Context, name string, excludeCode *string) (bool, error)
}
