package main

import (
	"database/sql"

	"demo-service/handler/department"
	"demo-service/handler/employee"

	depService "demo-service/service/department"
	empService "demo-service/service/employee"

	depStore "demo-service/store/department"
	empStore "demo-service/store/employee"

	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gofr.dev/pkg/gofr"
)

func main() {
	fmt.Println("Hello World")

	app := gofr.New()

	dsn := "root:password@tcp(localhost:3306)/demo-service?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	//defer func(db *sql.DB) {
	//	err := db.Close()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}(db)

	employeeStore := empStore.Init(db)
	departmentStore := depStore.Init(db)

	employeeService := empService.New(employeeStore, departmentStore)
	departmentService := depService.New(departmentStore, employeeStore)

	employeeHandler := employee.New(employeeService)
	departmentHandler := department.New(departmentService)

	app.GET("/employees", employeeHandler.Get)
	app.GET("/employees/{id}", employeeHandler.GetById)
	app.POST("/employees", employeeHandler.Create)
	app.PUT("/employees/{id}", employeeHandler.Update)
	app.DELETE("/employees/{id}", employeeHandler.Delete)

	// ---- DEPARTMENT ROUTES ----
	app.GET("/departments", departmentHandler.Get)
	app.GET("/departments/{code}", departmentHandler.GetByCode)
	app.POST("/departments", departmentHandler.Create)
	app.PUT("/departments/{code}", departmentHandler.Update)
	app.DELETE("/departments/{code}", departmentHandler.Delete)

	app.Run()
}
