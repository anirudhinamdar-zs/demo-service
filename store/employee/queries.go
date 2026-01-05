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
		department,
		deleted_at
	FROM employees
`

const queryByID = `
	SELECT
		id,
		name,
		email,
		phone_number,
		dob,
		major,
		city,
		department,
		deleted_at
	FROM employees
	WHERE id = ? AND deleted_at IS NULL
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

const softDeleteQueryById = `
	UPDATE employees
	SET deleted_at = CURRENT_DATE
	WHERE id = ? AND deleted_at IS NULL
`

var selectQueryForEmail = `
	SELECT COUNT(1)
	FROM employees
	WHERE email = ? AND deleted_at IS NULL
`

const selectQueryCountByDepartment = `
	SELECT COUNT(1)
	FROM employees
	WHERE department = ? AND deleted_at IS NULL
`
