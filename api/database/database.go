package database

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main/models"
	"time"
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

func (receiver Song) Update(dto models.SongDto) error {
	_, err := receiver.GetById(dto.Id)
	if err != nil {
		return fmt.Errorf("song with id = '%s' was not found", dto.Id)
	}
	return receiver.database.Save(&dto).Error
}

func (receiver Song) Delete(id string) error {
	_, err := receiver.GetById(id)
	if err != nil {
		return fmt.Errorf("song with id = '%s' was not found", id)
	}

	return receiver.database.Delete(&models.SongDto{}, "id = ?", id).Error
}

func (receiver Song) Ping() error {
	db, err := receiver.database.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	return db.PingContext(ctx)
}
