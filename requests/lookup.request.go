package requests

type GenderCreateRequest struct {
	Name string `json:"name" binding:"required,max=10"`
}

type GenderUpdateRequest struct {
	Name string `json:"name" binding:"required,max=10"`
}

type PrefixCreateRequest struct {
	Name string `json:"name" binding:"required,max=20"`
}

type PrefixUpdateRequest struct {
	Name string `json:"name" binding:"required,max=20"`
}
