package bot

import (
	"context"
	"log"

	telegrambot "github.com/go-telegram/bot"
	telegramodels "github.com/go-telegram/bot/models"
	"github.com/st2l/AlcoBot/internal/config"
	"github.com/st2l/AlcoBot/internal/handlers"
	"go.mongodb.org/mongo-driver/mongo"
)

// Bot represents the Telegram bot
type Bot struct {
	client      *telegrambot.Bot
	mongoClient *mongo.Client
}

// New creates a new Bot instance
func New() (*Bot, error) {
	// Create a new bot client
	opts := []telegrambot.Option{}

	client, err := telegrambot.New(config.AppConfig.TelegramToken, opts...)
	if err != nil {
		return nil, err
	}

	// Register command handlers
	client.RegisterHandler(telegrambot.HandlerTypeMessageText, "/start", telegrambot.MatchTypePrefix, handlers.HandleStart)
	client.RegisterHandler(telegrambot.HandlerTypeMessageText, "/help", telegrambot.MatchTypePrefix, handlers.HandleHelp)
	client.RegisterHandler(telegrambot.HandlerTypeMessageText, "/add_group", telegrambot.MatchTypePrefix, handlers.HandleAddGroup)
	client.RegisterHandler(telegrambot.HandlerTypeMessageText, "/getid", telegrambot.MatchTypeExact, handlers.HandleGetID)
	client.RegisterHandler(telegrambot.HandlerTypeMessageText, "/del_group", telegrambot.MatchTypePrefix, handlers.HandleDelGroup)
	client.RegisterHandler(telegrambot.HandlerTypeMessageText, "/drunk", telegrambot.MatchTypePrefix, handlers.HandleDrunk)
	client.RegisterHandler(telegrambot.HandlerTypeMessageText, "/stats", telegrambot.MatchTypePrefix, handlers.HandleStats)
	client.RegisterHandler(telegrambot.HandlerTypeCallbackQueryData, "user_", telegrambot.MatchTypePrefix, handlers.StatsCallbackHandler)
	client.RegisterHandler(telegrambot.HandlerTypeCallbackQueryData, "back_to_users", telegrambot.MatchTypeExact, handlers.HandleBackToUsersList)

	// Set bot commands
	commands := []telegramodels.BotCommand{
		{Command: "start", Description: "Начните работу с ботом / Зарегистрировать себя в группе"},
		{Command: "help", Description: "Показать все команды"},
		{Command: "stats", Description: "Показать статистику по пользователям этой группы"},
		{Command: "drunk", Description: "Обнулить свой счетчик"},
		// Add more commands as needed
	}
	_, err = client.SetMyCommands(context.Background(), &telegrambot.SetMyCommandsParams{
		Commands: commands,
	})
	if err != nil {
		log.Printf("Warning: Failed to set commands: %v", err)
	}

	return &Bot{
		client: client,
	}, nil
}

func (b *Bot) Start() error {
	// Start the bot in a separate goroutine
	go b.client.Start(context.Background())
	log.Println("Bot started")

	return nil
}
