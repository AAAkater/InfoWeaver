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

func (this *DatasetService) GetDatasetByID(ctx context.Context, id uint, ownerID uint) (dbDataset *models.Dataset, err error) {

	*dbDataset, err = gorm.G[models.Dataset](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", id, ownerID).
		First(ctx)
	switch err {
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

func (this *DatasetService) ListDatasetsByOwner(ctx context.Context, ownerID uint) (datasets []models.Dataset, err error) {

	datasets, err = gorm.G[models.Dataset](db.PgSqlDB).
		Where("owner_id = ?", ownerID).
		Find(ctx)
	if err != nil {
		utils.Logger.Errorf("Failed to list datasets for owner ID %d: %v", ownerID, err)
		return nil, err
	}
	return datasets, nil
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
