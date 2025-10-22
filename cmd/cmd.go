package cmd

import (
	"easy-attend-service/configs"
	"easy-attend-service/controller"
	"easy-attend-service/middlewares"
	"easy-attend-service/utils/logger"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "easy-attend-service",
	Short: "Easy Attend Service API",
	Long:  "A REST API service for Easy Attend System",
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  "Start the HTTP server to serve the API endpoints",
	Args:  NotReqArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Load environment variables
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, using system environment variables")
		}

		// Initialize logger
		logger.InitLogger()
		logger.LogInfo("Starting Easy Attend Service", nil)

		// Connect to database
		configs.ConnectDatabase()

		// Setup Gin mode
		ginMode := os.Getenv("GIN_MODE")
		if ginMode == "" {
			ginMode = "debug"
		}
		gin.SetMode(ginMode)

		// Create Gin router
		r := gin.Default()

		// Add logging middleware
		r.Use(middlewares.LoggingMiddleware())

		// Setup routes
		setupRoutes(r)

		// Get port from environment
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		log.Printf("Server starting on port %s", port)
		log.Fatal(r.Run(":" + port))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(Migrate())
}

func Execute() error {
	return rootCmd.Execute()
}

// NotReqArgs Not required arguments
func NotReqArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("not required arguments")
	}
	return nil
}

func setupRoutes(r *gin.Engine) {
	// Initialize controllers
	authController := controller.NewAuthController()
	teacherController := controller.NewTeacherController()
	studentController := controller.NewStudentController()
	schoolController := controller.NewSchoolController()
	genderController := controller.NewGenderController()
	prefixController := controller.NewPrefixController()
	classroomController := controller.NewClassroomController()
	classroomMemberController := controller.NewClassroomMemberController()
	attendanceController := controller.NewAttendanceController()
	logController := controller.NewLogController()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Easy Attend Service is running",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
		}

		// Test routes (public) - for testing only
		test := v1.Group("/test")
		{
			test.POST("/students", studentController.CreateStudent)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middlewares.AuthMiddleware())
		{
			// Auth profile and logout routes
			protected.GET("/auth/profile", authController.GetProfile)
			protected.POST("/auth/logout", authController.Logout)

			// Teacher routes
			teachers := protected.Group("/teachers")
			{
				teachers.GET("", teacherController.GetAllTeachers)
				teachers.POST("", teacherController.CreateTeacher)
				teachers.GET("/:id", teacherController.GetTeacherByID)
				teachers.PUT("/:id", teacherController.UpdateTeacher)
				teachers.DELETE("/:id", teacherController.DeleteTeacher)
			}

			// Student routes
			students := protected.Group("/students")
			{
				students.GET("", studentController.GetAllStudents)
				students.POST("", studentController.CreateStudent)
				students.GET("/:id", studentController.GetStudentByID)
				students.PUT("/:id", studentController.UpdateStudent)
				students.DELETE("/:id", studentController.DeleteStudent)
			}

			// School routes
			schools := protected.Group("/schools")
			{
				schools.GET("", schoolController.GetAllSchools)
				schools.POST("", schoolController.CreateSchool)
				schools.GET("/:id", schoolController.GetSchoolByID)
				schools.PUT("/:id", schoolController.UpdateSchool)
				schools.DELETE("/:id", schoolController.DeleteSchool)
			}

			// Gender routes
			genders := protected.Group("/genders")
			{
				genders.GET("", genderController.GetAllGenders)
				genders.POST("", genderController.CreateGender)
				genders.GET("/:id", genderController.GetGenderByID)
				genders.PUT("/:id", genderController.UpdateGender)
				genders.DELETE("/:id", genderController.DeleteGender)
			}

			// Prefix routes
			prefixes := protected.Group("/prefixes")
			{
				prefixes.GET("", prefixController.GetAllPrefixes)
				prefixes.POST("", prefixController.CreatePrefix)
				prefixes.GET("/:id", prefixController.GetPrefixByID)
				prefixes.PUT("/:id", prefixController.UpdatePrefix)
				prefixes.DELETE("/:id", prefixController.DeletePrefix)
			}

			// Classroom routes
			classrooms := protected.Group("/classrooms")
			{
				classrooms.GET("", classroomController.GetAllClassrooms)
				classrooms.POST("", classroomController.CreateClassroom)
				classrooms.GET("/:id", classroomController.GetClassroomByID)
				classrooms.PUT("/:id", classroomController.UpdateClassroom)
				classrooms.DELETE("/:id", classroomController.DeleteClassroom)
			}

			// Classroom Member routes
			classroomMembers := protected.Group("/classroom-members")
			{
				classroomMembers.GET("", classroomMemberController.GetAllClassroomMembers)
				classroomMembers.GET("/classroom/:classroom_id", classroomMemberController.GetClassroomMembersByClassroomID)
				classroomMembers.POST("", classroomMemberController.CreateClassroomMember)
				classroomMembers.PUT("/:classroom_id/:member_id", classroomMemberController.UpdateClassroomMember)
				classroomMembers.DELETE("/:classroom_id/:member_id", classroomMemberController.DeleteClassroomMember)
			}

			// Attendance routes
			attendances := protected.Group("/attendances")
			{
				attendances.GET("", attendanceController.GetAllAttendances)
				attendances.POST("", attendanceController.CreateAttendance)
				attendances.GET("/:id", attendanceController.GetAttendanceByID)
				attendances.PUT("/:id", attendanceController.UpdateAttendance)
				attendances.DELETE("/:id", attendanceController.DeleteAttendance)
				attendances.GET("/classroom/:classroom_id", attendanceController.GetAttendancesByClassroom)
				attendances.GET("/student/:student_id", attendanceController.GetAttendancesByStudent)
			}

			// Log routes (read-only + insert only - logs cannot be modified)
			logs := protected.Group("/logs")
			{
				logs.GET("", logController.GetAllLogs)                          // ดึงข้อมูล log ทั้งหมด
				logs.POST("", logController.CreateLog)                          // เพิ่ม log ใหม่
				logs.GET("/:id", logController.GetLogByID)                      // ดึงข้อมูล log ตาม ID
				logs.GET("/teacher/:teacherId", logController.GetLogsByTeacher) // ดึง logs ตาม teacher
				logs.GET("/action", logController.GetLogsByAction)              // ดึง logs ตาม action (query param)
			}
		}
	}
}
