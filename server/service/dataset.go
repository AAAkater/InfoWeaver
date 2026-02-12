package service

import (
	"context"
	"server/db"
	"server/models"

	"gorm.io/gorm"
)

var DatasetServiceApp = new(DatasetService)

type DatasetService struct{}

func (this *DatasetService) CreateNewDataset(ctx context.Context, icon string, datasetName string, description string, ownerID uint) error {

	dbDataset := models.Dataset{
		Name:        datasetName,
		Icon:        icon,
		Description: description,
		OwnerID:     ownerID,
	}
	return gorm.G[models.Dataset](db.PgSqlDB).Create(ctx, &dbDataset)

}

func (this *DatasetService) GetDatasetInfoByID(ctx context.Context, id uint, ownerID uint) (dbDataset *models.DatasetInfo, err error) {

	result := db.PgSqlDB.Model(&models.Dataset{}).
		Where("id = ? AND owner_id = ?", id, ownerID).
		First(&dbDataset)

	return dbDataset, result.Error
}

// ListDatasetsByOwnerID retrieves datasets for a specific user with pagination
// page: page number (1-indexed), pageSize: number of items per page
func (this *DatasetService) ListDatasetsByOwnerID(ctx context.Context, ownerID uint) (total int64, datasets []models.DatasetInfo, e error) {

	result := db.PgSqlDB.Model(&models.Dataset{}).
		Where("owner_id = ?", ownerID).
		Find(&datasets)

	return result.RowsAffected, datasets, result.Error
}

func (this *DatasetService) UpdateDataset(ctx context.Context, id uint, ownerID uint, icon string, name string, description string) error {

	newDatasetInfo := models.Dataset{
		Icon:        icon,
		Name:        name,
		Description: description,
	}

	rowsAffected, err := gorm.G[models.Dataset](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", id, ownerID).
		Updates(ctx, newDatasetInfo)
	// id not found
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return err
}

func (this *DatasetService) DeleteDataset(ctx context.Context, id uint, ownerID uint) error {
	rowsAffected, err := gorm.G[models.Dataset](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", id, ownerID).
		Delete(ctx)
	// id not found
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return err
}

// GetDatasetInfoByName retrieves a dataset by its name and owner ID using fuzzy matching (contains)
func (this *DatasetService) ListDatasetsByName(ctx context.Context, ownerID uint, name string) (rowsAffected int64, dbDataset []models.DatasetInfo, err error) {
	result := db.PgSqlDB.Model(&models.Dataset{}).
		Where("name LIKE ? AND owner_id = ?", "%"+name+"%", ownerID).
		Find(&dbDataset)
	return result.RowsAffected, dbDataset, result.Error
}
