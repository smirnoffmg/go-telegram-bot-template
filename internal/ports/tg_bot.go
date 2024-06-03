package ports

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	svc "github.com/smirnoffmg/partner/internal/services"
)

type Bot struct {
	bot         *tgbotapi.BotAPI
	chatService *svc.ChatGPTService
}

const infoMessage = "I'm a bot that can help you with your USML exam preparation.\nJust ask me anything.\nYou can contact my author @not_again_please for more info."

func NewTGBot(tgBotToken string, chatService *svc.ChatGPTService) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(tgBotToken)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	bot.Debug = false
	log.Info().Msgf("Authorized as %s", bot.Self.UserName)

	return &Bot{
		bot:         bot,
		chatService: chatService,
	}, nil
}

func (b *Bot) Start() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := b.bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Command() == "info" {
			infoMessageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, infoMessage)
			if _, err := b.bot.Send(infoMessageConfig); err != nil {
				log.Error().Err(err)
				return err
			}
			continue
		}

		answer, err := b.chatService.GetAnswer(update.Message.Chat.ID, update.Message.Text)

		if err != nil {
			log.Error().Err(err)
			return err
		}

		answerMessageConfig := tgbotapi.NewMessage(update.Message.Chat.ID, answer)

		if _, err := b.bot.Send(answerMessageConfig); err != nil {
			log.Error().Err(err)
			return err
		}
	}

	return nil
}