package requests

import (
	"easy-attend-service/models"
)

// LogCreateRequest represents the request payload for creating a log
// Note: Logs are immutable - only creation is allowed, no updates or deletes
type LogCreateRequest struct {
	TeacherID uint             `json:"teacher_id" binding:"required" example:"1"`
	Action    models.LogAction `json:"action" binding:"required" example:"CREATE_STUDENT"`
	Detail    string           `json:"detail" binding:"max=500" example:"สร้างข้อมูลนักเรียนใหม่"`
}
