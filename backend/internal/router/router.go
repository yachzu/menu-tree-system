package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/xos/Projects/stk-project/backend/internal/handler"
	"github.com/xos/Projects/stk-project/backend/internal/middleware"
)

func Setup(h *handler.MenuHandler) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandler())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.GET("/menus", h.GetTree)
		api.GET("/menus/:id", h.GetByID)
		api.POST("/menus", h.Create)
		api.PUT("/menus/:id", h.Update)
		api.DELETE("/menus/:id", h.Delete)
		api.PATCH("/menus/:id/move", h.Move)
		api.PATCH("/menus/:id/reorder", h.Reorder)
	}

	return r
}
