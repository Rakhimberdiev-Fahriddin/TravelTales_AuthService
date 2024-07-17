package handler

import (
	"my_module/api/auth"
	pb "my_module/generated/auth_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Register user
// @Description create new users
// @Tags auth
// @Param info body auth_service.RegisterRequest true "User info"
// @Success 200 {object} auth_service.RegisterResponce
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "Server error"
// @Router /api/v1/auth/register [post]
func (h *Handler) RegisterUser(ctx *gin.Context) {
	h.Logger.Info("Register is starting")
	req := pb.RegisterRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		h.Logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	res, err := h.User.RegisterUser(&req)
	if err != nil {
		h.Logger.Error(err.Error())
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	h.Logger.Info("Register ended")
	ctx.JSON(http.StatusOK, res)
}

// @Summary login user
// @Description it generates new access and refresh tokens
// @Tags auth
// @Param LoginRequest body auth_service.LoginRequest true "username and password"
// @Success 200 {object} auth_service.Tokens
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/auth/login [post]
func (h *Handler) LoginUser(ctx *gin.Context) {
	h.Logger.Info("Login is working")
	req := pb.LoginRequest{}

	if err := ctx.BindJSON(&req); err != nil {
		h.Logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{"error1": err.Error()})
		return
	}

	res, err := h.User.GetUserByEmail(req.Email)
	if err != nil {
		h.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{"error2": err.Error()})
		return
	}

	var token pb.Tokens
	err = auth.GeneratedAccessJWTToken(res, &token)
	if err != nil {
		h.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{"error3": err.Error()})
		return
	}

	err = auth.GeneratedRefreshJWTToken(res, &token)
	if err != nil {
		h.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{"error4": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &token)
	h.Logger.Info("login is successfully ended")
}

// @Security ApiKeyAuth
// @Summary ResetPass user
// @Description it changes your password to new one
// @Tags auth
// @Param userinfo body auth_service.ResetPasswordRequest true "passwords"
// @Success 200 {object} string
// @Failure 400 {object} string "Invalid date"
// @Failure 401 {object} string "Invalid token"
// @Failure 500 {object} string "error while reading from server"
// @Router /api/v1/auth/reset-password [post]
func (h *Handler) ResetUserPassword(ctx *gin.Context) {
	h.Logger.Info("ResetPassword is working")

	accessToken := ctx.GetHeader("Authorization")
	id, err := auth.GetUserIdFromAccessToken(accessToken)
	req := pb.ResetPasswordRequest{UserId: id}
	if err != nil {
		h.Logger.Error(err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := ctx.BindJSON(&req); err != nil {
		h.Logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}
	_, err = h.User.ResetUserPassword(&req)
	if err != nil {
		h.Logger.Error(err.Error())
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "password successfully reset"})
	h.Logger.Info("ResetPassword ended")
}

// @Summary Refresh token
// @Description it changes your access token
// @Tags auth
// @Param userinfo body auth_service.RefreshTokenRequest true "token"
// @Success 200 {object} auth_service.Tokens
// @Failure 400 {object} string "Invalid date"
// @Failure 401 {object} string "Invalid token"
// @Failure 500 {object} string "error while reading from server"
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(ctx *gin.Context) {
	h.Logger.Info("Refresh is working")
	req := pb.RefreshTokenRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		h.Logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	_, err := auth.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		h.Logger.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	id, err := auth.GetUserIdFromRefreshToken(req.RefreshToken)
	if err != nil {
		h.Logger.Error(err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	res := pb.Tokens{Refreshtoken: req.RefreshToken}

	err = auth.GeneratedAccessJWTToken(&pb.LoginResponce{Id: id}, &res)
	if err != nil {
		h.Logger.Error(err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	ctx.JSON(http.StatusOK, &res)
}

// @Summary Logout user
// @Description you log out
// @Tags auth
// @Success 200 {object} string
// @Router /api/v1/auth/logout [post]
func (h *Handler) LogOutUser(ctx *gin.Context) {
	h.Logger.Info("Logout is working")
	ctx.JSON(http.StatusOK, gin.H{"message": "logged out"})
	h.Logger.Info("Logout ended")
}
