package migration

import (
	"database/sql"
	"embed"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migration/*.sql
var dbMigrations embed.FS

func Initiator(dbParam *sql.DB) {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "sql_migration",
	}

	n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)
	if errs != nil {
		panic(errs)
	}

	fmt.Println("Migration success, applied", n, "migrations!")
}
