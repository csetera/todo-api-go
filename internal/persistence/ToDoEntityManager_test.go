package persistence_test

import (
	"github.com/stretchr/testify/assert"

	"todo-api-go/entities"
	"todo-api-go/persistence"
	"todo-api-go/testsupport"

	"testing"
)

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	mgr := testsupport.CreateTestManager(t)
	assert.NotNilf(mgr, "manager should not be nil")

	item := &entities.ToDoItemEntity{
		Description: "test",
		Completed:   false,
		DueDate:     testsupport.ParseTestDate("2024-01-01"),
	}

	err := mgr.Create(item)
	assert.Nilf(err, "error should be nil, not %s", err)
	assert.NotEqualf(0, int(item.ID), "ID should not be 0")
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	mgr := testsupport.CreateTestManager(t)
	assert.NotNilf(mgr, "manager should not be nil")

	items, total, err := mgr.FindAll()
	assert.Nilf(err, "error should be nil, not %s", err)
	assert.Equalf(int64(10), total, "total length should be 10")
	assert.Equalf(10, len(items), "length should be 10")
	assert.ElementsMatchf([]uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, testsupport.CollectIds(items), "IDs should match")

	err = mgr.Delete(5)
	assert.Nilf(err, "error should be nil, not %s", err)

	items, total, err = mgr.FindAll()
	assert.Nilf(err, "error should be nil, not %s", err)
	assert.Equalf(int64(9), total, "total length should be 9")
	assert.Equalf(9, len(items), "length should be 9")
	assert.ElementsMatchf([]uint{1, 2, 3, 4, 6, 7, 8, 9, 10}, testsupport.CollectIds(items), "IDs should match")
}

func TestFindAll(t *testing.T) {
	assert := assert.New(t)

	mgr := testsupport.CreateTestManager(t)
	assert.NotNilf(mgr, "manager should not be nil")

	items, total, err := mgr.FindAll()
	assert.Nilf(err, "error should be nil, not %s", err)
	assert.Equalf(int64(10), total, "total length should be 10")
	assert.Equalf(10, len(items), "length should be 10")
	assert.ElementsMatchf([]uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, testsupport.CollectIds(items), "IDs should match")

	items, total, err = mgr.FindAll(func(options *persistence.PagingOptions) {
		options.Limit = 5
		options.Offset = 1
	})

	assert.Nilf(err, "error should be nil, not %s", err)
	assert.Equalf(int64(10), total, "total length should be 10")
	assert.Equalf(5, len(items), "length should be 5")
	assert.ElementsMatchf([]uint{2, 3, 4, 5, 6}, testsupport.CollectIds(items), "IDs should match")
}

func TestFindById(t *testing.T) {
	assert := assert.New(t)

	mgr := testsupport.CreateTestManager(t)
	assert.NotNilf(mgr, "manager should not be nil")

	item := &entities.ToDoItemEntity{
		Description: "test",
		Completed:   false,
		DueDate:     testsupport.ParseTestDate("2024-01-01"),
	}

	err := mgr.Create(item)
	assert.Nilf(err, "error should be nil, not %s", err)

	founditem, finderr := mgr.FineOne(int(item.ID))
	assert.Nilf(finderr, "error should be nil, not %s", finderr)
	assert.NotNilf(founditem, "found item should not be nil")
	assert.Equalf(item.ID, founditem.ID, "found item should have same ID")
}
