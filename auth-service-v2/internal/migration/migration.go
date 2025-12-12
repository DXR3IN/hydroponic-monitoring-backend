package migration

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// Migration tracks which migrations have been applied
type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Version   string    `gorm:"uniqueIndex;size:255;not null"`
	Name      string    `gorm:"size:255;not null"`
	AppliedAt time.Time `gorm:"not null"`
}

// MigrationFunc is a function that performs a migration
type MigrationFunc func(*gorm.DB) error

// MigrationDefinition defines a single migration
type MigrationDefinition struct {
	Version string
	Name    string
	Up      MigrationFunc
	Down    MigrationFunc
}

// Migrator handles database migrations
type Migrator struct {
	db         *gorm.DB
	migrations []MigrationDefinition
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: []MigrationDefinition{},
	}
}

// RegisterMigration adds a migration to the migrator
func (m *Migrator) RegisterMigration(version, name string, up, down MigrationFunc) {
	m.migrations = append(m.migrations, MigrationDefinition{
		Version: version,
		Name:    name,
		Up:      up,
		Down:    down,
	})
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	err := m.db.Exec(`
    CREATE TABLE IF NOT EXISTS migrations (
        id SERIAL PRIMARY KEY,
        version VARCHAR(255) NOT NULL,
        name VARCHAR(255) NOT NULL,
        applied_at TIMESTAMP WITH TIME ZONE NOT NULL
    );
    CREATE UNIQUE INDEX IF NOT EXISTS idx_migrations_version ON migrations(version);
	`).Error

	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	for _, migration := range m.migrations {
		// Check if migration already applied
		var count int64
		if err := m.db.Model(&Migration{}).Where("version = ?", migration.Version).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check migration %s: %w", migration.Version, err)
		}

		if count > 0 {
			log.Printf("Migration %s already applied, skipping", migration.Version)
			continue
		}

		// Run migration in transaction
		if err := m.db.Transaction(func(tx *gorm.DB) error {
			log.Printf("Applying migration %s: %s", migration.Version, migration.Name)

			if err := migration.Up(tx); err != nil {
				return fmt.Errorf("migration up failed: %w", err)
			}

			// Record migration as applied
			record := Migration{
				Version:   migration.Version,
				Name:      migration.Name,
				AppliedAt: time.Now(),
			}

			if err := tx.Create(&record).Error; err != nil {
				return fmt.Errorf("failed to record migration: %w", err)
			}

			return nil
		}); err != nil {
			return fmt.Errorf("migration %s failed: %w", migration.Version, err)
		}

		log.Printf("Successfully applied migration %s", migration.Version)
	}

	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	var lastMigration Migration
	if err := m.db.Order("applied_at DESC").First(&lastMigration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No migrations to roll back")
			return nil
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	// Find migration definition
	var migrationDef *MigrationDefinition
	for i := range m.migrations {
		if m.migrations[i].Version == lastMigration.Version {
			migrationDef = &m.migrations[i]
			break
		}
	}

	if migrationDef == nil {
		return fmt.Errorf("migration definition not found for version %s", lastMigration.Version)
	}

	// Run rollback in transaction
	if err := m.db.Transaction(func(tx *gorm.DB) error {
		log.Printf("Rolling back migration %s: %s", migrationDef.Version, migrationDef.Name)

		if err := migrationDef.Down(tx); err != nil {
			return fmt.Errorf("migration down failed: %w", err)
		}

		// Delete migration record
		if err := tx.Delete(&lastMigration).Error; err != nil {
			return fmt.Errorf("failed to delete migration record: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	log.Printf("Successfully rolled back migration %s", migrationDef.Version)
	return nil
}

// Status shows all applied migrations
func (m *Migrator) Status() error {
	var applied []Migration
	if err := m.db.Order("applied_at ASC").Find(&applied).Error; err != nil {
		return fmt.Errorf("failed to get migrations: %w", err)
	}

	if len(applied) == 0 {
		log.Println("No migrations applied yet")
		return nil
	}

	log.Println("Applied migrations:")
	for _, migration := range applied {
		log.Printf("  - %s: %s (applied at %s)", migration.Version, migration.Name, migration.AppliedAt.Format(time.RFC3339))
	}

	return nil
}
