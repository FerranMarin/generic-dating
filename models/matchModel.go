package models

import "gorm.io/gorm"

type Match struct {
	gorm.Model
	UserID1 uint `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ID;references:ID"`
	UserID2 uint `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:ID;references:ID"`
}
