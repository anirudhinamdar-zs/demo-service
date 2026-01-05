package migrations

import (
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type K20251231181007 struct {
}

func (k K20251231181007) Up(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(`
		CREATE TABLE IF NOT EXISTS employees (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			phone_number VARCHAR(10),
			dob DATE,
			major VARCHAR(100),
			city VARCHAR(100),
			department VARCHAR(10) NOT NULL,
		    deleted_at Date NULL,
			CONSTRAINT fk_employee_department
				FOREIGN KEY (department)
				REFERENCES departments(code)
		);
	`)

	return err
}

func (k K20251231181007) Down(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(`
		DROP TABLE employees;
	`)

	return err
}
