package api_test

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"todo-api-go/api"
	"todo-api-go/entities"
	"todo-api-go/persistence"
	"todo-api-go/testsupport"

	"net/http"
	"net/http/httptest"

	"bytes"
	"testing"
)

type MockAuthorizer struct {
}

func (mock *MockAuthorizer) RequiresRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	mgr := testsupport.CreateTestManager(t)
	assert.NotNilf(mgr, "manager should not be nil")

	item := &entities.ToDoItemEntity{
		Description: "New Todo Item",
		Completed:   false,
		DueDate:     testsupport.ParseTestDate("2024-01-01"),
	}
	marshalled, err := json.Marshal(item)
	assert.Nilf(err, "error should be nil")

	req, _ := http.NewRequest("POST", "/api/todo", bytes.NewBuffer(marshalled))
	recorder := makeRequest(mgr, req)
	assert.Equalf(201, recorder.Code, "Expected successful response")

}

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	mgr := testsupport.CreateTestManager(t)
	assert.NotNilf(mgr, "manager should not be nil")

	req, _ := http.NewRequest("DELETE", "/api/todo/5", nil)
	recorder := makeRequest(mgr, req)
	assert.Equalf(200, recorder.Code, "Expected successful response")

	req, _ = http.NewRequest("GET", "/api/todo", nil)
	recorder = makeRequest(mgr, req)
	assert.Equalf(200, recorder.Code, "Expected successful response")

	var response api.FindResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Nilf(err, "error should be nil")
	assert.Equalf(9, int(response.Meta.Total), "total length should be 9")
	assert.Equalf(9, len(response.Data), "length should be 9")
	assert.ElementsMatchf([]uint{1, 2, 3, 4, 6, 7, 8, 9, 10}, testsupport.CollectIds(response.Data), "IDs should match")
}

func TestGetAll(t *testing.T) {
	assert := assert.New(t)

	mgr := testsupport.CreateTestManager(t)
	assert.NotNilf(mgr, "manager should not be nil")

	req, _ := http.NewRequest("GET", "/api/todo", nil)
	recorder := makeRequest(mgr, req)
	assert.Equalf(200, recorder.Code, "Expected successful response")

	var response api.FindResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Nilf(err, "error should be nil")
	assert.Equalf(10, int(response.Meta.Total), "total length should be 10")
	assert.Equalf(10, len(response.Data), "length should be 10")
	assert.ElementsMatchf([]uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, testsupport.CollectIds(response.Data), "IDs should match")

	req, _ = http.NewRequest("GET", "/api/todo?limit=5&offset=1", nil)
	recorder = makeRequest(mgr, req)
	assert.Equalf(200, recorder.Code, "Expected successful response")

	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Nilf(err, "error should be nil")
	assert.Equalf(10, int(response.Meta.Total), "total length should be 10")
	assert.Equalf(5, len(response.Data), "length should be 5")
	assert.ElementsMatchf([]uint{2, 3, 4, 5, 6}, testsupport.CollectIds(response.Data), "IDs should match")
}

func TestGetByID(t *testing.T) {
	assert := assert.New(t)

	mgr := testsupport.CreateTestManager(t)
	assert.NotNilf(mgr, "manager should not be nil")

	req, _ := http.NewRequest("GET", "/api/todo/1", nil)
	recorder := makeRequest(mgr, req)
	assert.Equalf(200, recorder.Code, "Expected successful response")

	var item entities.ToDoItemEntity
	err := json.Unmarshal(recorder.Body.Bytes(), &item)
	assert.Nilf(err, "error should be nil")
	assert.Equalf(1, int(item.ID), "IDs should match")
	assert.Equalf("Todo Item 0", item.Description, "descriptions should match")
}

func makeRequest(mgr *persistence.ToDoEntityManager, request *http.Request) *httptest.ResponseRecorder {
	var mock MockAuthorizer

	router := gin.Default()
	api.RegisterRoutes(router, mgr, &mock)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	return recorder
}
