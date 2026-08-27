package handlers

import (
	"backend/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

func (h *UserHandler) GetUsers(c *gin.Context) {

	users, err := h.Service.GetAllUsers()

	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{
			"error": "failed to fetch users",
		})
		return
	}

	c.JSON(200, gin.H{
		"users": users,
	})
}