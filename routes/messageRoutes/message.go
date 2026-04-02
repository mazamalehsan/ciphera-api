package messageRoutes

import (
	"ciphera-api/response"
	"ciphera-api/services/messageService"
	"ciphera-api/types"

	"github.com/gin-gonic/gin"
)

func HandleMessageRoutes(r *gin.RouterGroup, msgSvc *messageService.Service) {
	messages := r.Group("/messages")

	messages.POST("/send", func(c *gin.Context) {
		userId, _ := c.Get("userId")

		var body types.Message
		if err := c.ShouldBindJSON(&body); err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		msg, err := msgSvc.SendMessage(c.Request.Context(), userId.(string), body)
		if err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		response.HandleSuccessResponse(c, msg)
	})

	messages.POST("/conversation", func(c *gin.Context) {
		userId, _ := c.Get("userId")

		var body struct {
			OtherUserId string
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		msgs, err := msgSvc.GetConversation(c.Request.Context(), userId.(string), body.OtherUserId)
		if err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		response.HandleSuccessResponse(c, msgs)
	})

	messages.POST("/group", func(c *gin.Context) {
		var body struct {
			GroupId string
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		msgs, err := msgSvc.GetGroupMessages(c.Request.Context(), body.GroupId)
		if err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		response.HandleSuccessResponse(c, msgs)
	})
}
