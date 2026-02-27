package service

import (
	"context"
	"server/db"
	"server/models"
	"server/utils"

	"gorm.io/gorm"
)

var UserServiceApp = new(UserService)

type UserService struct{}

func (this *UserService) CreateNewUser(ctx context.Context, username string, password string, email string) (err error) {
	dbUser := &models.User{
		Username: username,
		Password: utils.BcryptHash(password),
		Email:    email,
	}
	return gorm.G[models.User](db.PgSqlDB).Create(ctx, dbUser)
}

func (this *UserService) GetUserInfoByEmail(ctx context.Context, email string) (*models.User, error) {

	db_user, err := gorm.G[models.User](db.PgSqlDB).
		Where("email = ?", email).
		First(ctx)
	return &db_user, err
}

func (this *UserService) GetUserInfoByID(ctx context.Context, id uint) (*models.User, error) {
	db_user, err := gorm.G[models.User](db.PgSqlDB).
		Where("id = ?", id).
		First(ctx)
	return &db_user, err
}

func (this *UserService) GetUserInfoByUsername(ctx context.Context, username string) (*models.User, error) {
	db_user, err := gorm.G[models.User](db.PgSqlDB).
		Where("username = ?", username).
		First(ctx)
	return &db_user, err
}

func (this *UserService) ResetUserPassword(ctx context.Context, userID uint, newPassword string) error {
	hashed_password := utils.BcryptHash(newPassword)
	_, err := gorm.G[models.User](db.PgSqlDB).
		Where("id = ?", userID).
		Update(ctx, "password", hashed_password)
	return err
}

func (this *UserService) UpdateUserInfo(ctx context.Context, userID uint, newUsername string, newEmail string) error {
	new_user_info := models.User{Username: newUsername, Email: newEmail}
	_, err := gorm.G[models.User](db.PgSqlDB).
		Where("id = ?", userID).
		Updates(ctx, new_user_info)
	return err
}

func (this *UserService) CheckUserExistsByEmail(ctx context.Context, ownerID uint, email string) (bool, error) {
	cnt, err := gorm.G[models.User](db.PgSqlDB).
		Where("email = ? AND id != ?", email, ownerID).
		Count(ctx, "*")
	return cnt > 0, err
}
