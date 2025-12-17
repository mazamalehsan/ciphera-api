package authRoutes

import (
	"ciphera-api/response"
	"ciphera-api/services/authService"
	"ciphera-api/types"

	"github.com/gin-gonic/gin"
)

func HandleAuthRoutes(r *gin.Engine, authSvc *authService.Service) {
	auth := r.Group("/v1/auth")
	auth.POST("/register", func(c *gin.Context) {
		var body types.User
		if err := c.ShouldBindJSON(&body); err != nil {
			response.HandleErrorResponse(c, err)
		}
		err := authSvc.RegisterUser(
			c.Request.Context(),
			body,
		)

		if err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		response.HandleSuccessResponse(c, nil)
		return
	})

	auth.POST("/get-login-challange", func(c *gin.Context) {
		var body types.LoginBody
		if err := c.ShouldBindJSON(&body); err != nil {
			response.HandleErrorResponse(c, err)
		}
		challenge, loginId, err := authSvc.GenerateLoginChallenge(
			c.Request.Context(),
			body,
		)

		if err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		response.HandleSuccessResponse(c, gin.H{"challenge": challenge, "loginId": loginId})
		return
	})

	auth.POST("/login", func(c *gin.Context) {
		var body types.LoginBody
		if err := c.ShouldBindJSON(&body); err != nil {
			response.HandleErrorResponse(c, err)
		}
		token, err := authSvc.LoginUser(
			c.Request.Context(),
			body,
		)

		if err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		response.HandleSuccessResponse(c, gin.H{"token": token})
		return
	})
}
