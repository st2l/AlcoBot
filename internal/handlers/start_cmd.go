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

// HandleStart handles the /start command
func HandleStart(ctx context.Context, b *telegrambot.Bot, update *telegramodels.Update) {
	if IsPrivateChat(update) && config.AppConfig.IsAdmin(update.Message.From.ID) {
		txt := `Приветствую администратор!
Для работы с ботом тебе требуется:
1) Инициализировать нужную группу при помощи /add_group <group_id> ( в нужной группе можно прописать /getid )
2) Удостовериться что бот обладает правами администратора в группе
3) Попросить всех участников прописать команду /start (они таким образом отметятся в системе)
4) Радоваться!1!!`

		b.SendMessage(ctx, &telegrambot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   txt,
		})
	} else if IsGroupChat(update) {
		client, err := mongodb.NewMongoClient(config.AppConfig.MongoUri, config.AppConfig.DBName, 3*time.Second)
		if err != nil {
			b.SendMessage(ctx, &telegrambot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Проблема при подключении к БД...",
			})
			return
		}
		defer client.Disconnect()

		groupWorking, err := client.CheckWorkingGroup(fmt.Sprintf("%d", update.Message.Chat.ID))
		if err != nil {
			b.SendMessage(ctx, &telegrambot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Проблема при проверке группы на подключенную",
			})
			return
		}

		if !groupWorking {
			b.SendMessage(ctx, &telegrambot.SendMessageParams{
				ChatID: config.AppConfig.AdminIDs[0],
				Text:   fmt.Sprintf("Неавторизованная группа - %d", update.Message.Chat.ID),
			})
			return
		}

		_, err = client.GetOrCreateUser(fmt.Sprintf("%d", update.Message.Chat.ID), update.Message.From.ID, update.Message.From.FirstName)
		if err != nil {
			b.SendMessage(ctx, &telegrambot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Ошибка при создании или получении текущего пользователя",
			})
			return
		}

		txt := `Привет! Я бот который будет учитывать кол-во вашего трезвого времени!
Если вы ввели эту команду впервые - вы зарегистрированы в системе и учет начался
Для получения справки по всем доступным командам введите - /help`
		b.SendMessage(ctx, &telegrambot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   txt,
		})
	}

}
