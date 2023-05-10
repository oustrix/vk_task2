package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
	"vk_task2/internal/domain"
	"vk_task2/internal/keyboards"
	"vk_task2/internal/usecase"
)

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
			msg := h.handleUpdate(update)
			_, err := h.bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}(update)
	}
}

func (h *handler) handleUpdate(update tgbotapi.Update) tgbotapi.Chattable {
	switch {
	case update.Message.IsCommand():
		return h.handleCommand(update)
	case update.Message != nil && update.Message.Text != "":
		return h.handleMessage(update)
	default:
		return h.handleUnknown(update)
	}
}

func (h *handler) handleCommand(update tgbotapi.Update) tgbotapi.Chattable {
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

func (h *handler) handleMessage(update tgbotapi.Update) tgbotapi.Chattable {
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

func (h *handler) handleStart(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Приветствую! Я бот, который поможет Вам управлять Вашими паролями и логами. Для работы со мною используйте кнопки ниже, либо команды.")
	msg.ReplyMarkup = keyboards.Basic()
	return msg
}

func (h *handler) handleUnknown(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не понимаю, что Вы от меня хотите."+
		"Для работы со мной используйте команды, либо кнопки снизу.")
	msg.ReplyMarkup = keyboards.Basic()
	return msg
}

func (h *handler) handleCreateHelp(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Для добавления нового логина и пароля используйте команду `/set`.\n\n"+
			"Пример: `/set <имя сервиса> <логин> <пароль>`\n\n"+
			"ВНИМАНИЕ: Имя сервиса, логин и пароль не должны содержать пробелов. ")
	msg.ReplyMarkup = keyboards.Basic()
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

func (h *handler) handleGetHelp(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Для получения логина и пароля используйте команду `/get`.\n\n"+
			"Пример: `/get <имя сервиса>`\n\n"+
			"Список сохраненных сервисов можно посмотреть, нажав на кнопку `Список сервисов`. ")
	msg.ReplyMarkup = keyboards.Basic()
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

func (h *handler) handleDeleteHelp(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Для удаления логинов и паролей для конкретного сервиса используйте команду `/delete`.\n\n"+
			"Пример: `/del <имя сервиса>`\n\n"+
			"Список сохраненных сервисов можно посмотреть, нажав на кнопку `Список сервисов`. ")
	msg.ReplyMarkup = keyboards.Basic()
	msg.ParseMode = tgbotapi.ModeMarkdown
	return msg
}

func (h *handler) handleList(update tgbotapi.Update) tgbotapi.Chattable {
	list, err := h.uc.Service.GetUserServices(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Произошла ошибка. Попробуйте позже.")
	}

	if len(list) == 0 {
		return tgbotapi.NewMessage(update.Message.Chat.ID, "У Вас нет сохраненных сервисов.")
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список Ваших сохраненных сервисов:\n")
	for _, service := range list {
		msg.Text += " - `" + service.Name + "`\n"
	}
	msg.ParseMode = tgbotapi.ModeMarkdown

	return msg
}

func (h *handler) handleSet(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyMarkup = keyboards.Basic()
	msg.ParseMode = tgbotapi.ModeMarkdown

	args := strings.Split(update.Message.CommandArguments(), " ")
	if len(args) < 3 {
		msg.Text = "Вы указали недостаточно параметров. Попробуйте еще раз."
		return msg
	} else if len(args) > 3 {
		msg.Text = "Вы указали слишком много параметров. Попробуйте еще раз."
		return msg
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
		return msg
	}

	msg.Text = "Сервис успешно добавлен."
	return msg
}

func (h *handler) handleGet(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = keyboards.Basic()

	args := strings.Split(update.Message.CommandArguments(), " ")
	if len(args) < 1 {
		msg.Text = "Вы забыли указать имя сервиса. Попробуйте еще раз."
		return msg
	}

	services, err := h.uc.Service.GetUserServicesByName(update.Message.Chat.ID, args[0])
	if err != nil {
		log.Println(err)
		msg.Text = "Произошла ошибка. Попробуйте позже."
		return msg
	}

	if len(services) == 0 {
		msg.Text = "Сервис с таким названием не найден."
		return msg
	}

	emoji := rune(128346)
	for _, service := range services {
		msg.Text += "- Логин: `" + service.Login + "` | Пароль: `" + service.Password + "`\n"
	}

	msg.Text += "\n" + "Данное сообщение будет удалено через 5 минут" + string(emoji)

	return msg
}

func (h *handler) handleDelete(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ReplyMarkup = keyboards.Basic()
	msg.ParseMode = tgbotapi.ModeMarkdown

	args := strings.Split(update.Message.CommandArguments(), " ")
	if len(args) < 1 {
		msg.Text = "Вы забыли указать имя сервиса. Попробуйте еще раз."
		return msg
	}

	err := h.uc.Service.Delete(update.Message.Chat.ID, args[0])
	if err != nil {
		log.Println(err)
		msg.Text = "Произошла ошибка. Попробуйте позже."
		return msg
	}

	msg.Text = "Значения для указанного сервиса успешно удалены."
	return msg
}
