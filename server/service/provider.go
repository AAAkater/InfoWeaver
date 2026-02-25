package service

import (
	"context"
	"server/db"
	"server/models"
	"server/utils"

	"gorm.io/gorm"
)

var ProviderServiceApp = new(ProviderService)

type ProviderService struct{}

func (this *ProviderService) CreateProvider(ctx context.Context, ownerID uint, name string, baseURL string, apiKey string, mode string) error {

	dbProvider := &models.Provider{
		OwnerID: ownerID,
		Name:    name,
		BaseURL: baseURL,
		APIKey:  utils.BcryptHash(apiKey),
		Mode:    mode,
	}
	return gorm.G[models.Provider](db.PgSqlDB).Create(ctx, dbProvider)
}

func (this *ProviderService) UpdateProvider(ctx context.Context, providerID uint, ownerID uint, name string, baseURL string, apiKey string, mode string) error {
	newProvider := models.Provider{
		Name:    name,
		BaseURL: baseURL,
		APIKey:  apiKey,
		Mode:    mode,
	}
	rows, err := gorm.G[models.Provider](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", providerID, ownerID).
		Updates(ctx, newProvider)
	if rows == 0 {
		return ErrNotFound
	}
	return err
}
func (this *ProviderService) GetProviderByID(ctx context.Context, providerID uint, ownerID uint) (dbProvider *models.ProviderInfo, e error) {
	result := db.PgSqlDB.Model(&models.Provider{}).
		Where("id = ? AND owner_id = ?", providerID, ownerID).
		First(dbProvider)
	return dbProvider, result.Error
}

func (this *ProviderService) GetProviderByName(ctx context.Context, name string, ownerID uint) (dbProvider *models.ProviderInfo, err error) {
	result := db.PgSqlDB.Model(&models.Provider{}).
		Where("name = ? AND owner_id = ?", name, ownerID).
		First(dbProvider)
	return dbProvider, result.Error
}

func (this *ProviderService) GetAllProviders(ctx context.Context, ownerID uint) (cows int64, dbProviders []models.ProviderInfo, err error) {

	result := db.PgSqlDB.Model(&models.Provider{}).
		Where("owner_id = ?", ownerID).
		Find(&dbProviders)

	return result.RowsAffected, dbProviders, result.Error
}

func (this *ProviderService) DeleteProvider(ctx context.Context, providerID uint, ownerID uint) error {
	_, err := gorm.G[models.Provider](db.PgSqlDB).
		Where("id = ? AND owner_id = ?", providerID, ownerID).
		Delete(ctx)
	return err
}
