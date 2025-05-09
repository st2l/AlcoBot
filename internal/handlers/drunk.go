package handlers

import (
	"context"
	"fmt"
	"time"

	telegrambot "github.com/go-telegram/bot"
	telegramodels "github.com/go-telegram/bot/models"
	"github.com/st2l/AlcoBot/internal/config"
	"github.com/st2l/AlcoBot/internal/storage/mongodb"
)

func HandleDrunk(ctx context.Context, b *telegrambot.Bot, update *telegramodels.Update) {
	client, err := mongodb.NewMongoClient(config.AppConfig.MongoUri, config.AppConfig.DBName, 3*time.Second)
	if err != nil {
		b.SendMessage(ctx, &telegrambot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Проблема в подключении к БД",
		})
		return
	}
	defer client.Disconnect()

	workingGroup, err := IsWorkingGroup(ctx, update, client)
	if err != nil {
		b.SendMessage(ctx, &telegrambot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Проблема при проверке группы...",
		})
		return
	}
	if !workingGroup {
		b.SendMessage(ctx, &telegrambot.SendMessageParams{
			ChatID: config.AppConfig.AdminIDs[0],
			Text:   fmt.Sprintf("Неавторизованная группа - %d", update.Message.Chat.ID),
		})
		return
	}

	drinkEvent, err := client.AddDrinkEvent(fmt.Sprintf("%d", update.Message.Chat.ID), update.Message.From.ID)
	if err != nil {
		b.SendMessage(ctx, &telegrambot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Проблема при создании пьянства TwT",
		})
		return
	}

	appTimeZone, _ := time.LoadLocation("Europe/Moscow")

	b.SendMessage(ctx, &telegrambot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text: fmt.Sprintf("🥃 *%s* не удержался и выпил!\n📅 *Время*: %s",
			update.Message.From.FirstName,
			drinkEvent.Timestamp.In(appTimeZone).Format("02.01.2006 15:04:05")),
		ParseMode: "Markdown",
	})

}
