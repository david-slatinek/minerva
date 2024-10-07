package controller

import (
	"github.com/gin-gonic/gin"
	"main/database"
	"main/logging"
	"main/models"
	"net/http"
)

type Health struct {
	db     *database.Song
	logger *logging.Logging
}

func NewHealth(db *database.Song, logger *logging.Logging) *Health {
	return &Health{
		db:     db,
		logger: logger,
	}
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
		receiver.logger.Set(c, logging.Error, http.StatusInternalServerError, err.Error())

		c.JSON(http.StatusInternalServerError, models.Health{Healthy: false})
		return
	}
	c.JSON(http.StatusOK, models.Health{Healthy: true})
}
