package requests

// ClassroomMemberCreateRequest represents the request payload for adding a member to classroom
type ClassroomMemberCreateRequest struct {
	ClassroomID uint  `json:"classroom_id" binding:"required"`
	TeacherID   *uint `json:"teacher_id,omitempty"`
	StudentID   *uint `json:"student_id,omitempty"`
}

// ClassroomMemberUpdateRequest represents the request payload for updating classroom member
type ClassroomMemberUpdateRequest struct {
	TeacherID *uint `json:"teacher_id,omitempty"`
	StudentID *uint `json:"student_id,omitempty"`
}
