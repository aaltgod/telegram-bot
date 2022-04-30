package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	User struct {
		UUID     uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
		Name     string
		ID       int64
		IsAdmin  bool
		Requests []Request `gorm:"foreignKey:UserUUID"`
	}

	CreateUser struct {
		Name    string
		ID      int64
		IsAdmin bool
	}

	UpdateUser struct {
		IsAdmin bool
	}

	Request struct {
		gorm.Model
		IP       string
		Response string
		UserUUID uuid.UUID
	}

	DeleteRequest struct {
		IP string
	}
)
