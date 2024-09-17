package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/database"
	"main/models"
	"net/http"
)

type Song struct {
	db *database.Song
}

func NewSong(db *database.Song) *Song {
	return &Song{
		db: db,
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
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}

	songDto, err := receiver.db.Create(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, songDto)
}

// GetById godoc
//
//	@Description	Get a song by id.
//	@Summary		Get a song by id
//	@Accept			json
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
		c.AbortWithStatusJSON(http.StatusBadRequest, models.Error{Message: fmt.Sprintf("song with id = '%s' was not found", id)})
		return
	}

	c.JSON(http.StatusOK, songDto)
}
