package requests

type AuthRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	FirstName  string `json:"firstname" binding:"required"`
	LastName   string `json:"lastname" binding:"required"`
	Phone      string `json:"phone"`
	SchoolName string `json:"school_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type TeacherRequest struct {
	Page   int64  `json:"page" form:"page"`
	Size   int64  `json:"size" form:"size"`
	Search string `json:"search" form:"search"`
}

type TeacherIdRequest struct {
	ID string `json:"id" uri:"id" binding:"required"`
}

type TeacherCreateRequest struct {
	SchoolID  uint   `json:"school_id" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
	Phone     string `json:"phone"`
}

type TeacherUpdateRequest struct {
	SchoolName string `json:"school_name"`
	Email      string `json:"email" binding:"omitempty,email"`
	Password   string `json:"password,omitempty"`
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	Phone      string `json:"phone"`
}
