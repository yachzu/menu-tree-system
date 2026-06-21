// @title STK Menu Tree System API
// @version 1.0
// @description RESTful API for hierarchical menu tree management
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@stk.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api
package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/xos/Projects/stk-project/backend/docs"
	"github.com/xos/Projects/stk-project/backend/internal/config"
	"github.com/xos/Projects/stk-project/backend/internal/handler"
	"github.com/xos/Projects/stk-project/backend/internal/model"
	"github.com/xos/Projects/stk-project/backend/internal/repository"
	"github.com/xos/Projects/stk-project/backend/internal/router"
	"github.com/xos/Projects/stk-project/backend/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&model.Menu{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	menuRepo := repository.NewMenuRepository(db)
	menuSvc := service.NewMenuService(menuRepo)
	menuHandler := handler.NewMenuHandler(menuSvc)

	r := router.Setup(menuHandler)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
