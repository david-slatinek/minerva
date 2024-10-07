package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/database"
	"main/logging"
	"main/models"
	"net/http"
)

type Song struct {
	db     *database.Song
	logger *logging.Logging
}

func NewSong(db *database.Song, logger *logging.Logging) *Song {
	return &Song{
		db:     db,
		logger: logger,
	}
}

// Create godoc
//
//	@Description	Create a new song.
//	@Summary		Create a new song
//	@Accept			json
//	@Produce		json
//	@Tags			song
//	@Param			requestBody	body		models.Song	true	"Song info"
//	@Success		201			{object}	models.SongDto
//	@Failure		400			{object}	models.Error
//	@Failure		500			{object}	models.Error
//	@Router			/songs [POST]
func (receiver Song) Create(c *gin.Context) {
	var req models.Song
	if err := c.ShouldBindJSON(&req); err != nil {
		receiver.logger.Set(c, logging.Error, http.StatusBadRequest, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}

	songDto, err := receiver.db.Create(req)
	if err != nil {
		receiver.logger.Set(c, logging.Error, http.StatusBadRequest, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}

	receiver.logger.Set(c, logging.Info, http.StatusOK, "")
	c.JSON(http.StatusOK, songDto)
}

// GetById godoc
//
//	@Description	Get a song by id.
//	@Summary		Get a song by id
//	@Produce		json
//	@Tags			song
//	@Param			id	path		string	true	"Song ID"
//	@Success		200	{object}	models.SongDto
//	@Failure		400	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/songs/{id} [GET]
func (receiver Song) GetById(c *gin.Context) {
	id := c.Param("id")

	songDto, err := receiver.db.GetById(id)
	if err != nil {
		message := fmt.Sprintf("song with id = '%s' was not found", id)
		receiver.logger.Set(c, logging.Error, http.StatusBadRequest, message)
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: message})
		return
	}

	receiver.logger.Set(c, logging.Info, http.StatusOK, "")
	c.JSON(http.StatusOK, songDto)
}

// GetAll godoc
//
//	@Description	Get all songs.
//	@Summary		Get all songs
//	@Produce		json
//	@Tags			song
//	@Success		200	{object}	[]models.SongDto	"An array of songs"
//	@Success		204	"No Content"
//	@Failure		500	{object}	models.Error
//	@Router			/songs [GET]
func (receiver Song) GetAll(c *gin.Context) {
	songs, err := receiver.db.GetAll()
	if err != nil {
		receiver.logger.Set(c, logging.Error, http.StatusInternalServerError, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
		return
	}

	if len(*songs) == 0 {
		receiver.logger.Set(c, logging.Info, http.StatusNoContent, "")
		c.Status(http.StatusNoContent)
		return
	}

	receiver.logger.Set(c, logging.Info, http.StatusOK, "")
	c.JSON(http.StatusOK, songs)
}

// Update godoc
//
//	@Description	Update song.
//	@Summary		Update song
//	@Accept			json
//	@Produce		json
//	@Tags			song
//	@Param			id			path		string		true	"Song ID"
//	@Param			requestBody	body		models.Song	true	"Updated song"
//	@Success		200			{object}	models.SongDto
//	@Failure		400			{object}	models.Error
//	@Failure		500			{object}	models.Error
//	@Router			/songs/{id} [PUT]
func (receiver Song) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.Song
	if err := c.ShouldBindJSON(&req); err != nil {
		receiver.logger.Set(c, logging.Error, http.StatusBadRequest, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}

	song := models.SongDto{
		Id: id,
		Song: models.Song{
			Title:    req.Title,
			Duration: req.Duration,
			Release:  req.Release,
			Author:   req.Author,
		},
	}

	err := receiver.db.Update(song)
	if err != nil {
		receiver.logger.Set(c, logging.Error, http.StatusBadRequest, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}

	receiver.logger.Set(c, logging.Info, http.StatusOK, "")
	c.JSON(http.StatusOK, song)
}

// Delete godoc
//
//	@Description	Delete song.
//	@Summary		Delete song
//	@Produce		json
//	@Tags			song
//	@Param			id	path	string	true	"Song ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/songs/{id} [DELETE]
func (receiver Song) Delete(c *gin.Context) {
	id := c.Param("id")

	err := receiver.db.Delete(id)
	if err != nil {
		receiver.logger.Set(c, logging.Error, http.StatusBadRequest, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}

	receiver.logger.Set(c, logging.Info, http.StatusNoContent, "")
	c.Status(http.StatusNoContent)
}
