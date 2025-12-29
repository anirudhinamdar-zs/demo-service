package store

import (
	"context"
	"demo-service/models/department"
	"demo-service/models/employee"
)

type Employee interface {
	Create(ctx context.Context, emp *employee.NewEmployee) (*employee.Employee, error)

	Get(ctx context.Context, filter employee.Filter) ([]*employee.Employee, error)

	GetById(ctx context.Context, id int) (*employee.Employee, error)

	Update(ctx context.Context, id int, emp *employee.NewEmployee) (*employee.Employee, error)

	Delete(ctx context.Context, id int) (string, error)

	ExistsByEmail(ctx context.Context, email string, excludeID *int) (bool, error)

	CountByDepartment(ctx context.Context, deptCode string) (int, error)
}

type Department interface {
	Create(ctx context.Context, dep *department.Department) (*department.Department, error)

	Get(ctx context.Context) ([]*department.Department, error)

	GetByCode(ctx context.Context, code string) (*department.Department, error)

	Update(ctx context.Context, code string, dep *department.NewDepartment) (*department.Department, error)

	Delete(ctx context.Context, code string) (string, error)

	ExistsByName(ctx context.Context, name string, excludeCode *string) (bool, error)
}
