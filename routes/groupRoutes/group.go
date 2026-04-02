package groupRoutes

import (
	"ciphera-api/response"
	"ciphera-api/services/groupService"

	"github.com/gin-gonic/gin"
)

func HandleGroupRoutes(r *gin.RouterGroup, grpSvc *groupService.Service) {
	groups := r.Group("/groups")

	groups.POST("/create", func(c *gin.Context) {
		userId, _ := c.Get("userId")

		var body struct {
			Name    string
			Members []string
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		group, err := grpSvc.CreateGroup(c.Request.Context(), userId.(string), body.Name, body.Members)
		if err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		response.HandleSuccessResponse(c, group)
	})

	groups.GET("/my-groups", func(c *gin.Context) {
		userId, _ := c.Get("userId")

		grps, err := grpSvc.GetMyGroups(c.Request.Context(), userId.(string))
		if err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		response.HandleSuccessResponse(c, grps)
	})

	groups.POST("/members", func(c *gin.Context) {
		var body struct {
			GroupId string
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		members, err := grpSvc.GetGroupMembers(c.Request.Context(), body.GroupId)
		if err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		response.HandleSuccessResponse(c, members)
	})

	groups.POST("/add-member", func(c *gin.Context) {
		userId, _ := c.Get("userId")

		var body struct {
			GroupId  string
			Username string
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		err := grpSvc.AddMember(c.Request.Context(), userId.(string), body.GroupId, body.Username)
		if err != nil {
			response.HandleErrorResponse(c, err)
			return
		}

		response.HandleSuccessResponse(c, nil)
	})
}
