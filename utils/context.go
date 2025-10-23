package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetTeacherIDFromContext ดึง teacher ID จาก JWT context (รูปแบบ uint)
func GetTeacherIDFromContext(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return 0, errors.New("invalid user ID format in context")
	}

	// Parse as uint
	id, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return 0, errors.New("failed to parse teacher ID")
	}

	return uint(id), nil
}

// GetTeacherIDUUIDFromContext ดึง teacher ID จาก JWT context (รูปแบบ UUID - สำหรับระบบเก่า)
func GetTeacherIDUUIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return uuid.Nil, errors.New("invalid user ID format in context")
	}

	teacherID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("failed to parse teacher ID")
	}

	return teacherID, nil
}

// GetSchoolIDFromContext ดึง school ID จาก teacher context (ถ้ามี)
func GetSchoolIDFromContext(c *gin.Context) (*uuid.UUID, error) {
	schoolID, exists := c.Get("school_id")
	if !exists {
		return nil, nil // ไม่บังคับต้องมี school_id ใน context
	}

	schoolIDStr, ok := schoolID.(string)
	if !ok {
		return nil, errors.New("invalid school ID format in context")
	}

	parsedSchoolID, err := uuid.Parse(schoolIDStr)
	if err != nil {
		return nil, errors.New("failed to parse school ID")
	}

	return &parsedSchoolID, nil
}
