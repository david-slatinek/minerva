package controller

import (
	"github.com/gin-gonic/gin"
	"main/logging"
	"main/models"
	"net/http"
)

type Version struct {
	Version string
}

func NewVersion(version string) *Version {
	return &Version{
		Version: version,
	}
}

// GetVersion godoc
//
//	@Description	Get API version.
//	@Summary		Get API version
//	@Produce		json
//	@Tags			version
//	@Success		200	{object}	models.Version
//	@Failure		500	{object}	models.Error
//	@Router			/version [GET]
func (receiver Version) GetVersion(c *gin.Context) {
	if receiver.Version == "" {
		set(c, logging.Error, http.StatusInternalServerError, "version not found")
		c.JSON(http.StatusInternalServerError, models.Error{Message: "version not found"})
		return
	}

	set(c, logging.Info, http.StatusOK, "")
	c.JSON(http.StatusOK, models.Version{Version: receiver.Version})
}
