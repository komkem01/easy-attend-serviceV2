package services

import (
	"easy-attend-service/configs"
	"easy-attend-service/models"
	"easy-attend-service/requests"
	"easy-attend-service/utils/logger"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type GenderService struct{}

func NewGenderService() *GenderService {
	return &GenderService{}
}

func (s *GenderService) GetAllGenders() ([]models.Gender, error) {
	logger.LogInfo("Fetching all genders", nil)

	var genders []models.Gender
	if err := configs.DB.Where("deleted_at IS NULL").Find(&genders).Error; err != nil {
		logger.LogError(err, "Failed to fetch genders", nil)
		return nil, errors.New("failed to fetch genders")
	}

	logger.LogInfo("Genders fetched successfully", logrus.Fields{
		"count": len(genders),
	})

	return genders, nil
}

func (s *GenderService) GetGenderByID(id uint) (*models.Gender, error) {
	logger.LogInfo("Fetching gender by ID", logrus.Fields{
		"gender_id": fmt.Sprintf("%d", id),
	})

	var gender models.Gender
	if err := configs.DB.Where("id = ? AND deleted_at IS NULL", id).First(&gender).Error; err != nil {
		logger.LogWarning("Gender not found", logrus.Fields{
			"gender_id": id,
		})
		return nil, errors.New("gender not found")
	}

	return &gender, nil
}

func (s *GenderService) CreateGender(req *requests.GenderCreateRequest) (*models.Gender, error) {
	logger.LogInfo("Creating new gender", logrus.Fields{
		"name": req.Name,
	})

	// Check if gender already exists
	var existingGender models.Gender
	if err := configs.DB.Where("name = ? AND deleted_at IS NULL", req.Name).First(&existingGender).Error; err == nil {
		logger.LogWarning("Gender creation failed - name already exists", logrus.Fields{
			"name": req.Name,
		})
		return nil, errors.New("gender with this name already exists")
	}

	// Create new gender
	gender := models.Gender{
		Name:      req.Name,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := configs.DB.Create(&gender).Error; err != nil {
		logger.LogError(err, "Failed to create gender", logrus.Fields{
			"name": req.Name,
		})
		return nil, errors.New("failed to create gender")
	}

	logger.LogInfo("Gender created successfully", logrus.Fields{
		"gender_id": fmt.Sprintf("%d", gender.ID),
		"name":      gender.Name,
	})

	return &gender, nil
}

func (s *GenderService) UpdateGender(id uint, req *requests.GenderUpdateRequest) (*models.Gender, error) {
	logger.LogInfo("Updating gender", logrus.Fields{
		"gender_id": id,
		"name":      req.Name,
	})

	var gender models.Gender
	if err := configs.DB.Where("id = ? AND deleted_at IS NULL", id).First(&gender).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogWarning("Gender not found for update", logrus.Fields{
				"gender_id": id,
			})
			return nil, errors.New("gender not found")
		}
		logger.LogError(err, "Failed to find gender for update", logrus.Fields{
			"gender_id": id,
		})
		return nil, errors.New("failed to find gender")
	}

	// Check if name is being changed and if it already exists
	if req.Name != gender.Name {
		var existingGender models.Gender
		if err := configs.DB.Where("name = ? AND id != ? AND deleted_at IS NULL", req.Name, id).First(&existingGender).Error; err == nil {
			logger.LogWarning("Gender update failed - name already exists", logrus.Fields{
				"gender_id": id,
				"name":      req.Name,
			})
			return nil, errors.New("gender with this name already exists")
		}
	}

	// Update fields
	gender.Name = req.Name
	gender.UpdatedAt = time.Now().Unix()

	if err := configs.DB.Save(&gender).Error; err != nil {
		logger.LogError(err, "Failed to update gender", logrus.Fields{
			"gender_id": id,
		})
		return nil, errors.New("failed to update gender")
	}

	logger.LogInfo("Gender updated successfully", logrus.Fields{
		"gender_id": fmt.Sprintf("%d", gender.ID),
		"name":      gender.Name,
	})

	return &gender, nil
}

func (s *GenderService) DeleteGender(id uint) error {
	logger.LogInfo("Deleting gender", logrus.Fields{
		"gender_id": id,
	})

	var gender models.Gender
	if err := configs.DB.Where("id = ? AND deleted_at IS NULL", id).First(&gender).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogWarning("Gender not found for deletion", logrus.Fields{
				"gender_id": id,
			})
			return errors.New("gender not found")
		}
		logger.LogError(err, "Failed to find gender for deletion", logrus.Fields{
			"gender_id": id,
		})
		return errors.New("failed to find gender")
	}

	// Soft delete
	deleteTime := time.Now().Unix()
	gender.DeletedAt = &deleteTime

	if err := configs.DB.Save(&gender).Error; err != nil {
		logger.LogError(err, "Failed to delete gender", logrus.Fields{
			"gender_id": id,
		})
		return errors.New("failed to delete gender")
	}

	logger.LogInfo("Gender deleted successfully", logrus.Fields{
		"gender_id": id,
	})

	return nil
}
