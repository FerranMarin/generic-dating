package models

import "gorm.io/gorm"

type Swipe struct {
	gorm.Model
	UserID       uint `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ID;references:ID"`
	SwipedUserID uint `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ID;references:ID"`
	Preference   bool
}
