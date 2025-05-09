package handlers

import (
	"context"
	"fmt"

	telegrambot "github.com/go-telegram/bot"
	telegramodels "github.com/go-telegram/bot/models"
)

func HandleGetID(ctx context.Context, b *telegrambot.Bot, update *telegramodels.Update) {
	b.SendMessage(ctx, &telegrambot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("ID чата в котором отправлена команда: %d", update.Message.Chat.ID),
	})
}
