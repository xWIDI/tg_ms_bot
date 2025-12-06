package cfg

import (
	"errors"
)

type Config struct {
	BotToken string
	ChatID   int64
}

func LoadConfig() (*Config, error) {
	botToken := "6522008688:AAGB5Y9Mre0tvijMnMn6Qg0IBcxqBsHsb2I"

	if botToken == "" {
		return nil, errors.New("необходимые переменные окружения отсутствуют")
	}

	//валидация токена бота + чата(если надо будет)

	return &Config{BotToken: botToken}, nil
}
