// internal/handlers/commands.go
package handlers

import (
	"context"

	telegrambot "github.com/go-telegram/bot"
	telegramodels "github.com/go-telegram/bot/models"
	"github.com/st2l/AlcoBot/internal/config"
)

// HandleHelp handles the /help command
func HandleHelp(ctx context.Context, b *telegrambot.Bot, update *telegramodels.Update) {

	help_text := ""
	if IsPrivateChat(update) && config.AppConfig.IsAdmin(update.Message.From.ID) {
		help_text = `АДМИНСКАЯ ПАНЕЛЬ!
Доступные команды:
/start - Начать работу с ботом (зачем)
/help - Получить эту сводку
/add_group - Добавить группу в которой будет работать бот (Убедитесь что бот является админом группы)
/del_group - Удалить группу в которой будет работать бот (Все данные будут утеряны)
`

	} else {
		help_text = `Доступные команды:
/start - Начать работу с ботом
/help - Показать эту справку
/stats - Показать статистику всех
/drunk - Обнулить свой счетчик
`
	}

	b.SendMessage(ctx, &telegrambot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   help_text,
	})
}
