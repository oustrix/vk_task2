package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
	"log"
	"vk_task2/config"
	serviceRepostiry "vk_task2/internal/repository/postgres/service"
	"vk_task2/internal/usecase"
	postgresConnect "vk_task2/pkg/postgres"
)

type App struct {
	cfg      *config.Config
	bot      *tgbotapi.BotAPI
	db       *gorm.DB
	usecases *usecase.Usecases
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
}
