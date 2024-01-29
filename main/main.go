package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"todo-api-go/api"
	"todo-api-go/entities"
	"todo-api-go/oidc"
	"todo-api-go/persistence"

	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

var (
	// flags to be provided for running the example server
	domain = flag.String("domain", "", "your ZITADEL instance domain (in the form: <instance>.zitadel.cloud or <yourdomain>)")
	key    = flag.String("key", "", "path to your key.json")
	port   = flag.String("port", "8089", "port to run the server on (default is 8089)")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	// Initiate the authorization by providing a zitadel configuration and a verifier.
	// This example will use OAuth2 Introspection for this, therefore you will also need to provide the downloaded api key.json
	authZ, err := authorization.New(ctx, zitadel.New(*domain, zitadel.WithInsecure("8088")), oauth.DefaultAuthorization(*key))
	if err != nil {
		slog.Error("zitadel sdk could not initialize", "error", err)
		os.Exit(1)
	}

	// Initialize the HTTP middleware by providing the authorization
	mw := oidc.New(authZ)

	router := gin.Default()
	api.RegisterRoutes(router, createEntityManager(), mw)
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
