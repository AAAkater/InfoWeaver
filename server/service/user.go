package service

import (
	"errors"
	"server/db"
	"server/models"
	"server/utils"

	"gorm.io/gorm"
)

var UserServiceApp = new(UserService)

type UserService struct{}

func (this *UserService) Register(new_user models.User) error {
	var user models.User
	res := db.PgSqlDB.Where("username = ?", new_user.Username).First(&user)

	if res.Error != gorm.ErrRecordNotFound {
		return errors.New("this username has been already used")
	}

	new_user.Password = utils.BcryptHash(new_user.Password)
	err := db.PgSqlDB.Create(&new_user).Error

	return err
}
