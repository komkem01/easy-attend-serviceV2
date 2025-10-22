package models

// LogAction enum for log actions
type LogAction string

const (
	LogActionLogin           LogAction = "login"
	LogActionLogout          LogAction = "logout"
	LogActionAttendance      LogAction = "attendance"
	LogActionCreateClassroom LogAction = "create_classroom"
	LogActionUpdateClassroom LogAction = "update_classroom"
	LogActionDeleteClassroom LogAction = "delete_classroom"
	LogActionCreateStudent   LogAction = "create_student"
	LogActionUpdateStudent   LogAction = "update_student"
	LogActionDeleteStudent   LogAction = "delete_student"
	LogActionCreateTeacher   LogAction = "create_teacher"
	LogActionUpdateTeacher   LogAction = "update_teacher"
	LogActionDeleteTeacher   LogAction = "delete_teacher"
)

type Log struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TeacherID uint      `gorm:"not null" json:"teacher_id"`
	Action    LogAction `gorm:"type:varchar(100);not null" json:"action"`
	Detail    string    `gorm:"type:text" json:"detail"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at"`
	SchoolID  *uint     `gorm:"default:null" json:"school_id,omitempty"`
}

func (l *Log) TableName() string {
	return "logs"
}

// IsValidAction checks if the provided action is valid
func (l *Log) IsValidAction() bool {
	switch l.Action {
	case LogActionLogin, LogActionLogout, LogActionAttendance,
		LogActionCreateClassroom, LogActionUpdateClassroom, LogActionDeleteClassroom,
		LogActionCreateStudent, LogActionUpdateStudent, LogActionDeleteStudent,
		LogActionCreateTeacher, LogActionUpdateTeacher, LogActionDeleteTeacher:
		return true
	default:
		return false
	}
}
