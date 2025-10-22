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

type PrefixService struct{}

func NewPrefixService() *PrefixService {
	return &PrefixService{}
}

func (s *PrefixService) GetAllPrefixes() ([]models.Prefix, error) {
	logger.LogInfo("Fetching all prefixes", nil)

	var prefixes []models.Prefix
	if err := configs.DB.Where("deleted_at IS NULL").Find(&prefixes).Error; err != nil {
		logger.LogError(err, "Failed to fetch prefixes", nil)
		return nil, errors.New("failed to fetch prefixes")
	}

	logger.LogInfo("Prefixes fetched successfully", logrus.Fields{
		"count": len(prefixes),
	})

	return prefixes, nil
}

func (s *PrefixService) GetPrefixByID(id uint) (*models.Prefix, error) {
	logger.LogInfo("Fetching prefix by ID", logrus.Fields{
		"prefix_id": fmt.Sprintf("%d", id),
	})

	var prefix models.Prefix
	if err := configs.DB.Where("id = ? AND deleted_at IS NULL", id).First(&prefix).Error; err != nil {
		logger.LogWarning("Prefix not found", logrus.Fields{
			"prefix_id": id,
		})
		return nil, errors.New("prefix not found")
	}

	return &prefix, nil
}

func (s *PrefixService) CreatePrefix(req *requests.PrefixCreateRequest) (*models.Prefix, error) {
	logger.LogInfo("Creating new prefix", logrus.Fields{
		"name": req.Name,
	})

	// Check if prefix already exists
	var existingPrefix models.Prefix
	if err := configs.DB.Where("name = ? AND deleted_at IS NULL", req.Name).First(&existingPrefix).Error; err == nil {
		logger.LogWarning("Prefix creation failed - name already exists", logrus.Fields{
			"name": req.Name,
		})
		return nil, errors.New("prefix with this name already exists")
	}

	// Create new prefix
	prefix := models.Prefix{
		Name:      req.Name,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := configs.DB.Create(&prefix).Error; err != nil {
		logger.LogError(err, "Failed to create prefix", logrus.Fields{
			"name": req.Name,
		})
		return nil, errors.New("failed to create prefix")
	}

	logger.LogInfo("Prefix created successfully", logrus.Fields{
		"prefix_id": fmt.Sprintf("%d", prefix.ID),
		"name":      prefix.Name,
	})

	return &prefix, nil
}

func (s *PrefixService) UpdatePrefix(id uint, req *requests.PrefixUpdateRequest) (*models.Prefix, error) {
	logger.LogInfo("Updating prefix", logrus.Fields{
		"prefix_id": id,
		"name":      req.Name,
	})

	var prefix models.Prefix
	if err := configs.DB.Where("id = ? AND deleted_at IS NULL", id).First(&prefix).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogWarning("Prefix not found for update", logrus.Fields{
				"prefix_id": id,
			})
			return nil, errors.New("prefix not found")
		}
		logger.LogError(err, "Failed to find prefix for update", logrus.Fields{
			"prefix_id": id,
		})
		return nil, errors.New("failed to find prefix")
	}

	// Check if name is being changed and if it already exists
	if req.Name != prefix.Name {
		var existingPrefix models.Prefix
		if err := configs.DB.Where("name = ? AND id != ? AND deleted_at IS NULL", req.Name, id).First(&existingPrefix).Error; err == nil {
			logger.LogWarning("Prefix update failed - name already exists", logrus.Fields{
				"prefix_id": id,
				"name":      req.Name,
			})
			return nil, errors.New("prefix with this name already exists")
		}
	}

	// Update fields
	prefix.Name = req.Name
	prefix.UpdatedAt = time.Now().Unix()

	if err := configs.DB.Save(&prefix).Error; err != nil {
		logger.LogError(err, "Failed to update prefix", logrus.Fields{
			"prefix_id": id,
		})
		return nil, errors.New("failed to update prefix")
	}

	logger.LogInfo("Prefix updated successfully", logrus.Fields{
		"prefix_id": fmt.Sprintf("%d", prefix.ID),
		"name":      prefix.Name,
	})

	return &prefix, nil
}

func (s *PrefixService) DeletePrefix(id uint) error {
	logger.LogInfo("Deleting prefix", logrus.Fields{
		"prefix_id": id,
	})

	var prefix models.Prefix
	if err := configs.DB.Where("id = ? AND deleted_at IS NULL", id).First(&prefix).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.LogWarning("Prefix not found for deletion", logrus.Fields{
				"prefix_id": id,
			})
			return errors.New("prefix not found")
		}
		logger.LogError(err, "Failed to find prefix for deletion", logrus.Fields{
			"prefix_id": id,
		})
		return errors.New("failed to find prefix")
	}

	// Soft delete
	deleteTime := time.Now().Unix()
	prefix.DeletedAt = &deleteTime

	if err := configs.DB.Save(&prefix).Error; err != nil {
		logger.LogError(err, "Failed to delete prefix", logrus.Fields{
			"prefix_id": id,
		})
		return errors.New("failed to delete prefix")
	}

	logger.LogInfo("Prefix deleted successfully", logrus.Fields{
		"prefix_id": id,
	})

	return nil
}
