package v1

import "github.com/gin-gonic/gin"

func (h *Handler) initProjectsRoutes(api *gin.RouterGroup) {
	project := api.Group("/projects")
	{
		project.GET("/", h.notImplemented)
		project.POST("/", h.notImplemented)
		project.GET("/:id", h.notImplemented)
		project.PUT("/:id", h.notImplemented)
		project.DELETE("/:id", h.notImplemented)
	}
}
