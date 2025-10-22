package controller

import (
	"easy-attend-service/models"
	"easy-attend-service/requests"
	"easy-attend-service/response"
	"easy-attend-service/services"
	"easy-attend-service/utils"
	"easy-attend-service/utils/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

func (ac *AuthController) Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	result, err := ac.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Login failed", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Login successful", result))
}

func (ac *AuthController) Register(c *gin.Context) {
	var req requests.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request data", err.Error()))
		return
	}

	teacher, err := ac.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Registration failed", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Teacher registered successfully", teacher))
}

func (ac *AuthController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized", "User ID not found"))
		return
	}

	teacher, err := ac.authService.GetProfile(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse("Profile not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Profile retrieved successfully", teacher))
}

// Logout godoc
// @Summary Logout user
// @Description Logout user and log the activity
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /auth/logout [post]
func (ac *AuthController) Logout(c *gin.Context) {
	// Get teacher ID from context
	teacherID, err := utils.GetTeacherIDFromContext(c)
	if err != nil {
		logger.LogError(err, "Failed to get teacher ID from context", logrus.Fields{})
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized", err.Error()))
		return
	}

	// Get teacher info for school ID
	teacher, err := ac.authService.GetProfile(teacherID.String())
	if err != nil {
		logger.LogError(err, "Failed to get teacher profile for logout", logrus.Fields{
			"teacher_id": teacherID.String(),
		})
		// Still allow logout even if we can't get profile
	}

	// Log logout activity
	var schoolID *uint
	if teacher != nil {
		schoolID = teacher.SchoolID
	}

	var teacherIDUint *uint
	if teacher != nil {
		teacherIDUint = &teacher.ID
	}

	var teacherIDVal uint
	if teacherIDUint != nil {
		teacherIDVal = *teacherIDUint
	}

	logger.LogActivity(teacherIDVal, models.LogActionLogout, fmt.Sprintf("ออกจากระบบ"), schoolID)

	logger.LogInfo("User logged out successfully", logrus.Fields{
		"teacher_id": teacherID.String(),
	})

	c.JSON(http.StatusOK, response.SuccessResponse("Logged out successfully", nil))
}
