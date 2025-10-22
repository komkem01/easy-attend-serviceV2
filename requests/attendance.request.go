package requests

import (
	"easy-attend-service/models"
)

// AttendanceCreateRequest represents the request payload for creating attendance
type AttendanceCreateRequest struct {
	ClassroomID uint                    `json:"classroom_id" binding:"required"`
	TeacherID   uint                    `json:"teacher_id" binding:"required"`
	StudentID   uint                    `json:"student_id" binding:"required"`
	SessionDate string                  `json:"session_date" binding:"required"` // YYYY-MM-DD
	Status      models.AttendanceStatus `json:"status" binding:"required"`
	CheckedAt   int64                   `json:"checked_at" binding:"required"`
	Remark      string                  `json:"remark"`
}

// AttendanceUpdateRequest represents the request payload for updating attendance
type AttendanceUpdateRequest struct {
	Status    models.AttendanceStatus `json:"status" binding:"required,oneof=present absent late leave"`
	CheckedAt int64                   `json:"checked_at" binding:"required"`
	Remark    string                  `json:"remark"`
}
