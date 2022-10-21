package models

import "gorm.io/gorm"

// Product model info
// @Description Product Information
type Cart struct {
	gorm.Model
	Id       uint     `json:"id" validate:"required"`
	ProductId     int  `form:"productid" json:"productid" validate:"required"`
	UserId     int  `form:"userid" json:"userid" validate:"required"`
	Quantity int     `form:"quantity" json:"quantity" validate:"required"`
	Total    float64 `form:"total" json:"total" validate:"required"`
	Product Product `gorm:"foreignkey:ProductId;references:Id"`
	User User `gorm:"foreignkey:UserId;references:Id"`
}

func ViewCart(db *gorm.DB, cart *[]Cart) (err error) {
	err = db.Preload("User").Preload("Product").Find(cart).Error
	if err != nil {
		return err
	}
	return nil
}

func AddtoCart(db *gorm.DB, data *Cart) (err error) {
	err = db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteCart(db *gorm.DB, data *Cart, id int) (err error) {
	err = db.Where("id=?", id).Delete(data).Error
	if err != nil {
		return err
	}
	return nil
}