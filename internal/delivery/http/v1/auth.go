package v1

import (
	"errors"
	"net/http"

	"github.com/Alexander272/my-portfolio/internal/domain"
	"github.com/Alexander272/my-portfolio/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-out", h.signOut)
		auth.POST("/refresh", h.refresh)
	}
}

type SignUpInput struct {
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

// @Summary SignUp
// @Tags auth
// @Description create user account
// @ModuleID userSignUp
// @Accept  json
// @Produce  json
// @Param input body SignUpInput true "sign up info"
// @Success 201 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var inp SignUpInput
	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	if err := h.services.User.SignUp(c.Request.Context(), service.SignUpInput{
		Name:     inp.Name,
		Email:    inp.Email,
		Password: inp.Password,
	}); err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, statusResponse{"Created"})
}

type SignInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}
type Token struct {
	AccessToken string `json:"accessToken"`
}

// @Summary SignIn
// @Tags auth
// @Description user sign in
// @ModuleID userSignIn
// @Accept  json
// @Produce  json
// @Param input body SignInInput true "sign in info"
// @Success 200 {object} Token
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var inp SignInInput
	if err := c.BindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}
	ua := c.GetHeader("sec-ch-ua") + " " + c.GetHeader("sec-ch-ua-platform") + " " + c.GetHeader("User-Agent")
	ip := c.ClientIP()

	cookie, token, err := h.services.Auth.SignIn(c.Request.Context(), service.SignInInput{
		Email:    inp.Email,
		Password: inp.Password,
	}, ua, ip)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	c.JSON(http.StatusOK, Token{
		AccessToken: token.AccessToken,
	})
}

// @Summary SignOut
// @Tags auth
// @Description выход из системы
// @ID logout
// @Accept json
// @Produce json
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-out [post]
func (h *Handler) signOut(c *gin.Context) {
	token, err := c.Cookie("refreshToken")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	cookie, err := h.services.Auth.SingOut(token)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	c.JSON(http.StatusOK, statusResponse{"Sign out success"})
}

// @Summary Refresh
// @Tags auth
// @Description обновление токенов доступа
// @Id refresh
// @Accept json
// @Produce json
// @Success 200 {object} Token
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	ua := c.GetHeader("sec-ch-ua") + " " + c.GetHeader("sec-ch-ua-platform") + " " + c.GetHeader("User-Agent")
	ip := c.ClientIP()

	token, err := c.Cookie(service.CookieName)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}
	tokens, cookie, err := h.services.Auth.Refresh(token, ua, ip)

	if err != nil {
		newErrorResponse(c, http.StatusForbidden, "Invalid request")
		return
	}

	c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	c.JSON(http.StatusOK, Token{
		AccessToken: tokens.AccessToken,
	})
}
