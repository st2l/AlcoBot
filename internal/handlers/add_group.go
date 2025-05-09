package handlers

import (
	"context"
	"log"
	"strings"

	telegrambot "github.com/go-telegram/bot"
	telegramodels "github.com/go-telegram/bot/models"
	"github.com/st2l/AlcoBot/internal/config"
	"github.com/st2l/AlcoBot/internal/storage/mongodb"
)

func HandleAddGroup(ctx context.Context, b *telegrambot.Bot, update *telegramodels.Update) {
	// check if admin and in personal
	if !(IsPrivateChat(update) && config.AppConfig.IsAdmin(update.Message.From.ID)) {
		log.Println("Someone tried to acces /add_group method while is not admin")
		return
	}

	text := update.Message.Text
	params := strings.Split(text, " ")
	if len(params) != 2 {
		b.SendMessage(ctx, &telegrambot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Вы не указали id группы. Введите команду вида: /add_group <group_id>",
		})
		return
	}

	group_id := params[1]
	client, err := mongodb.NewMongoClient(config.AppConfig.MongoUri, config.AppConfig.DBName, 3)
	if err != nil {
		b.SendMessage(ctx, &telegrambot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка при подключении к БД...",
		})
		return
	}
	defer client.Disconnect()

	err = client.InitializeCollection(group_id)
	if err != nil {
		b.SendMessage(ctx, &telegrambot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Ошибка при создании коллекции, проверьте корректность введенных вами параметров",
		})
		return
	}

	b.SendMessage(ctx, &telegrambot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Чат успешно добавлен в БД!",
	})
}
