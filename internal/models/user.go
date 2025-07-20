package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"type:text;unique;not null"`
	PasswordHash string    `gorm:"type:text;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	ShortLinks []ShortLink `gorm:"foreignKey:UserID"`
}

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}
