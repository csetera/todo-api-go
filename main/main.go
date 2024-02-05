package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"todo-api-go/api"
	"todo-api-go/entities"
	"todo-api-go/oidc"
	"todo-api-go/persistence"
)

func main() {
	// Initialize the HTTP middleware for authorization
	authz, err := oidc.New()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	api.RegisterRoutes(router, createEntityManager(), authz)
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
	dialector, err := persistence.OpenDialectorFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	if strings.ToLower(os.Getenv("DB_AUTO_MIGRATE")) == "true" {
		err = db.AutoMigrate(&entities.ToDoItemEntity{})
		if err != nil {
			log.Fatal(err)
		}
	}

	return persistence.New(db)
}
