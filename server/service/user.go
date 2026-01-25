package service

import (
	"context"
	"errors"
	"server/db"
	"server/models"
	"server/utils"

	"gorm.io/gorm"
)

var UserServiceApp = new(UserService)

type UserService struct{}

func (this *UserService) CreateNewUser(ctx context.Context, username string, password string, email string) error {

	_, err := gorm.G[models.User](db.PgSqlDB).
		Select("id").
		Where("email = ?", email).
		First(ctx)

	switch err {
	case nil:
		return errors.New("this email has been already used")
	case gorm.ErrRecordNotFound: // user not found
		break
	default:
		utils.Logger.Error(err)
		return errors.New("Unknown error")
	}

	db_user := &models.User{
		Username: username,
		Password: utils.BcryptHash(password),
		Email:    email,
	}
	if err := gorm.G[models.User](db.PgSqlDB).Create(ctx, db_user); err != nil {
		utils.Logger.Error(err)
		return errors.New("Unknown error")
	}
	return nil
}

func (this *UserService) GetUserInfoByEmail(ctx context.Context, email string) (*models.User, error) {

	db_user, err := gorm.G[models.User](db.PgSqlDB).
		Where("email = ?", email).
		First(ctx)
	switch err {
	case nil:
		return &db_user, nil
	case gorm.ErrRecordNotFound:
		utils.Logger.Error(err)
		return nil, errors.New("User not found")
	default:
		utils.Logger.Error(err)
		return nil, errors.New("Unknown error")
	}

}

func (this *UserService) GetUserInfoByID(ctx context.Context, id uint) (*models.User, error) {
	db_user, err := gorm.G[models.User](db.PgSqlDB).
		Where("id = ?", id).
		First(ctx)
	switch err {
	case nil:
		return &db_user, nil
	case gorm.ErrRecordNotFound:
		utils.Logger.Error(err)
		return nil, errors.New("User not found")
	default:
		utils.Logger.Error(err)
		return nil, errors.New("Unknown error")
	}
}

func (this *UserService) GetUserInfoByUsername(ctx context.Context, username string) (*models.User, error) {
	db_user, err := gorm.G[models.User](db.PgSqlDB).
		Where("username = ?", username).
		First(ctx)
	switch err {
	case nil:
		return &db_user, nil
	case gorm.ErrRecordNotFound:
		utils.Logger.Error(err)
		return nil, errors.New("User not found")
	default:
		utils.Logger.Error(err)
		return nil, errors.New("Unknown error")
	}
}

func (this *UserService) ResetUserPassword(ctx context.Context, userID uint, newPassword string) error {

	hashed_password := utils.BcryptHash(newPassword)

	row, err := gorm.G[models.User](db.PgSqlDB).
		Where("id = ?", userID).
		Update(ctx, "password", hashed_password)
	utils.Logger.Debugf("rowsAffected :%d", row)
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return errors.New("User not found")
	default:
		utils.Logger.Error(err)
		return errors.New("Unknown error")
	}
}

func (this *UserService) UpdateUserInfo(ctx context.Context, userID uint, newUsername string, newEmail string) error {

	new_user_info := models.User{Username: newUsername, Email: newEmail}

	_, err := gorm.G[models.User](db.PgSqlDB).
		Where("id = ?", userID).
		Updates(ctx, new_user_info)
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return errors.New("User not found")
	default:
		utils.Logger.Error(err)
		return errors.New("Unknown error")
	}
}
