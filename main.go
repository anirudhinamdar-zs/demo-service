package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"demo-service/handler/department"
	"demo-service/handler/employee"
	depService "demo-service/service/department"
	empService "demo-service/service/employee"
	depStore "demo-service/store/department"
	empStore "demo-service/store/employee"
)

func main() {
	app := gofr.New()

	employeeStore := empStore.Init()
	departmentStore := depStore.Init()

	employeeService := empService.New(employeeStore, departmentStore)
	departmentService := depService.New(departmentStore, employeeStore)

	employeeHandler := employee.New(employeeService)
	departmentHandler := department.New(departmentService)

	app.GET("/employees", employeeHandler.Get)
	app.GET("/employees/{id}", employeeHandler.GetById)
	app.POST("/employees", employeeHandler.Create)
	app.PUT("/employees/{id}", employeeHandler.Update)
	app.DELETE("/employees/{id}", employeeHandler.Delete)

	app.GET("/departments", departmentHandler.Get)
	app.GET("/departments/{code}", departmentHandler.GetByCode)
	app.POST("/departments", departmentHandler.Create)
	app.PUT("/departments/{code}", departmentHandler.Update)
	app.DELETE("/departments/{code}", departmentHandler.Delete)

	app.Start()
}
