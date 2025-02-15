//go:build mage
// +build mage

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "words.db"

// SeedGroup represents a group in the seed file
type SeedGroup struct {
	Name string `json:"name"`
}

// SeedWord represents a word in the seed file
type SeedWord struct {
	Japanese string `json:"japanese"`
	Romaji   string `json:"romaji"`
	English  string `json:"english"`
	Parts    string `json:"parts"`
}

// SeedFile represents the structure of a word group seed file
type SeedFile struct {
	Group SeedGroup  `json:"group"`
	Words []SeedWord `json:"words"`
}

// ActivitySeedFile represents the structure of the activities seed file
type ActivitySeedFile struct {
	Activities []struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		ThumbnailURL string `json:"thumbnail_url"`
	} `json:"activities"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Available commands:")
		fmt.Println("  initdb   - Initialize the database")
		fmt.Println("  migrate  - Run database migrations")
		fmt.Println("  seed     - Seed the database with initial data")
		fmt.Println("  clean    - Remove the database")
		fmt.Println("  reset    - Reset the database (clean + init + migrate)")
		return
	}

	var err error
	switch os.Args[1] {
	case "initdb":
		err = InitDB()
	case "migrate":
		err = Migrate()
	case "seed":
		err = Seed()
	case "clean":
		err = Clean()
	case "reset":
		err = Reset()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// InitDB initializes the SQLite database
func InitDB() error {
	fmt.Println("Initializing database...")

	if _, err := os.Stat(dbName); err == nil {
		fmt.Printf("Database %s already exists\n", dbName)
		return nil
	}

	file, err := os.Create(dbName)
	if err != nil {
		return fmt.Errorf("error creating database file: %v", err)
	}
	file.Close()

	fmt.Printf("Created database %s\n", dbName)
	return nil
}

// Migrate runs all migration files in order
func Migrate() error {
	fmt.Println("Running migrations...")

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Get all migration files
	files, err := filepath.Glob("db/migrations/*.sql")
	if err != nil {
		return fmt.Errorf("error finding migration files: %v", err)
	}

	// Sort migration files by name
	sort.Strings(files)

	for _, file := range files {
		fmt.Printf("Applying migration %s...\n", filepath.Base(file))

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading migration file %s: %v", file, err)
		}

		// Split the file into separate statements
		statements := strings.Split(string(content), ";")

		// Execute each statement
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			_, err = db.Exec(stmt)
			if err != nil {
				return fmt.Errorf("error executing migration %s: %v", file, err)
			}
		}
	}

	fmt.Println("Migrations completed successfully")
	return nil
}

// Seed imports data from JSON files in the seeds directory
func Seed() error {
	fmt.Println("Seeding database...")

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Get all seed files
	files, err := filepath.Glob("db/seeds/*.json")
	if err != nil {
		return fmt.Errorf("error finding seed files: %v", err)
	}

	// Process activities first
	for _, file := range files {
		if strings.Contains(file, "activities") {
			if err := seedActivities(db, file); err != nil {
				return err
			}
		}
	}

	// Then process word groups
	for _, file := range files {
		if !strings.Contains(file, "activities") {
			if err := seedWordGroup(db, file); err != nil {
				return err
			}
		}
	}

	fmt.Println("Seeding completed successfully")
	return nil
}

func seedActivities(db *sql.DB, file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error reading seed file %s: %v", file, err)
	}

	var seedFile ActivitySeedFile
	if err := json.Unmarshal(content, &seedFile); err != nil {
		return fmt.Errorf("error parsing seed file %s: %v", file, err)
	}

	for _, activity := range seedFile.Activities {
		_, err := db.Exec(`
			INSERT INTO study_activities (id, study_session_id, group_id)
			VALUES (?, 1, 1)
		`, activity.ID)
		if err != nil {
			return fmt.Errorf("error inserting activity: %v", err)
		}
	}

	return nil
}

func seedWordGroup(db *sql.DB, file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("error reading seed file %s: %v", file, err)
	}

	var seedFile SeedFile
	if err := json.Unmarshal(content, &seedFile); err != nil {
		return fmt.Errorf("error parsing seed file %s: %v", file, err)
	}

	// Insert group
	result, err := db.Exec(`
		INSERT INTO groups (name)
		VALUES (?)
	`, seedFile.Group.Name)
	if err != nil {
		return fmt.Errorf("error inserting group: %v", err)
	}

	groupID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting group ID: %v", err)
	}

	// Insert words and create word-group associations
	for _, word := range seedFile.Words {
		result, err := db.Exec(`
			INSERT INTO words (japanese, romaji, english, parts)
			VALUES (?, ?, ?, ?)
		`, word.Japanese, word.Romaji, word.English, word.Parts)
		if err != nil {
			return fmt.Errorf("error inserting word: %v", err)
		}

		wordID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting word ID: %v", err)
		}

		_, err = db.Exec(`
			INSERT INTO words_groups (word_id, group_id)
			VALUES (?, ?)
		`, wordID, groupID)
		if err != nil {
			return fmt.Errorf("error inserting word-group association: %v", err)
		}
	}

	fmt.Printf("Seeded group '%s' with %d words\n", seedFile.Group.Name, len(seedFile.Words))
	return nil
}

// Clean removes the database file
func Clean() error {
	fmt.Printf("Removing database %s...\n", dbName)

	err := os.Remove(dbName)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error removing database: %v", err)
	}

	fmt.Println("Database cleaned successfully")
	return nil
}

// Reset cleans and reinitializes the database
func Reset() error {
	if err := Clean(); err != nil {
		return err
	}
	if err := InitDB(); err != nil {
		return err
	}
	if err := Migrate(); err != nil {
		return err
	}
	return nil
}
