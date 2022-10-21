package models

import "gorm.io/gorm"

// Product model info
// @Description Product Information

type User struct {
	gorm.Model
	Id       uint     `form:"id" json:"id" validate:"required"`
	Name     string  `form:"name" json:"name" validate:"required"`
	Email string     `form:"email" json:"email" validate:"required"`
	Username    string `form:"username" json:"username" validate:"required"`
	Password    string `form:"password" json:"password" validate:"required"`
}

func Register(db *gorm.DB, register *User) (err error) {
	err = db.Create(register).Error
	if err != nil {
		return err
	}
	return nil
}

func ReadUser(db *gorm.DB, users *[]User) (err error) {
	err = db.Find(users).Error
	if err != nil {
		return err
	}
	return nil
}

func ReadOneUser(db *gorm.DB, login *User, username string) (err error) {
	err = db.Where(&User{
		Username: username,
	}).First(login).Error
	if err != nil {
		return err
	}
	return nil
}


