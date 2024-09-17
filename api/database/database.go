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

func (receiver Song) GetById(id string) (*models.SongDto, error) {
	song := models.SongDto{}
	err := receiver.database.First(&song, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func (receiver Song) GetAll() (*[]models.SongDto, error) {
	var songs []models.SongDto
	err := receiver.database.Find(&songs).Error
	if err != nil {
		return nil, err
	}
	return &songs, err
}
