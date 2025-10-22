package requests

type SchoolRequest struct {
	Page   int64  `json:"page" form:"page"`
	Size   int64  `json:"size" form:"size"`
	Search string `json:"search" form:"search"`
}

type SchoolIdRequest struct {
	ID string `json:"id" uri:"id" binding:"required"`
}

type SchoolCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type SchoolUpdateRequest struct {
	Name string `json:"name" binding:"required"`
}
