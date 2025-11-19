package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal(err)
	}

	// Open SQLite connection
	db, err := sql.Open("sqlite3", "file:data/chat.db?_foreign_keys=1")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Ensure foreign keys are enforced
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Connected to SQLite database")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	fmt.Println("✅ All migrations completed successfully!")
}

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	Up      func(*sql.DB) error
}

// runMigrations executes all database migrations
func runMigrations(db *sql.DB) error {
	// Create migrations table to track which migrations have been run
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// List of all migrations
	migrations := []Migration{
		{Version: 1, Name: "create_users_table", Up: createUsersTable},
		{Version: 2, Name: "create_friendships_table", Up: createFriendshipsTable},
		// Messages now stored in MongoDB, not SQLite
		// Add new migrations here in the future
	}

	// Run each migration
	for _, migration := range migrations {
		if err := runMigration(db, migration); err != nil {
			return fmt.Errorf("migration %d (%s) failed: %v", migration.Version, migration.Name, err)
		}
	}

	return nil
}

// createMigrationsTable creates a table to track migration versions
func createMigrationsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}

// runMigration runs a single migration if it hasn't been run yet
func runMigration(db *sql.DB, migration Migration) error {
	// Check if migration has already been applied
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", migration.Version).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Printf("⏭️  Migration %d (%s) already applied, skipping", migration.Version, migration.Name)
		return nil
	}

	// Run the migration
	log.Printf("🔄 Running migration %d: %s", migration.Version, migration.Name)
	if err := migration.Up(db); err != nil {
		return err
	}

	// Record that migration was applied
	_, err = db.Exec("INSERT INTO schema_migrations (version, name) VALUES (?, ?)", migration.Version, migration.Name)
	if err != nil {
		return err
	}

	log.Printf("✅ Migration %d (%s) completed", migration.Version, migration.Name)
	return nil
}

// createUsersTable creates the users table
func createUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		email TEXT,
		bio TEXT,
		avatar_color TEXT DEFAULT '#8774e1',
		status TEXT DEFAULT 'offline' CHECK(status IN ('online', 'offline', 'away')),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	// Create trigger for updated_at
	trigger := `
	CREATE TRIGGER IF NOT EXISTS update_users_timestamp 
	AFTER UPDATE ON users
	FOR EACH ROW
	BEGIN
		UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
	END`
	_, err = db.Exec(trigger)
	return err
}

// createFriendshipsTable creates the friendships table
func createFriendshipsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS friendships (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		friend_id INTEGER NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'accepted', 'rejected')),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, friend_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (friend_id) REFERENCES users(id) ON DELETE CASCADE
	)`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_friendships_user_id ON friendships(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_friendships_friend_id ON friendships(friend_id)",
	}

	for _, idx := range indexes {
		if _, err := db.Exec(idx); err != nil {
			return err
		}
	}

	// Create trigger for updated_at
	trigger := `
	CREATE TRIGGER IF NOT EXISTS update_friendships_timestamp 
	AFTER UPDATE ON friendships
	FOR EACH ROW
	BEGIN
		UPDATE friendships SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
	END`
	_, err = db.Exec(trigger)
	return err
}
