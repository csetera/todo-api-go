package testsupport

import (
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"

	"todo-api-go/entities"
	"todo-api-go/persistence"

	"fmt"
	"testing"
)

// CreateTestManager creates a ToDoEntityManager instance for testing.
//
// This function takes a testing.T instance as a parameter and returns a *persistence.ToDoEntityManager.
func CreateTestManager(t *testing.T) *persistence.ToDoEntityManager {
	dsn := "file::memory:?cache=shared"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		t.Fatal("failed to connect database")
	}

	if db == nil {
		t.Fatal("db is nil")
	}

	// Set up the schema and load some test data
	err = db.AutoMigrate(&entities.ToDoItemEntity{})
	if err != nil {
		t.Fatal(err)
	}
	loadTestData(t, db)

	// Wrap the database connection in a ToDoEntityManager
	mgr := persistence.New(db)
	if mgr == nil {
		t.Fatal("mgr is nil")
	}

	// Defer the closing of the database connection until
	// the end of the test
	t.Cleanup(func() {
		mgr.Close()
	})

	return mgr
}

// loadTestData loads test data into the database using
// standard SQL.
//
// Parameters:
// - t: The testing.T object for running tests.
// - db: The GORM database connection.
// Return type: None.
func loadTestData(t *testing.T, db *gorm.DB) {
	for i := 0; i < 10; i++ {
		description := fmt.Sprintf("Todo Item %d", i)
		db.Exec("INSERT INTO to_do_item_entities (id, description, completed, due_date) VALUES (?, ?, ?, ?)",
			(i + 1), description, false, ParseTestDate("2025-01-01"))
	}
}
