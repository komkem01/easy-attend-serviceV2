package requests

type StudentRequest struct {
	Page   int64  `json:"page" form:"page"`
	Size   int64  `json:"size" form:"size"`
	Search string `json:"search" form:"search"`
}

type StudentIdRequest struct {
	ID string `json:"id" uri:"id" binding:"required"`
}

type StudentCreateRequest struct {
	SchoolName string `json:"school_name" binding:"required"`
	StudentNo  string `json:"student_no"` // Made optional, will auto-generate if empty
	Firstname  string `json:"firstname" binding:"required"`
	Lastname   string `json:"lastname" binding:"required"`
	GenderID   *uint  `json:"gender_id"`
	PrefixID   *uint  `json:"prefix_id"`
}

type StudentUpdateRequest struct {
	SchoolName string `json:"school_name" binding:"required"`
	StudentNo  string `json:"student_no" binding:"required"`
	Firstname  string `json:"firstname" binding:"required"`
	Lastname   string `json:"lastname" binding:"required"`
}
