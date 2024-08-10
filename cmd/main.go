package main

import (
	"github.com/MESMUR/fixed-term-track-web-server/config"
	"github.com/MESMUR/fixed-term-track-web-server/controllers"
	"github.com/MESMUR/fixed-term-track-web-server/cron"
	"github.com/MESMUR/fixed-term-track-web-server/internal/clients"
	"github.com/MESMUR/fixed-term-track-web-server/internal/database"
	"github.com/MESMUR/fixed-term-track-web-server/pkg/logger"
	"github.com/MESMUR/fixed-term-track-web-server/repositories"
	routes "github.com/MESMUR/fixed-term-track-web-server/router"
	"github.com/MESMUR/fixed-term-track-web-server/services"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.Initialize()
	defer logger.Sync()

	config.LoadConfig()

	db := database.ConnectPostgres()

	investmentRepository := repositories.NewInvestmentRepository(db)
	investmentReturnRepository := repositories.NewInvestmentReturnRepository(db)
	eventRepository := repositories.NewEventRepository(db)
	investmentService := services.NewInvestmentService(investmentRepository, investmentReturnRepository, eventRepository)
	investmentController := controllers.NewUserController(investmentService)

	router := routes.SetupRouter(investmentController)

	telegramSdk := clients.CreateTelegramSdk(config.AppConfig.TelegramBotToken, config.AppConfig.TelegramChatID)
	eventReader := cron.NewEventReader(eventRepository, telegramSdk)
	go eventReader.CheckEvents()

	logger.Sugar.Infof("Starting server on port %s", config.AppConfig.Port)

	err := router.Run(config.AppConfig.Port)

	if err != nil {
		panic(err)
	}
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}
