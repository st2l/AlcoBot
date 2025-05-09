package handlers

import (
	"context"
	"log"
	"strings"
	"time"

	telegrambot "github.com/go-telegram/bot"
	telegramodels "github.com/go-telegram/bot/models"
	"github.com/st2l/AlcoBot/internal/config"
	"github.com/st2l/AlcoBot/internal/storage/mongodb"
)

func HandleDelGroup(ctx context.Context, b *telegrambot.Bot, update *telegramodels.Update) {
	if !(IsPrivateChat(update) && config.AppConfig.IsAdmin(update.Message.From.ID)) {
		log.Println("Someone tried to access admin method /del_group")
		return
	}

	text := update.Message.Text
	params := strings.Split(text, " ")
	if len(params) != 2 {
		b.SendMessage(ctx, &telegrambot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Вы не указали id группы. Введите команду вида: /del_group <group_id>",
		})
		return
	}

	group_id := params[1]
	client, err := mongodb.NewMongoClient(config.AppConfig.MongoUri, config.AppConfig.DBName, 3*time.Second)
	if err != nil {
		b.SendMessage(ctx, &telegrambot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка при подключении к БД...",
		})
		return
	}
	defer client.Disconnect()

	err = client.DeleteCollection(group_id)
	if err != nil {
		b.SendMessage(ctx, &telegrambot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка при удалении коллекции - проверьте корректность введеных параметров",
		})
		return
	}

	b.SendMessage(ctx, &telegrambot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Чат успешно удален из БД!",
	})
}
