package service

import (
	"context"
	"errors"
	"server/db"
	"server/models"
	"server/utils"

	"gorm.io/gorm"
)

var DatasetServiceApp = new(DatasetService)

type DatasetService struct{}

func (this *DatasetService) CreateNewDataset(ctx context.Context, datasetName string, description string, ownerID uint) error {

	dbDataset := models.Dataset{
		Name:        datasetName,
		Description: description,
		OwnerID:     ownerID,
	}
	err := gorm.G[models.Dataset](db.PgSqlDB).Create(ctx, &dbDataset)
	switch err {
	case nil: // Dataset created successfully
		return nil
	case gorm.ErrDuplicatedKey:
		utils.Logger.Errorf("Dataset with name '%s' already exists for owner ID %d", datasetName, ownerID)
		return errors.New("dataset with the same name already exists")
	default:
		utils.Logger.Errorf("Failed to create dataset record in database: %v", err)
		return errors.New("unknown error")
	}

}

func (this *DatasetService) GetDatasetInfoByID(ctx context.Context, id uint, ownerID uint) (dbDataset *models.DatasetInfo, err error) {

	result := db.PgSqlDB.Model(&models.Dataset{}).
		Where("id = ? AND owner_id = ?", id, ownerID).
		Find(&dbDataset)

	switch result.Error {
	case nil: // Dataset found successfully
		return dbDataset, nil
	case gorm.ErrRecordNotFound:
		utils.Logger.Warnf("Dataset with ID %d not found for owner ID %d", id, ownerID)
		return nil, errors.New("dataset not found")
	default:
		utils.Logger.Errorf("Failed to get dataset by ID %d: %v", id, err)
		return nil, errors.New("unknown error")
	}
}

// ListDatasetsByOwnerID retrieves datasets for a specific user with pagination
// page: page number (1-indexed), pageSize: number of items per page
func (this *DatasetService) ListDatasetsByOwnerID(ctx context.Context, ownerID uint) (total int64, datasets []models.DatasetInfo, e error) {

	result := db.PgSqlDB.Model(&models.Dataset{}).
		Where("owner_id = ?", ownerID).
		Find(&datasets)

	switch result.Error {
	case nil:
		return result.RowsAffected, datasets, nil
	case gorm.ErrRecordNotFound:
		utils.Logger.Errorf("No datasets found for user %d: %v", ownerID, result.Error)
		return 0, datasets, nil
	default:
		utils.Logger.Errorf("Failed to get dataset list for user %d: %v", ownerID, result.Error)
		return 0, nil, errors.New("Unknown error occurred while retrieving dataset list")
	}
}

func (this *DatasetService) UpdateDataset(ctx context.Context, id uint, ownerID uint, name string, description string) error {

	newDatasetInfo := models.Dataset{Name: name, Description: description}

	_, err := gorm.G[models.Dataset](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", id, ownerID).
		Updates(ctx, newDatasetInfo)

	switch err {
	case nil:
		// Dataset updated successfully
		return nil
	case gorm.ErrRecordNotFound:
		utils.Logger.Infof("Dataset with ID %d not found for owner ID %d", id, ownerID)
		return errors.New("dataset not found")
	case gorm.ErrDuplicatedKey:
		utils.Logger.Errorf("Dataset with name '%s' already exists for owner ID %d", name, ownerID)
		return errors.New("dataset with the same name already exists")
	default:
		utils.Logger.Errorf("Failed to update dataset ID %d: %v", id, err)
		return errors.New("unknown error")
	}
}

func (this *DatasetService) DeleteDataset(ctx context.Context, id uint, ownerID uint) error {
	_, err := gorm.G[models.Dataset](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", id, ownerID).
		Delete(ctx)
	switch err {
	case nil: // Dataset deleted successfully
		return nil
	case gorm.ErrRecordNotFound:
		utils.Logger.Infof("Dataset with ID %d not found for owner ID %d", id, ownerID)
		return errors.New("dataset not found")
	default:
		utils.Logger.Errorf("Failed to delete dataset ID %d: %v", id, err)
		return errors.New("unknown error")
	}
}
