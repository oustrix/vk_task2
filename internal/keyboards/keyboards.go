package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func Basic() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить"),
			tgbotapi.NewKeyboardButton("Получить"),
			tgbotapi.NewKeyboardButton("Удалить"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Список сервисов"),
		),
	)
}
