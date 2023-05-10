package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
	"vk_task2/config"
	"vk_task2/internal/handlers"
	serviceRepostiry "vk_task2/internal/repository/postgres/service"
	"vk_task2/internal/usecase"
	postgresConnect "vk_task2/pkg/postgres"
)

type App struct {
	cfg      *config.Config
	bot      *tgbotapi.BotAPI
	db       *gorm.DB
	usecases *usecase.Usecases
	updates  *tgbotapi.UpdatesChannel
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run() {
	// Bot connection.
	var err error
	a.bot, err = tgbotapi.NewBotAPI(a.cfg.Bot.Token)
	if err != nil {
		log.Fatalf("error while creating bot: %s", err)
	}
	log.Println("bot connected")

	// DB connection.
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		a.cfg.Postgres.Host, a.cfg.Postgres.Port, a.cfg.Postgres.User, a.cfg.Postgres.Password, a.cfg.Postgres.Database)
	a.db, err = postgresConnect.Connect(dsn, &gorm.Config{})
	if err != nil {
		log.Fatalf("error while connecting to db: %s", err)
	}
	log.Println("db connected")

	// Migrate.
	err = a.migrate()
	if err != nil {
		log.Fatalf("error while migrating: %s", err)
	}
	log.Println("db migrated")

	// Service.
	sRepository := serviceRepostiry.NewRepository(a.db)
	sUsecase := usecase.NewServiceUsecase(sRepository)

	// Usecases.
	a.usecases = &usecase.Usecases{
		Service: sUsecase,
	}

	// Updates.
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := a.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("error while getting updates: %s", err)
	}
	a.updates = &updates
	log.Println("got updates channel")

	// Handler.
	handler := handlers.NewHandler(a.bot, a.usecases, a.updates)
	go handler.HandleUpdates()

	// All below is for graceful shutdown.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		log.Printf("interrupt signal: %s\n", sig.String())
	case err = <-handler.Notify():
		log.Printf("handler error: %s\n", err)
	}

	log.Println("shutting down")
	handler.Shutdown(a.cfg.Handler.ShutdownTimeout)
	log.Println("shutdown complete")
}
