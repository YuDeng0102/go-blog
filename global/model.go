package global

import (
	"gorm.io/gorm"
	"time"
)

type MODEL struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
