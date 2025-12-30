package department

import (
	"demo-service/models/department"
	"demo-service/store"
	"errors"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
)

type Department struct {
	store         store.Department
	employeeStore store.Employee
}

func New(store store.Department, employeeStore store.Employee) *Department {
	return &Department{store: store, employeeStore: employeeStore}
}

func (d *Department) Create(ctx *gofr.Context, dep *department.Department) (*department.Department, error) {
	if !department.IsValidCode(dep.Code) {
		return nil, errors.New("invalid department code")
	}

	// Optional but correct: uniqueness check
	exists, err := d.store.ExistsByName(ctx, dep.Name, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("department already exists")
	}

	return d.store.Create(ctx, dep)
}

func (d *Department) Get(ctx *gofr.Context) ([]*department.Department, error) {
	return d.store.Get(ctx)
}

func (d *Department) GetByCode(ctx *gofr.Context, code string) (*department.Department, error) {
	return d.store.GetByCode(ctx, code)
}

func (d *Department) Update(
	ctx *gofr.Context,
	code string,
	dep *department.NewDepartment,
) (*department.Department, error) {
	return d.store.Update(ctx, code, dep)
}

func (d *Department) Delete(ctx *gofr.Context, code string) (string, error) {
	count, err := d.employeeStore.CountByDepartment(ctx, code)
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "", errors.New("department has employees mapped")
	}

	return d.store.Delete(ctx, code)
}
