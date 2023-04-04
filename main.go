package main

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

type Meeting struct {
	Date     string
	Location string
	Comment  string
	User     string
}

func main() {
	bot, err := tgbotapi.NewBotAPI("API_KEY") // Замените на API-ключ от бота
	if err != nil {
		log.Panic(err)
	}

	channelUsername := "@channelName" // Замените на название канала, чата, пользователя
	chat, err := bot.GetChat(tgbotapi.ChatConfig{SuperGroupUsername: channelUsername})
	if err != nil {
		// обработка ошибки
	}
	channelID := chat.ID

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		meeting := Meeting{}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите дату встречи:") // Запрос даты
		bot.Send(msg)

		date := <-updates
		meeting.Date = strings.TrimSpace(date.Message.Text)
		meeting.User = strings.TrimSpace("@" + date.Message.From.UserName)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите место встречи:") // Запрос места
		bot.Send(msg)

		location := <-updates
		meeting.Location = strings.TrimSpace(location.Message.Text)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите комментарий к встрече:") // Запрос доп. информации, которая объяснит суть встречи
		bot.Send(msg)

		comment := <-updates
		meeting.Comment = strings.TrimSpace(comment.Message.Text)

		messageText := fmt.Sprintf("%s Объявляет сходон. \nДата: %s, \nМесто: %s, \nКомментарий гения: %s", meeting.User, meeting.Date, meeting.Location, meeting.Comment) // Меняйте текст итогового сообщения на свой вкус и цвет

		msg = tgbotapi.NewMessage(channelID, messageText)
		bot.Send(msg)
	}
}
