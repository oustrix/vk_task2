package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"time"
	"vk_task2/internal/domain"
	"vk_task2/internal/keyboards"
	"vk_task2/internal/usecase"
)

type config struct {
	isDelete      bool
	deleteTimeout time.Duration
}

type Handler interface {
	HandleUpdates()
}

type handler struct {
	bot     *tgbotapi.BotAPI
	uc      *usecase.Usecases
	updates *tgbotapi.UpdatesChannel
}

func NewHandler(bot *tgbotapi.BotAPI, uc *usecase.Usecases, updates *tgbotapi.UpdatesChannel) Handler {
	return &handler{
		bot:     bot,
		uc:      uc,
		updates: updates,
	}
}

func (h *handler) HandleUpdates() {
	for update := range *h.updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		go func(update tgbotapi.Update) {
			msg, cfg := h.handleUpdate(update)

			sent, err := h.bot.Send(msg)
			if err != nil {
				log.Println(err)
			}

			if cfg.isDelete {
				go func(chatID int64, messageID int, timeout time.Duration) {
					time.Sleep(timeout)
					_, err2 := h.bot.DeleteMessage(tgbotapi.NewDeleteMessage(chatID, messageID))
					if err2 != nil {
						log.Println(err2)
					}
				}(sent.Chat.ID, sent.MessageID, cfg.deleteTimeout)
			}
		}(update)
	}
}

func (h *handler) handleUpdate(update tgbotapi.Update) (tgbotapi.Chattable, *config) {
	switch {
	case update.Message.IsCommand():
		return h.handleCommand(update)
	case update.Message != nil && update.Message.Text != "":
		return h.handleMessage(update)
	default:
		return h.handleUnknown(update)
	}
}

func (h *handler) handleCommand(update tgbotapi.Update) (tgbotapi.Chattable, *config) {
	switch update.Message.Command() {
	case "start":
		return h.handleStart(update)
	case "set":
		return h.handleSet(update)
	case "get":
		return h.handleGet(update)
	case "del":
		return h.handleDelete(update)
	default:
		return h.handleUnknown(update)
	}
}

func (h *handler) handleMessage(update tgbotapi.Update) (tgbotapi.Chattable, *config) {
	switch update.Message.Text {
	case "Добавить":
		return h.handleCreateHelp(update)
	case "Получить":
		return h.handleGetHelp(update)
	case "Удалить":
		return h.handleDeleteHelp(update)
	case "Список сервисов":
		return h.handleList(update)
	default:
		return h.handleUnknown(update)
	}
}

func (h *handler) handleStart(update tgbotapi.Update) (tgbotapi.Chattable, *config) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую! Я бот, который поможет Вам управлять Вашими паролями и логами. Для работы со мною используйте кнопки ниже, либо команды.")
	msg.ReplyMarkup = keyboards.Basic()
	return msg, &config{}
}

func (h *handler) handleUnknown(update tgbotapi.Update) (tgbotapi.Chattable, *config) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не понимаю, что Вы от меня хотите."+
		"Для работы со мной используйте команды, либо кнопки снизу.")
	msg.ReplyMarkup = keyboards.Basic()
	return msg, &config{}
}

func (h *handler) handleCreateHelp(update tgbotapi.Update) (tgbotapi.Chattable, *config) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Для добавления нового логина и пароля используйте команду `/set`.\n\n"+
			"Пример: `/set <имя сервиса> <логин> <пароль>`\n\n"+
			"ВНИМАНИЕ: Имя сервиса, логин и пароль не должны содержать пробелов. ")
	msg.ReplyMarkup = keyboards.Basic()
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg, &config{}
}

func (h *handler) handleGetHelp(update tgbotapi.Update) (tgbotapi.Chattable, *config) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Для получения логина и пароля используйте команду `/get`.\n\n"+
			"Пример: `/get <имя сервиса>`\n\n"+
			"Список сохраненных сервисов можно посмотреть, нажав на кнопку `Список сервисов`. ")
	msg.ReplyMarkup = keyboards.Basic()
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg, &config{}
}

