package employee

import (
	"demo-service/models/department"
	"demo-service/models/employee"
	"demo-service/store"
	"errors"

	"gofr.dev/pkg/gofr"
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
		return nil, errors.New("invalid department")
	}

	// Department must exist
	if _, err := e.departmentStore.GetByCode(ctx, emp.Department); err != nil {
		return nil, errors.New("department does not exist")
	}

	// Email uniqueness
	exists, err := e.store.ExistsByEmail(ctx, emp.Email, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	return e.store.Create(ctx, emp)
}

func (e *Employee) Get(ctx *gofr.Context, filter employee.Filter) ([]*employee.Employee, error) {
	// Validate department filter
	if filter.Department != nil {
		if !department.IsValidCode(*filter.Department) {
			return nil, errors.New("invalid department filter")
		}
	}

	return e.store.Get(ctx, filter)
}

func (e *Employee) GetById(ctx *gofr.Context, employeeId int) (*employee.Employee, error) {
	return e.store.GetById(ctx, employeeId)
}

func (e *Employee) Update(ctx *gofr.Context, id int, emp *employee.NewEmployee) (*employee.Employee, error) {
	if !department.IsValidCode(emp.Department) {
		return nil, errors.New("invalid department")
	}

	// Target department must exist
	if _, err := e.departmentStore.GetByCode(ctx, emp.Department); err != nil {
		return nil, errors.New("department does not exist")
	}

	// Email uniqueness except self
	exists, err := e.store.ExistsByEmail(ctx, emp.Email, &id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	return e.store.Update(ctx, id, emp)
}

func (e *Employee) Delete(ctx *gofr.Context, employeeId int) (string, error) {
	return e.store.Delete(ctx, employeeId)
}
