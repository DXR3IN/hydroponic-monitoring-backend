package migration

import (
	"gorm.io/gorm"
)

// SetupMigrations registers all migrations
func SetupMigrations(migrator *Migrator) {
	// Migration 001: Enable pgcrypto extension
	migrator.RegisterMigration(
		"001",
		"enable_pgcrypto_extension",
		func(db *gorm.DB) error {
			return db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error
		},
		func(db *gorm.DB) error {
			// Be careful with dropping extensions in production
			return db.Exec("DROP EXTENSION IF EXISTS \"pgcrypto\" CASCADE").Error
		},
	)

	// Migration 002: Create users table
	migrator.RegisterMigration(
		"002",
		"create_users_table",
		func(db *gorm.DB) error {
			// Create table using raw SQL to ensure exact schema
			sql := `
				CREATE TABLE IF NOT EXISTS users (
					id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
					name VARCHAR(255) NOT NULL,
					email VARCHAR(255) NOT NULL,
					password VARCHAR(255) NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				)
			`
			if err := db.Exec(sql).Error; err != nil {
				return err
			}

			// Create unique index
			return db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error
		},
		func(db *gorm.DB) error {
			return db.Exec("DROP TABLE IF EXISTS users CASCADE").Error
		},
	)

	// Example: Migration 003 - Add phone column (for future use)
	migrator.RegisterMigration(
		"003",
		"add_phone_to_users",
		func(db *gorm.DB) error {
			return db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20)").Error
		},
		func(db *gorm.DB) error {
			return db.Exec("ALTER TABLE users DROP COLUMN IF EXISTS phone").Error
		},
	)
}

// RunMigrations executes all pending migrations
func RunMigrations(db *gorm.DB) error {
	migrator := NewMigrator(db)
	SetupMigrations(migrator)
	return migrator.Up()
}

// RollbackLastMigration rolls back the last applied migration
func RollbackLastMigration(db *gorm.DB) error {
	migrator := NewMigrator(db)
	SetupMigrations(migrator)
	return migrator.Down()
}

// ShowMigrationStatus displays all applied migrations
func ShowMigrationStatus(db *gorm.DB) error {
	migrator := NewMigrator(db)
	return migrator.Status()
}