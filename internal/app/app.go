package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"vk_task2/config"
)

type App struct {
	cfg *config.Config
	bot *tgbotapi.BotAPI
}

func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run() {
	var err error
	a.bot, err = tgbotapi.NewBotAPI(a.cfg.Bot.Token)
	if err != nil {
		log.Fatalf("error while creating bot: %s", err)
	}
}
