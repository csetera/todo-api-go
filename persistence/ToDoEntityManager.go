package persistence

import (
	"context"

	"gorm.io/gorm"

	"todo-api-go/entities"
)

type ToDoEntityManager struct {
	orm *gorm.DB
}

// Close closes the ToDoEntityManager and associated database connection.
//
// The Close function does not take any parameters.
// It does not return any values.
func (mgr *ToDoEntityManager) Close() {
	db, err := mgr.orm.DB()
	if err != nil {
		return
	}

	db.Close()
}

// Create creates a ToDoItemEntity in the database.
//
// It takes a pointer to a ToDoItemEntity as a parameter.
// It returns an error if there was an issue creating the entity.
func (mgr *ToDoEntityManager) Create(item *entities.ToDoItemEntity) error {
	return mgr.orm.Create(item).Error
}

// Delete a ToDoItemEntity from the database by its ID.
//
// Parameters:
// - id: the ID of the ToDoItemEntity to be deleted.
//
// Returns:
// - error: an error if the deletion operation fails.
func (mgr *ToDoEntityManager) Delete(id uint) error {
	return mgr.orm.Delete(&entities.ToDoItemEntity{}, id).Error
}

// FindAll retrieves all ToDoItemEntity objects from the database based on the provided paging configuration.
//
// The function accepts optional PagingConfigurator arguments to configure the pagination of the results.
// It returns a slice of ToDoItemEntity objects and an error if any occurred.
func (mgr *ToDoEntityManager) FindAll(configurators ...PagingConfigurator) ([]entities.ToDoItemEntity, int64, error) {
	var count int64

	err := mgr.orm.Model(&entities.ToDoItemEntity{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	var items []entities.ToDoItemEntity
	err = mgr.orm.Scopes(Paginate(configurators...)).Order("id asc").Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, count, nil
}

// FineOne returns a ToDoItemEntity and an error.
//
// It takes an integer `id` as a parameter.
// The function retrieves a ToDoItemEntity with the given `id` from the database using the `mgr.orm.First` method.
// If the retrieval is successful, it returns a pointer to the retrieved ToDoItemEntity and a `nil` error.
// If an error occurs during the retrieval, it returns a `nil` ToDoItemEntity and the error encountered.
func (mgr *ToDoEntityManager) FineOne(id int) (*entities.ToDoItemEntity, error) {
	var item entities.ToDoItemEntity

	err := mgr.orm.First(&item, id).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// func (mgr *ToDoEntityManager) ORM() *gorm.DB {
// 	return mgr.orm
// }

func (mgr *ToDoEntityManager) WithContext(ctx context.Context) *ToDoEntityManager {
	return &ToDoEntityManager{orm: mgr.orm.WithContext(ctx)}
}

// New creates a new instance of ToDoEntityManager.
//
// Parameters:
// - orm: A pointer to a gorm.DB object representing the underlying GORM ORM instance.
//
// Returns:
// - A pointer to a ToDoEntityManager object.
func New(orm *gorm.DB) *ToDoEntityManager {
	return &ToDoEntityManager{orm: orm}
}
