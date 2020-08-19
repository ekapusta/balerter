package api

import (
	"fmt"
	"github.com/balerter/balerter/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/proxy"
	"net/http"
	"time"
)

const (
	defaultHTTPClientTimeout = time.Second * 5
)

type API struct {
	api        *tgbotapi.BotAPI
}

func New(cfg *config.ChannelTelegram) (*API, error) {
	botAPI, err := tgbotapi.NewBotAPI(cfg.Token)

	if err != nil {
		return nil, fmt.Errorf("error create bot api, %w", err)
	}

	a := &API{
		api: botAPI,
	}

	var tr *http.Transport

	if cfg.Proxy != nil {
		var proxyAuth *proxy.Auth

		if cfg.Proxy.Auth != nil {
			proxyAuth = &proxy.Auth{
				User:     cfg.Proxy.Auth.Username,
				Password: cfg.Proxy.Auth.Password,
			}
		}

		d, err := proxy.SOCKS5("tcp4", cfg.Proxy.Address, proxyAuth, nil)
		if err != nil {
			return nil, fmt.Errorf("error create proxy, %w", err)
		}

		tr = &http.Transport{
			Proxy:       nil,
			DialContext: nil,
			Dial:        d.Dial,
		}
	}

	a.api.Client.Transport = tr

	if a.api.Client.Timeout == 0 {
		a.api.Client.Timeout = defaultHTTPClientTimeout
	}

	return a, nil
}
