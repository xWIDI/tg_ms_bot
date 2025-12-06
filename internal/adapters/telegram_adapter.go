package adapters

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramAdapter struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramAdapter(botToken string) (*TelegramAdapter, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = true

	fmt.Printf("Authorized on account %s", bot.Self.UserName)

	return &TelegramAdapter{bot: bot}, nil
}

// Тут же описываем все методы работы бота
// Для группировки изображений из одного сообщения
var mediaGroups = make(map[string][]string)
var mediaMutex sync.Mutex

func (t *TelegramAdapter) GetMessage() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			fmt.Println(update.Message.Text)

			if update.Message.Text == "stop" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Бот выключается")
				t.bot.Send(msg)
				cleanImagesFolder()
				os.Exit(0)
			}
		}

		// Обработка одиночного изображения
		if update.Message.Photo != nil && update.Message.MediaGroupID == "" {
			processSinglePhoto(t.bot, update.Message)
		}

		// Обработка группы изображений
		if update.Message.Photo != nil && update.Message.MediaGroupID != "" {
			processGroupedPhoto(t.bot, update.Message)
		}
	}
}

func processSinglePhoto(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	photo := (message.Photo)[len(message.Photo)-1]
	downloadPhoto(bot, photo.FileID, message.Chat.ID, "single_")
	bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Одиночное изображение сохранено!"))
}

func processGroupedPhoto(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	mediaGroupID := message.MediaGroupID
	photo := (message.Photo)[len(message.Photo)-1]

	mediaMutex.Lock()
	mediaGroups[mediaGroupID] = append(mediaGroups[mediaGroupID], photo.FileID)
	mediaMutex.Unlock()

	// Ждем немного, чтобы собрать все изображения группы
	go func(mediaGroupID string, chatID int64) {
		time.Sleep(2 * time.Second)

		mediaMutex.Lock()
		defer mediaMutex.Unlock()

		if files, exists := mediaGroups[mediaGroupID]; exists {
			for i, fileID := range files {
				downloadPhoto(bot, fileID, chatID, fmt.Sprintf("group_%s_%d_", mediaGroupID, i))
			}

			bot.Send(tgbotapi.NewMessage(chatID,
				fmt.Sprintf("Сохранено %d изображений из группы!", len(files))))

			delete(mediaGroups, mediaGroupID)
		}
	}(mediaGroupID, message.Chat.ID)
}

func downloadPhoto(bot *tgbotapi.BotAPI, fileID string, chatID int64, prefix string) {
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		fmt.Printf("Ошибка получения файла: %v\n", err)
		return
	}

	fileURL := "https://api.telegram.org/file/bot" + bot.Token + "/" + file.FilePath

	resp, err := http.Get(fileURL)
	if err != nil {
		fmt.Printf("Ошибка скачивания: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Создаем уникальное имя файла
	fileName := prefix + strconv.Itoa(int(time.Now().UnixNano())) + ".jpg"
	filePath := filepath.Join("../tmp/", fileName)
	os.MkdirAll("../tmp/", os.ModePerm)

	out, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v\n", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
		return
	}

	fmt.Printf("Изображение сохранено: %s\n", filePath)
}

func cleanImagesFolder() {
	dir := "../tmp"

	// Читаем содержимое папки
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Ошибка чтения папки: %v\n", err)
		return
	}

	// Удаляем все файлы в папке
	deletedCount := 0
	for _, file := range files {
		if !file.IsDir() {
			err := os.Remove(filepath.Join(dir, file.Name()))
			if err != nil {
				fmt.Printf("Ошибка удаления файла %s: %v\n", file.Name(), err)
			} else {
				deletedCount++
				fmt.Printf("Удален файл: %s\n", file.Name())
			}
		}
	}

	fmt.Printf("Очистка завершена. Удалено файлов: %d\n", deletedCount)
}
