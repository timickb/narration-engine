package migrations

import (
	"embed"
	"github.com/timickb/go-stateflow/pkg/db"
)

//go:embed *.sql
var migrations embed.FS

var Migrator = db.NewMigrator(".", migrations)
