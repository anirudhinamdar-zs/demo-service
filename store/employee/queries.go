package employee

var baseGetEmployeesQuery = `
	SELECT
		id,
		name,
		email,
		phone_number,
		dob,
		major,
		city,
		department
	FROM employees
`

const queryByID = `
	SELECT id, name, email, phone_number, dob, major, city, department
	FROM employees
	WHERE id = ?
`

const updateQueryById = `
		UPDATE employees
		SET
			name = ?,
			email = ?,
			phone_number = ?,
			dob = ?,
			major = ?,
			city = ?,
			department = ?
		WHERE id = ?
	`

const deleteQueryById = `DELETE FROM employees WHERE id = ?`

var selectQueryForEmail = `SELECT COUNT(1) FROM employees WHERE email = ?`

const selectQueryCountByDepartment = `
		SELECT COUNT(1)
		FROM employees
		WHERE department = ?
	`
