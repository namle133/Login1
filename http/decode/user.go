package decode

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/namle133/Login1.git/LOGIN/domain"
)

func InputUser(c *gin.Context) *domain.UserInit {
	var creds *domain.UserInit
	err := c.BindJSON(&creds)
	if err != nil {
		c.String(http.StatusBadRequest, "400")
		return nil
	}
	if creds == nil {
		c.String(http.StatusBadRequest, "400")
		return nil
	}
	return creds
}
