package reservation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (this *implHealthAPI) GetHealth(ctx *gin.Context) {
  ctx.JSON(http.StatusOK, true)
}