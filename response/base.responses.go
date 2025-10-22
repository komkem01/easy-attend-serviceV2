package response

import (
	"easy-attend-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Status StatusResponse `json:"status"`
	Data   any            `json:"data,omitempty"`
}

type ResponsePaginate struct {
	Status     StatusResponse  `json:"status"`
	Data       any             `json:"data,omitempty"`
	Pagination models.Paginate `json:"paginate"`
}

type ResponsePaginate0 struct {
	Status     StatusResponse `json:"status"`
	Data       any            `json:"data,omitempty"`
	Pagination any            `json:"pagination"`
}

func SuccessWithPaginate(ctx *gin.Context, data any, pagination models.Paginate) {
	if pagination.Total == 0 {
		ctx.JSON(http.StatusOK, ResponsePaginate0{
			Status: StatusResponse{
				Code:    200,
				Message: "Success",
			},
			Data:       []any{},
			Pagination: gin.H{},
		})
		return
	} else {
		ctx.JSON(http.StatusOK, ResponsePaginate{
			Status: StatusResponse{
				Code:    200,
				Message: "Success",
			},
			Data:       data,
			Pagination: pagination,
		})
	}
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{StatusResponse{
		Code:    200,
		Message: "Success",
	}, data})
}

func SuccessResponse(message string, data any) Response {
	return Response{
		Status: StatusResponse{
			Code:    200,
			Message: message,
		},
		Data: data,
	}
}

func ErrorResponse(message string, detail string) StatusResponse {
	return StatusResponse{
		Code:    400,
		Message: message + ": " + detail,
	}
}

// InternalError ส่งผลลัพธ์เมื่อมีข้อผิดพลาดภายใน
func InternalError(ctx *gin.Context, message any, payloadCode ...string) {
	ctx.JSON(http.StatusInternalServerError, StatusResponse{
		Code:    500,
		Message: message.(string), // Set the message directly here
	})
}

func NotFound(ctx *gin.Context, message any, payloadCode ...string) {
	ctx.JSON(http.StatusNotFound, StatusResponse{
		Code:    404,
		Message: message.(string), // Set the message directly here
	})
}

// BadRequest ส่งผลลัพธ์เมื่อมีข้อผิดพลาดจากการขอข้อมูลที่ไม่ถูกต้อง
func BadRequest(ctx *gin.Context, message any, payloadCode ...string) {
	ctx.JSON(http.StatusBadRequest, StatusResponse{
		Code:    400,
		Message: message.(string), // Set the message directly here
	})
}

func Unauthorized(ctx *gin.Context, message any, payloadCode ...string) {
	ctx.JSON(http.StatusUnauthorized, StatusResponse{
		Code:    401,
		Message: message.(string), // Set the message directly here
	})
}

// Teacher response models
type TeacherResponses struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LoginResponse struct {
	Token   string           `json:"token"`
	Teacher TeacherResponses `json:"teacher"`
}
