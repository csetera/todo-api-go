package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"todo-api-go/entities"
	"todo-api-go/oidc"
	"todo-api-go/persistence"
)

type ListMetadata struct {
	Total int64
}

type FindResponse struct {
	Meta ListMetadata
	Data []entities.ToDoItemEntity
}

// RegisterRoutes registers the ToDo API routes for the Gin engine.
//
// gin: The Gin engine to register the routes with.
// mgr: The ToDo entity manager.
// Returns the registered Gin engine.
func RegisterRoutes(gin *gin.Engine, mgr *persistence.ToDoEntityManager, mw *oidc.OIDCMiddleware) *gin.Engine {
	gin.DELETE("/api/todo/:id", mw.RequiresRole("delete"), deleteToDoItemHandler(mgr))
	gin.GET("/api/todo", mw.RequiresRole("retrieve"), getAllToDoItemsHandler(mgr))
	gin.GET("/api/todo/:id", mw.RequiresRole("retreive"), getToDoByIdHandler(mgr))
	gin.POST("/api/todo", mw.RequiresRole("create"), createToDoItemHandler(mgr))

	return gin
}

// createToDoItemHandler creates a HandlerFunc function for creating a ToDoItemEntity.
//
// It takes a manager of type *persistence.ToDoEntityManager as a parameter.
// The function returns a gin.HandlerFunc.
func createToDoItemHandler(manager *persistence.ToDoEntityManager) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		var item entities.ToDoItemEntity
		err := c.BindJSON(&item)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = manager.Create(&item)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusCreated, item)
	})
}

// deleteToDoItemHandler creates a HandlerFunc function for deleting a ToDoItemEntity
// by identifier
//
// It takes a manager of type *persistence.ToDoEntityManager as a parameter.
// The function returns a gin.HandlerFunc.
func deleteToDoItemHandler(manager *persistence.ToDoEntityManager) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		id_value := c.Param("id")
		id, err := strconv.Atoi(id_value)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = manager.Delete(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusOK)
	})
}

// getAllToDoItemsHandler creates a HandlerFunc function for getting all ToDoItemEntity's with
// pagination.
//
// It takes a manager of type *persistence.ToDoEntityManager as a parameter.
// The function returns a gin.HandlerFunc.
func getAllToDoItemsHandler(manager *persistence.ToDoEntityManager) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		items, total, err := manager.FindAll(getPagingConfigurator(c))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := FindResponse{
			Meta: ListMetadata{Total: total},
			Data: items,
		}
		c.IndentedJSON(http.StatusOK, response)
	})
}

// getPagingConfigurator generates a function that configures the paging options for a given Gin request context.
//
// It takes a Gin context object as a parameter and returns a function that takes a pointer to a PagingOptions object.
// The PagingOptions object is modified based on the "limit" and "offset" query parameters from the Gin context.
func getPagingConfigurator(c *gin.Context) func(*persistence.PagingOptions) {
	return func(options *persistence.PagingOptions) {
		limit, err := strconv.Atoi(c.Query("limit"))
		if err == nil {
			options.Limit = limit
		}

		offset, err := strconv.Atoi(c.Query("offset"))
		if err == nil {
			options.Offset = offset
		}
	}
}

// getToDoByIdHandler returns a Gin handler function that retrieves a to-do item by its ID.
//
// The function takes a pointer to a `persistence.ToDoEntityManager` as its parameter.
// It expects a Gin context object `c`, which contains the ID of the to-do item as a URL parameter.
// The function first parses the ID from the URL parameter and handles any parsing errors.
// It then calls the `FineOne` method of the `manager` to retrieve the to-do item with the given ID.
// If any error occurs during the retrieval process, it returns a JSON response with the corresponding error message.
// Otherwise, it returns a JSON response with the retrieved to-do item.
func getToDoByIdHandler(manager *persistence.ToDoEntityManager) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		id_value := c.Param("id")
		id, err := strconv.Atoi(id_value)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := manager.FineOne(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, todo)
	})
}
