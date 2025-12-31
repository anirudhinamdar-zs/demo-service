package migrations

import (
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/log"
)

type K20251231180829 struct {
}

func (k K20251231180829) Up(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(`
		CREATE TABLE IF NOT EXISTS departments (
			code VARCHAR(10) PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL,
			floor INT NOT NULL,
			description TEXT
		);
	`)

	return err
}

func (k K20251231180829) Down(d *datastore.DataStore, logger log.Logger) error {
	_, err := d.DB().Exec(`
		DROP TABLE departments;
	`)

	return err
}
