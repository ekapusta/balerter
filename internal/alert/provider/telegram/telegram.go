package telegram

import (
	"fmt"
	"github.com/balerter/balerter/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type API interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type Telegram struct {
	name   string
	chatID int64
	logger *zap.Logger
	api    API
}

func New(cfg *config.ChannelTelegram, logger *zap.Logger) (*Telegram, error) {
	tg := &Telegram{
		name:   cfg.Name,
		chatID: cfg.ChatID,
		logger: logger,
	}

	var err error
	tg.api, err = tgbotapi.NewBotAPI(cfg.Token)

	if err != nil {
		return nil, fmt.Errorf("error create bot api, %w", err)
	}

	return tg, nil
}

func (tg *Telegram) Name() string {
	return tg.name
}
