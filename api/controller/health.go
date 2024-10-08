package controller

import (
	"github.com/gin-gonic/gin"
	"main/database"
	"main/logging"
	"main/models"
	"net/http"
)

type Health struct {
	db *database.Song
}

func NewHealth(db *database.Song) *Health {
	return &Health{
		db: db,
	}
}

func set(c *gin.Context, level string, statusCode int, message string) {
	c.Set(logging.Level, level)
	c.Set(logging.StatusCode, statusCode)
	c.Set(logging.Message, message)
}

// Check godoc
//
//	@Description	Perform healthcheck.
//	@Summary		Perform healthcheck
//	@Produce		json
//	@Tags			health
//	@Success		200	{object}	models.Health
//	@Failure		500	{object}	models.Health
//	@Router			/health [GET]
func (receiver Health) Check(c *gin.Context) {
	if err := receiver.db.Ping(); err != nil {
		set(c, logging.Error, http.StatusInternalServerError, err.Error())

		c.JSON(http.StatusInternalServerError, models.Health{Healthy: false})
		return
	}

	set(c, logging.Info, http.StatusOK, "")
	c.JSON(http.StatusOK, models.Health{Healthy: true})
}
