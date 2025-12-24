package migrations

import "database/sql"

func ManagementUp(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE departments (
			code VARCHAR(10) PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL,
			floor INT NOT NULL,
			description TEXT
		);

		CREATE TABLE employees (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			phone_number VARCHAR(10),
			dob DATE,
			major VARCHAR(100),
			city VARCHAR(100),
			department VARCHAR(10) NOT NULL,
			CONSTRAINT fk_employee_department
				FOREIGN KEY (department)
				REFERENCES departments(code)
		);
	`)
	return err
}

func ManagementDown(tx *sql.Tx) error {
	_, err := tx.Exec(`
		DROP TABLE employees;
		DROP TABLE departments;
	`)
	return err
}
