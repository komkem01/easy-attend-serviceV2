package models

// AttendanceStatus enum for attendance status
type AttendanceStatus string

const (
	AttendanceStatusPresent AttendanceStatus = "present" // Present
	AttendanceStatusAbsent  AttendanceStatus = "absent"  // Absent
	AttendanceStatusLate    AttendanceStatus = "late"    // Late
	AttendanceStatusLeave   AttendanceStatus = "leave"   // Leave
)

type Attendance struct {
	ID          uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	ClassroomID *uint            `gorm:"not null" json:"classroom_id"`
	TeacherID   *uint            `gorm:"not null" json:"teacher_id"`
	StudentID   *uint            `gorm:"not null" json:"student_id"`
	SessionDate string           `gorm:"type:date;not null" json:"session_date"` // YYYY-MM-DD format
	Status      AttendanceStatus `gorm:"type:varchar(20);not null" json:"status"`
	CheckedAt   int64            `gorm:"not null" json:"checked_at"`
	Remark      string           `gorm:"type:text" json:"remark"`
	CreatedAt   int64            `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   int64            `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *int64           `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key Relationships
	Classroom *Classroom `gorm:"foreignKey:ClassroomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"classroom,omitempty"`
	Teacher   *Teacher   `gorm:"foreignKey:TeacherID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"teacher,omitempty"`
	Student   *Student   `gorm:"foreignKey:StudentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"student,omitempty"`
}

func (a *Attendance) TableName() string {
	return "attendances"
}

// IsValidStatus checks if the provided status is valid
func (a *Attendance) IsValidStatus() bool {
	switch a.Status {
	case AttendanceStatusPresent, AttendanceStatusAbsent, AttendanceStatusLate, AttendanceStatusLeave:
		return true
	default:
		return false
	}
}
