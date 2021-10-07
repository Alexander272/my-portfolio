package v1

import "github.com/gin-gonic/gin"

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	user := api.Group("/user", h.userIdentity)
	{
		user.PUT("/:id", h.notImplemented)
		user.DELETE("/:id", h.notImplemented)
	}
}
