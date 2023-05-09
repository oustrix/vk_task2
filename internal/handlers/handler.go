package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"vk_task2/internal/usecase"
)

type Handler interface {
	HandleUpdates()
}

type handler struct {
	uc      *usecase.Usecases
	updates *tgbotapi.UpdatesChannel
}

func NewHandler(uc *usecase.Usecases, updates *tgbotapi.UpdatesChannel) Handler {
	return &handler{
		uc:      uc,
		updates: updates,
	}
}

func (h *handler) HandleUpdates() {
	for update := range *h.updates {
		go h.handleUpdate(update)
	}
}

func (h *handler) handleUpdate(update tgbotapi.Update) {

}
