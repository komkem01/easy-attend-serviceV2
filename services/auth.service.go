package services

import (
	"easy-attend-service/configs"
	"easy-attend-service/models"
	"easy-attend-service/requests"
	"easy-attend-service/utils"
	"easy-attend-service/utils/jwt"
	"easy-attend-service/utils/logger"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

type LoginResponse struct {
	Token     string         `json:"token"`
	Teacher   models.Teacher `json:"teacher"`
	ExpiresAt time.Time      `json:"expires_at"`
}

func (s *AuthService) Login(req *requests.LoginRequest) (*LoginResponse, error) {
	logger.LogInfo("User login attempt", logrus.Fields{
		"email": req.Email,
	})

	var teacher models.Teacher

	// Find teacher by email
	if err := configs.DB.Where("email = ?", req.Email).First(&teacher).Error; err != nil {
		logger.LogWarning("Login failed - user not found", logrus.Fields{
			"email": req.Email,
		})
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, teacher.Password) {
		logger.LogWarning("Login failed - invalid password", logrus.Fields{
			"email":   req.Email,
			"user_id": fmt.Sprintf("%d", teacher.ID),
		})
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	claims := jwt.CustomClaims{
		UserID:   fmt.Sprintf("%d", teacher.ID),
		Email:    teacher.Email,
		UserType: "teacher",
	}

	token, expiresAt, err := jwt.GenerateToken(claims)
	if err != nil {
		logger.LogError(err, "Failed to generate JWT token", logrus.Fields{
			"user_id": fmt.Sprintf("%d", teacher.ID),
			"email":   teacher.Email,
		})
		return nil, errors.New("failed to generate token")
	}

	logger.LogInfo("User login successful", logrus.Fields{
		"user_id": fmt.Sprintf("%d", teacher.ID),
		"email":   teacher.Email,
	})

	// Log activity automatically
	logger.LogActivity(teacher.ID, models.LogActionLogin, fmt.Sprintf("เข้าสู่ระบบด้วยอีเมล: %s", req.Email), teacher.SchoolID)

	return &LoginResponse{
		Token:     token,
		Teacher:   teacher,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) Register(req *requests.AuthRequest) (*models.Teacher, error) {
	// Check if teacher already exists
	var existingTeacher models.Teacher
	if err := configs.DB.Where("email = ?", req.Email).First(&existingTeacher).Error; err == nil {
		return nil, errors.New("teacher with this email already exists")
	}

	// Find or create school
	var school models.School
	err := configs.DB.Where("name = ?", req.SchoolName).First(&school).Error
	if err != nil {
		// School doesn't exist, create new one
		school = models.School{
			Name: req.SchoolName,
		}
		if err := configs.DB.Create(&school).Error; err != nil {
			return nil, errors.New("failed to create school")
		}
	}

	// Find or create gender if provided
	var genderID *uint
	if req.GenderName != "" {
		var gender models.Gender
		err := configs.DB.Where("name = ?", req.GenderName).First(&gender).Error
		if err != nil {
			// Gender doesn't exist, create new one
			gender = models.Gender{
				Name: req.GenderName,
			}
			if err := configs.DB.Create(&gender).Error; err != nil {
				return nil, errors.New("failed to create gender")
			}
		}
		genderID = &gender.ID
	}

	// Find or create prefix if provided
	var prefixID *uint
	if req.PrefixName != "" {
		var prefix models.Prefix
		err := configs.DB.Where("name = ?", req.PrefixName).First(&prefix).Error
		if err != nil {
			// Prefix doesn't exist, create new one
			prefix = models.Prefix{
				Name: req.PrefixName,
			}
			if err := configs.DB.Create(&prefix).Error; err != nil {
				return nil, errors.New("failed to create prefix")
			}
		}
		prefixID = &prefix.ID
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create new teacher
	teacher := models.Teacher{
		SchoolID:  &school.ID,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		GenderID:  genderID,
		PrefixID:  prefixID,
	}

	if err := configs.DB.Create(&teacher).Error; err != nil {
		return nil, errors.New("failed to create teacher")
	}

	// Log activity automatically
	logger.LogActivity(teacher.ID, models.LogActionCreateTeacher, fmt.Sprintf("ลงทะเบียนครูใหม่: %s %s (%s)", req.FirstName, req.LastName, req.Email), teacher.SchoolID)

	return &teacher, nil
}

func (s *AuthService) GetProfile(userID uint) (*models.Teacher, error) {
	var teacher models.Teacher
	if err := configs.DB.Preload("School").Preload("Gender").Preload("Prefix").Where("id = ?", userID).First(&teacher).Error; err != nil {
		return nil, errors.New("teacher not found")
	}
	return &teacher, nil
}
