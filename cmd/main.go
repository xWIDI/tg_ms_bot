package main

import (
	"bot/internal/adapters"
	"bot/internal/cfg"
	"bot/internal/entities"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Алгоритм передачи данных карты в мс
func sendToMoySklad(card entities.Card, depot string) error {
	//Валидация карточки
	validating := entities.CardValidating(card)
	if validating != nil {
		return validating
	}
	//Создание описаний
	var mainDescription, subDescription string
	if len(depot) == 0 {
		return errors.New("ошибка заполнения поля склад, пустое поле")
	} else if strings.ToLower(depot) != "perm" && strings.ToLower(depot) != "moscow" {
		return errors.New("ошибка заполнения поля склад, такого склада не существует")
	} else {
		data, err := os.ReadFile(fmt.Sprintf("../internal/descriptions/main%s.txt", strings.ToUpper(depot)[0:1]+strings.ToLower(depot)[1:]))
		if err != nil {
			fmt.Println("Ошибка чтения файла")
		}
		mainDescription = string(data)
		data, err = os.ReadFile(fmt.Sprintf("../internal/descriptions/sub%s.txt", strings.ToUpper(depot)[0:1]+strings.ToLower(depot)[1:]))
		if err != nil {
			fmt.Println("Ошибка чтения файла")
		}
		subDescription = string(data)
	}
	fmt.Println(card, "\n", mainDescription, "\n\n", subDescription)
	//Логика передачи
	//bla bla bla
	//Возврат пустой ошибки в случае удачи
	return nil
}

func main() {
	cfg, err := cfg.LoadConfig()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации", err)
		os.Exit(1)
	}

	telegramAdapter, err := adapters.NewTelegramAdapter(cfg.BotToken)
	if err != nil {
		fmt.Println("Не удалось инициализировать TelegramAdapter", err)
		os.Exit(1)
	}

	telegramAdapter.GetMessage()
}
