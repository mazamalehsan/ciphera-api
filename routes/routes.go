package routes

import (
	"ciphera-api/constants"
	"ciphera-api/db/models/groupModel"
	"ciphera-api/db/models/loginChallengeModel"
	"ciphera-api/db/models/messageModel"
	"ciphera-api/db/models/userModel"
	"ciphera-api/middlewares"
	"ciphera-api/response"
	"ciphera-api/routes/authRoutes"
	"ciphera-api/routes/groupRoutes"
	"ciphera-api/routes/messageRoutes"
	"ciphera-api/routes/userRoutes"
	"ciphera-api/services/authService"
	"ciphera-api/services/groupService"
	"ciphera-api/services/messageService"
	"ciphera-api/services/userService"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func HandleRoutes(r *gin.Engine) {

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = constants.AllowedOrigins()
	corsConfig.AllowMethods = constants.AllowedMethods()
	corsConfig.AllowHeaders = constants.AllowedHeaders()
	r.Use(cors.New(corsConfig))

	r.Use(middlewares.TimeoutMiddleware(30 * time.Second))

	r.GET("/status", func(c *gin.Context) {
		response.HandleSuccessResponse(c, "ok")
		return
	})

	// Repos
	userRepo := userModel.New()
	loginChallengeRepo := loginChallengeModel.New()
	messageRepo := messageModel.New()
	groupRepo := groupModel.New()

	// Services
	authSvc := authService.New(userRepo, loginChallengeRepo)
	userSvc := userService.New(userRepo)
	msgSvc := messageService.New(messageRepo, userRepo)
	grpSvc := groupService.New(groupRepo, userRepo)

	// Public routes (no auth)
	authRoutes.HandleAuthRoutes(r, authSvc)

	// Protected routes (require JWT)
	protected := r.Group("/v1")
	protected.Use(middlewares.AuthMiddleware())

	userRoutes.HandleUserRoutes(r, userSvc)
	messageRoutes.HandleMessageRoutes(protected, msgSvc)
	groupRoutes.HandleGroupRoutes(protected, grpSvc)
}
