package v1

import (
	"mime/multipart"
	"net/http"

	"github.com/Alexander272/my-portfolio/internal/domain"
	"github.com/Alexander272/my-portfolio/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	user := api.Group("/user", h.userIdentity)
	{
		user.GET("/:id", h.getUserById)
		user.PUT("/:id", h.updateUserById)
		user.DELETE("/:id", h.removeUserById)
	}
}

// @Summary Get User By Id
// @Security ApiKeyAuth
// @Tags user
// @Description получение данных пользователя
// @ModuleID getUserById
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} domain.User
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/{id} [get]
func (h *Handler) getUserById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, http.StatusBadRequest, "empty id param")
		return
	}
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	user, err := h.services.User.GetById(c, userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

type UserUpdateInput struct {
	Name     string         `form:"name" json:"name" binding:"min=2,max=64"`
	Email    string         `form:"email" json:"email" binding:"email,max=64"`
	Password string         `form:"password" json:"password" binding:"min=8,max=64"`
	Role     string         `form:"role" json:"role"`
	Avatar   multipart.File `json:"avatar"`
}

// @Summary Update User By Id
// @Security ApiKeyAuth
// @Tags user
// @Description обновление данных пользователя по его id
// @ModuleID updateUserById
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Param input body UserUpdateInput true "user info"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/{id} [put]
func (h *Handler) updateUserById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, http.StatusBadRequest, "empty id param")
		return
	}
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}
	var input UserUpdateInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	logger.Debug(input)

	// file, header, err := c.Request.FormFile("file")
	// if err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, "no file is received")
	// 	return
	// }

	err = h.services.User.UpdateById(c, userId, domain.UserUpdate{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     input.Role,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"Updated"})
}

// @Summary Remove User By Id
// @Security ApiKeyAuth
// @Tags user
// @Description удаление пользователя и всех его данных
// @ModuleID removeUserById
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /user/{id} [delete]
func (h *Handler) removeUserById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, http.StatusBadRequest, "empty id param")
		return
	}
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err = h.services.User.RemoveById(c, userId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"Removed"})
}
