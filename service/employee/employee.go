package employee

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"demo-service/models/department"
	"demo-service/models/employee"
	"demo-service/store"
)

type Employee struct {
	store           store.Employee
	departmentStore store.Department
}

func New(store store.Employee, departmentStore store.Department) *Employee {
	return &Employee{store: store, departmentStore: departmentStore}
}

func (e *Employee) Create(ctx *gofr.Context, emp *employee.NewEmployee) (*employee.Employee, error) {
	if !department.IsValidCode(emp.Department) {
		return nil, errors.InvalidParam{Param: []string{"department"}}
	}

	// Department must exist
	if _, err := e.departmentStore.GetByCode(ctx, emp.Department); err != nil {
		return nil, errors.EntityNotFound{Entity: emp.Department}
	}

	// Email uniqueness
	exists, err := e.store.ExistsByEmail(ctx, emp.Email, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.EntityAlreadyExists{}
	}

	return e.store.Create(ctx, emp)
}

func (e *Employee) Get(ctx *gofr.Context, filter employee.Filter) ([]*employee.Employee, error) {
	// Validate department filter
	if filter.Department != nil {
		if !department.IsValidCode(*filter.Department) {
			return nil, errors.EntityNotFound{Entity: *filter.Department}
		}
	}

	return e.store.Get(ctx, filter)
}

func (e *Employee) GetById(ctx *gofr.Context, employeeId int) (*employee.Employee, error) {
	return e.store.GetById(ctx, employeeId)
}

func (e *Employee) Update(ctx *gofr.Context, id int, emp *employee.NewEmployee) (*employee.Employee, error) {
	if emp.Department != "" {
		if !department.IsValidCode(emp.Department) {
			return nil, errors.EntityNotFound{Entity: emp.Department}
		}

		if _, err := e.departmentStore.GetByCode(ctx, emp.Department); err != nil {
			return nil, errors.EntityNotFound{Entity: emp.Department}
		}
	}

	if emp.Email != "" {
		exists, err := e.store.ExistsByEmail(ctx, emp.Email, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.EntityAlreadyExists{}
		}
	}

	return e.store.Update(ctx, id, emp)
}

func (e *Employee) Delete(ctx *gofr.Context, employeeId int) (string, error) {
	return e.store.Delete(ctx, employeeId)
}
