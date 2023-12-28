package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"todo-api-go/api"
	"todo-api-go/entities"
	"todo-api-go/persistence"

	"os"
	"strings"
)

func main() {
	router := gin.Default()
	api.RegisterRoutes(router, createEntityManager())
	router.Run(":8080")
}

// createEntityManager creates and returns a configured instance of ToDoEntityManager.
//
// It opens a dialector from the environment variables and uses it to open a new gorm DB connection.
// If the connection is successful, it checks if the DB_AUTO_MIGRATE environment variable is set to "true".
// If it is, it auto migrates the entities.ToDoItemEntity table.
// Finally, it returns a new instance of ToDoEntityManager using the created DB connection.
//
// Returns:
// *persistence.ToDoEntityManager - The newly created instance of ToDoEntityManager.
func createEntityManager() *persistence.ToDoEntityManager {
	dialector := persistence.OpenDialectorFromEnv()

	db, err := gorm.Open(dialector, &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		panic(err)
	}

	if strings.ToLower(os.Getenv("DB_AUTO_MIGRATE")) == "true" {
		err = db.AutoMigrate(&entities.ToDoItemEntity{})
		if err != nil {
			panic(err)
		}
	}

	return persistence.New(db)
}