func (h *handler) handleDeleteHelp(update tgbotapi.Update) (tgbotapi.Chattable, *config) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Для удаления логинов и паролей для конкретного сервиса используйте команду `/delete`.\n\n"+
			"Пример: `/del <имя сервиса>`\n\n"+
			"Список сохраненных сервисов можно посмотреть, нажав на кнопку `Список сервисов`. ")
	msg.ReplyMarkup = keyboards.Basic()
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg, &config{}
}

func (h *handler) handleList(update tgbotapi.Update) (tgbotapi.Chattable, *config) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyMarkup = keyboards.Basic()
	msg.ParseMode = tgbotapi.ModeMarkdown

	list, err := h.uc.Service.GetUserServices(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		msg.Text = "Произошла ошибка при получении списка сервисов."
		return msg, &config{}
	}

	if len(list) == 0 {
		msg.Text = "У Вас нет сохраненных сервисов."
		return msg, &config{}
	}

	msg.Text += "Список Ваших сохраненных сервисов:\n"
	for _, service := range list {
		msg.Text += " - `" + service.Name + "`\n"
	}

	return msg, &config{}
}

func (h *handler) handleSet(update tgbotapi.Update) (tgbotapi.Chattable, *config) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyMarkup = keyboards.Basic()
	msg.ParseMode = tgbotapi.ModeMarkdown

	args := strings.Split(update.Message.CommandArguments(), " ")
	if len(args) < 3 {
		msg.Text = "Вы указали недостаточно параметров. Попробуйте еще раз."
		return msg, &config{}
	} else if len(args) > 3 {
		msg.Text = "Вы указали слишком много параметров. Попробуйте еще раз."
		return msg, &config{}
	}

	service := domain.Service{
		Name:     args[0],
		Login:    args[1],
		Password: args[2],
		UserID:   update.Message.Chat.ID,
	}

	err := h.uc.Service.Create(&service)
	if err != nil {
		log.Println(err)
		msg.Text = "Произошла ошибка. Попробуйте позже."
		return msg, &config{}
	}

	msg.Text = "Сервис успешно добавлен."
	return msg, &config{}
}

func (h *handler) handleGet(update tgbotapi.Update) (tgbotapi.Chattable, *config) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = keyboards.Basic()

	args := strings.Split(update.Message.CommandArguments(), " ")
	if len(args) < 1 {
		msg.Text = "Вы забыли указать имя сервиса. Попробуйте еще раз."
		return msg, &config{}
	}

	services, err := h.uc.Service.GetUserServicesByName(update.Message.Chat.ID, args[0])
	if err != nil {
		log.Println(err)
		msg.Text = "Произошла ошибка. Попробуйте позже."
		return msg, &config{}
	}

	if len(services) == 0 {
		msg.Text = "Сервис с таким названием не найден."
		return msg, &config{}
	}

	emoji := rune(128346)
	for _, service := range services {
		msg.Text += "- Логин: `" + service.Login + "` | Пароль: `" + service.Password + "`\n"
	}

	msg.Text += "\n" + "Данное сообщение будет удалено через 5 минут" + string(emoji)

	cfg := &config{
		isDelete:      true,
		deleteTimeout: 5 * time.Minute,
	}

	return msg, cfg
}

func (h *handler) handleDelete(update tgbotapi.Update) (tgbotapi.Chattable, *config) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyMarkup = keyboards.Basic()
	msg.ParseMode = tgbotapi.ModeMarkdown

	args := strings.Split(update.Message.CommandArguments(), " ")
	if len(args) < 1 {
		msg.Text = "Вы забыли указать имя сервиса. Попробуйте еще раз."
		return msg, &config{}
	}

	err := h.uc.Service.Delete(update.Message.Chat.ID, args[0])
	if err != nil {
		log.Println(err)
		msg.Text = "Произошла ошибка. Попробуйте позже."
		return msg, &config{}
	}

	msg.Text = "Значения для указанного сервиса успешно удалены."
	return msg, &config{}
}
