package infra

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"

	"github.com/jmoiron/sqlx"
	"github.com/ray-d-song/yan/internal/embedfs"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB(config *Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect(config.DB.Driver, config.DB.DSN)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// WithTx executes a function within a database transaction.
// It automatically handles commit on success and rollback on error or panic.
//
// Example usage:
//
//	err := infra.WithTx(ctx, db, func(tx *sqlx.Tx) error {
//	    // Your transactional operations here
//	    return repo.CreateUser(ctx, tx, user)
//	})
func WithTx(
	ctx context.Context,
	db *sqlx.DB,
	fn func(tx *sqlx.Tx) error,
) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// Panic occurred, rollback and re-panic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// Error occurred, rollback
			_ = tx.Rollback()
		} else {
			// Success, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// AutoMigrate runs all SQL migration files from the embedded filesystem
func AutoMigrate(db *sqlx.DB, logger *Logger) error {
	logger.Info("Starting database migration...")

	// Read all .up.sql files from the embedded filesystem
	var migrations []string
	err := fs.WalkDir(embedfs.SQLFile, "sql/migrate", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".sql" && filepath.Base(path)[len(filepath.Base(path))-7:] == ".up.sql" {
			migrations = append(migrations, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// Sort migrations by filename to ensure they run in order
	sort.Strings(migrations)

	// Execute each migration
	for _, migrationPath := range migrations {
		logger.Infof("Running migration: %s", filepath.Base(migrationPath))

		// Read the SQL file content
		content, err := fs.ReadFile(embedfs.SQLFile, migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migrationPath, err)
		}

		// Execute the SQL
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migrationPath, err)
		}

		logger.Infof("Successfully executed migration: %s", filepath.Base(migrationPath))
	}

	logger.Infof("Database migration completed. Executed %d migrations.", len(migrations))
	return nil
}
