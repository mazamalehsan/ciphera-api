package routes

import (
	"ciphera-api/constants"
	"ciphera-api/db/models/loginChallengeModel"
	"ciphera-api/db/models/userModel"
	"ciphera-api/middlewares"
	"ciphera-api/response"
	"ciphera-api/routes/authRoutes"
	"ciphera-api/routes/userRoutes"
	"ciphera-api/services/authService"
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

	//use timeout middleware

	r.Use(middlewares.TimeoutMiddleware(30 * time.Second))

	//generic server status check route//

	r.GET("/status", func(c *gin.Context) {
		response.HandleSuccessResponse(c, "ok")
		return
	})

	//user routes handling//

	userRepo := userModel.New()
	userSvc := userService.New(userRepo)

	loginChallengeRepo := loginChallengeModel.New()

	authSvc := authService.New(userRepo, loginChallengeRepo)

	userRoutes.HandleUserRoutes(r, userSvc)
	authRoutes.HandleAuthRoutes(r, authSvc)

}
