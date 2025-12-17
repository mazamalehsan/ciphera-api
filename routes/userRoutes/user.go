package userRoutes

import (
	"ciphera-api/response"
	"ciphera-api/services/userService"
	"ciphera-api/types"

	"github.com/gin-gonic/gin"
)

func HandleUserRoutes(r *gin.Engine, userSvc *userService.Service) {
	users := r.Group("/v1/users")
	users.POST("/fetch-by-username", func(c *gin.Context) {
		var body types.User
		if err := c.ShouldBindJSON(&body); err != nil {
			response.HandleErrorResponse(c, err)
		}
		user, err := userSvc.GetByUsername(
			c.Request.Context(),
			body.Username,
		)

		if err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		response.HandleSuccessResponse(c, user)
		return
	})

	users.POST("/check-username-duplication", func(c *gin.Context) {
		var body types.User
		if err := c.ShouldBindJSON(&body); err != nil {
			response.HandleErrorResponse(c, err)
		}
		exists, err := userSvc.UsernameExists(
			c.Request.Context(),
			body.Username,
		)

		if err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		response.HandleSuccessResponse(c, gin.H{"exists": exists})
		return
	})
}
