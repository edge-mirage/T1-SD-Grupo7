package router

import (
	"github.com/gin-gonic/gin"

	"github.com/T1-SD/internal/handler"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/api/clients", handler.GetClients)
	r.POST("/api/clients", handler.CreateClient)
	r.GET("/api/clients/:id", handler.GetClientById)
	r.PUT("/api/clients/:id", handler.UpdateClientById)
	r.DELETE("/api/clients/:id", handler.DeleteClientById)

	r.GET("/api/users", handler.GetUsers)
	r.POST("/register", handler.CreateUser)
	r.POST("/login", handler.Login)
	r.GET("/api/users/:id", handler.GetUserById)
	r.PUT("/api/users/:id", handler.UpdateUserById)
	r.DELETE("/api/users/:id", handler.DeleteUserById)

	return r
}
