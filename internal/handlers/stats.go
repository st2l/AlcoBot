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

func HandleStats(ctx context.Context, b *telegrambot.Bot, update *telegramodels.Update) {
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

	listUsers, err := client.ListAllUsers(fmt.Sprintf("%d", update.Message.Chat.ID))
	if err != nil {
		b.SendMessage(ctx, &telegrambot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Проблема при получении списка пользователей",
		})
		return
	}

	var inlineKeyboard [][]telegramodels.InlineKeyboardButton

	for _, user := range listUsers {
		buttonText := fmt.Sprintf("User: %s", user.FirstName)
		callbackData := fmt.Sprintf("user_%d", user.TelegramID)
		row := []telegramodels.InlineKeyboardButton{
			{
				Text:         buttonText,
				CallbackData: callbackData,
			},
		}
		inlineKeyboard = append(inlineKeyboard, row)
	}

	b.SendMessage(ctx, &telegrambot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Выберите пользователя для просмотра статистики:",
		ReplyMarkup: &telegramodels.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard},
	})
}

func StatsCallbackHandler(ctx context.Context, b *telegrambot.Bot, update *telegramodels.Update) {
	if update.CallbackQuery == nil || len(update.CallbackQuery.Data) <= 5 || update.CallbackQuery.Data[:5] != "user_" {
		return
	}

	var userID int64
	_, err := fmt.Sscanf(update.CallbackQuery.Data, "user_%d", &userID)
	if err != nil {
		return
	}

	client, err := mongodb.NewMongoClient(config.AppConfig.MongoUri, config.AppConfig.DBName, 3*time.Second)
	if err != nil {
		b.AnswerCallbackQuery(ctx, &telegrambot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Проблема в подключении к БД",
			ShowAlert:       true,
		})
		return
	}
	defer client.Disconnect()

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	chatIDStr := fmt.Sprintf("%d", chatID)

	user, err := client.GetUserByTelegramID(chatIDStr, userID)
	if err != nil {
		b.AnswerCallbackQuery(ctx, &telegrambot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Пользователь не найден",
			ShowAlert:       true,
		})
		return
	}

	drinkCount := len(user.DrinkEvents)
	firstName := user.FirstName
	soberDuration := time.Since(user.LastSoberResetAt)
	soberDays := int(soberDuration.Hours() / 24)
	soberHours := int(soberDuration.Hours()) % 24

	statsText := fmt.Sprintf("Статистика пользователя: %s\n\n", firstName)
	statsText += fmt.Sprintf("Всего напитков записано: %d\n", drinkCount)
	statsText += fmt.Sprintf("Время трезвости: %d дней %d часов\n", soberDays, soberHours)

	if drinkCount > 0 {
		lastDrink := user.DrinkEvents[drinkCount-1].Timestamp
		statsText += fmt.Sprintf("Последний напиток: %s\n", lastDrink.Format("02.01.2006 15:04"))
	}

	b.EditMessageText(ctx, &telegrambot.EditMessageTextParams{
		ChatID:    chatID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      statsText,
		ReplyMarkup: &telegramodels.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegramodels.InlineKeyboardButton{
				{
					{Text: "← Назад к списку", CallbackData: "back_to_users"},
				},
			},
		},
	})

	b.AnswerCallbackQuery(ctx, &telegrambot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
}

func HandleBackToUsersList(ctx context.Context, b *telegrambot.Bot, update *telegramodels.Update) {
	if update.CallbackQuery == nil || update.CallbackQuery.Data != "back_to_users" {
		return
	}

	client, err := mongodb.NewMongoClient(config.AppConfig.MongoUri, config.AppConfig.DBName, 3*time.Second)
	if err != nil {
		b.AnswerCallbackQuery(ctx, &telegrambot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Проблема в подключении к БД",
			ShowAlert:       true,
		})
		return
	}
	defer client.Disconnect()

	chatID := update.CallbackQuery.Message.Message.Chat.ID
	listUsers, err := client.ListAllUsers(fmt.Sprintf("%d", chatID))
	if err != nil {
		b.AnswerCallbackQuery(ctx, &telegrambot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			Text:            "Проблема при получении списка пользователей",
			ShowAlert:       true,
		})
		return
	}

	var inlineKeyboard [][]telegramodels.InlineKeyboardButton
	for _, user := range listUsers {
		buttonText := fmt.Sprintf("User: %s", user.FirstName)
		callbackData := fmt.Sprintf("user_%d", user.TelegramID)
		row := []telegramodels.InlineKeyboardButton{
			{
				Text:         buttonText,
				CallbackData: callbackData,
			},
		}
		inlineKeyboard = append(inlineKeyboard, row)
	}

	b.EditMessageText(ctx, &telegrambot.EditMessageTextParams{
		ChatID:      chatID,
		MessageID:   update.CallbackQuery.Message.Message.ID,
		Text:        "Выберите пользователя для просмотра статистики:",
		ReplyMarkup: &telegramodels.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard},
	})

	b.AnswerCallbackQuery(ctx, &telegrambot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})
}
