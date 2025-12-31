package department

const createQuery = `INSERT INTO departments (code, name, floor, description) VALUES (?, ?, ?, ?)`

const selectQuery = `
	SELECT code, name, floor, description
	FROM departments
`

const selectQueryByCode = `
	SELECT code, name, floor, description
	FROM departments
	WHERE code = ?
`

const updateQueryByCode = `
	UPDATE departments
	SET name = ?, floor = ?, description = ?
	WHERE code = ?
`

const countQueryByCode = `
	SELECT COUNT(1)
	FROM employees
	WHERE department = ?
`

const deleteQuery = `
	DELETE FROM departments
	WHERE code = ?
`

var countByDepartmentQuery = `SELECT COUNT(1) FROM departments WHERE name = ?`
