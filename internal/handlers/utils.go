package handlers

import (
	"context"
	"fmt"

	telegramodels "github.com/go-telegram/bot/models"
	"github.com/st2l/AlcoBot/internal/storage/mongodb"
)

// IsGroupChat returns true if the message came from a group or supergroup
func IsGroupChat(update *telegramodels.Update) bool {
	if update.Message == nil {
		return false
	}

	// In Telegram API, chat.Type can be "private", "group", "supergroup", or "channel"
	return update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup"
}

// IsPrivateChat returns true if the message came from a private chat
func IsPrivateChat(update *telegramodels.Update) bool {
	if update.Message == nil {
		return false
	}

	return update.Message.Chat.Type == "private"
}

func IsWorkingGroup(ctx context.Context, update *telegramodels.Update, client *mongodb.MongoClient) (bool, error) {
	workingGroup, err := client.CheckWorkingGroup(fmt.Sprintf("%d", update.Message.Chat.ID))
	if err != nil {
		return false, err
	}

	if !(IsGroupChat(update) && workingGroup) {
		return false, nil
	}

	return true, nil
}
