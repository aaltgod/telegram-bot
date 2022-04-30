package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	logger *logrus.Logger
	botApi *tgbotapi.BotAPI
	api    *Api
}

func NewBot(logger *logrus.Logger, botApi *tgbotapi.BotAPI, api *Api) *Bot {
	return &Bot{
		logger: logger,
		botApi: botApi,
		api:    api,
	}
}

func (b *Bot) Start() error {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.botApi.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.logger.Warnln(err)
				b.handleError(update.Message)
			}
		} else {
			if err := b.handleText(update.Message); err != nil {
				b.logger.Warnln(err)
				b.handleError(update.Message)
			}
		}
	}

	return nil
}
