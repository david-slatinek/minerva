package database

import (
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main/models"
)

type Song struct {
	database *gorm.DB
}

func NewSong(dsn string) (*Song, error) {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	return &Song{database: db}, nil
}

func (receiver Song) Create(song models.Song) (models.SongDto, error) {
	sDto := models.SongDto{
		Id:   uuid.New().String(),
		Song: song,
	}
	return sDto, receiver.database.Create(&sDto).Error
}
