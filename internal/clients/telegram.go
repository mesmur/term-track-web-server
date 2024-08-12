package clients

import (
	"github.com/MESMUR/fixed-term-track-web-server/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type telegramSdk struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

type TelegramSdk interface {
	SendMessage(message string) error
}

func CreateTelegramSdk(botToken string, chatID int64) TelegramSdk {
	bot, err := tgbotapi.NewBotAPI(botToken)

	if err != nil {
		panic(err)
	}

	return &telegramSdk{
		bot:    bot,
		chatID: chatID,
	}
}

func (t *telegramSdk) SendMessage(message string) error {
	logger.Sugar.Info("Creating message")
	msg := tgbotapi.NewMessage(t.chatID, message)

	logger.Sugar.Info("Sending  message")
	if _, err := t.bot.Send(msg); err != nil {
		log.Panic(err)
	}

	logger.Log.Info("Message sent successfully!")

	return nil
}
