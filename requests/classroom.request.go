package requests

// ClassroomCreateRequest represents the request payload for creating a new classroom
type ClassroomCreateRequest struct {
	SchoolID  uint   `json:"school_id" binding:"required"`
	TeacherID uint   `json:"teacher_id" binding:"required"`
	Name      string `json:"name" binding:"required,min=1,max=255"`
	Grade     string `json:"grade" binding:"required,min=1,max=10"` // ชั้นเรียน
}

// ClassroomUpdateRequest represents the request payload for updating a classroom
type ClassroomUpdateRequest struct {
	SchoolID  uint   `json:"school_id" binding:"required"`
	TeacherID uint   `json:"teacher_id" binding:"required"`
	Name      string `json:"name" binding:"required,min=1,max=255"`
	Grade     string `json:"grade" binding:"required,min=1,max=10"` // ชั้นเรียน
}
