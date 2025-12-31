// This is auto-generated file using 'gofr migrate' tool. DO NOT EDIT.
package migrations

import (
	"developer.zopsmart.com/go/gofr/cmd/gofr/migration/dbMigration"
)

func All() map[string]dbmigration.Migrator {
	return map[string]dbmigration.Migrator{

		"20251231180829": K20251231180829{},
		"20251231181007": K20251231181007{},
	}
}
