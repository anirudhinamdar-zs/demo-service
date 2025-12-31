package main

import (
	"demo-service/migrations"

	"developer.zopsmart.com/go/gofr/cmd/gofr/migration"
	dbmigration "developer.zopsmart.com/go/gofr/cmd/gofr/migration/dbMigration"
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

	syncMigration(app)
	db := dbmigration.NewGorm(app.GORM())

	err := migration.Migrate("demo-service", db, migrations.All(), "UP", app.Logger)

	if err != nil {
		app.Logger.Errorf("Error from migration: %v", err)
	}

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

func syncMigration(ctx *gofr.Gofr) {
	query := "CREATE TABLE IF NOT EXISTS `gofr_migrations`(" +
		"`app` varchar(191) NOT NULL," +
		"`version` bigint(20) NOT NULL," +
		"`start_time` datetime(3) DEFAULT NULL," +
		"`end_time` datetime(3) DEFAULT NULL," +
		"`method` varchar(191) NOT NULL," +
		"PRIMARY KEY (`app`,`version`,`method`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8;"

	_, err := ctx.DB().Exec(query)
	if err != nil {
		ctx.Logger.Errorf("error while create gofr migration table in order service.[Error]%v", err)
		return
	}

	query = "INSERT IGNORE INTO gofr_migrations(app,version,start_time,end_time,method) " +
		"SELECT  app,version,start_time,end_time,'UP' as method FROM zs_migrations ;"

	_, err = ctx.DB().Exec(query)
	if err != nil {
		ctx.Logger.Errorf("error while insert data in gofr migration table in order service.[Error]%v", err)
		return
	}
}
