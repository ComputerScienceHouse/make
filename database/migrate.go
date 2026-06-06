package database

import (
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate() error {
	log.Println("[DB] Beginning database migration")

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		log.Printf("[DB] Failed to load migrations: %v\n", err)
		return err
	}

	log.Println("[DB] Running migrations...")

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		source,
		"sqlite://database.db",
	)
	if err != nil {
		log.Printf("[DB] Failed to create migrator: %v\n", err)
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Println("[DB] Database already up to date")
		return nil
	}

	if err != nil {
		log.Printf("[DB] Migration failed: %v\n", err)
		return err
	}

	version, dirty, verr := m.Version()
	if verr == nil {
		log.Printf("[DB] Migration complete. Current version: %d (dirty=%v)\n", version, dirty)
	} else {
		log.Println("[DB] Migration complete.")
	}

	return nil
}
