package telegram

import (
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"strings"
	"time"
	"yarl_intern_bot/internal/user"
)

type TgBot struct {
	bot      *bot.Bot
	ctx      context.Context
	chanData chan any
}

func New(ctx context.Context, token string, chanData chan any) TgBot {
	b, err := bot.New(token)
	if err != nil {
		panic(err)
	}

	bot := TgBot{
		bot:      b,
		chanData: chanData,
		ctx:      ctx,
	}

	bot.registerCommandHandlers()

	return bot
}

func (b *TgBot) registerCommandHandlers() {
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/addchannels", bot.MatchTypePrefix, b.addChannelsHandler())
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/adduser", bot.MatchTypePrefix, b.addUserHandler())
	b.bot.RegisterHandler(bot.HandlerTypeMessageText, "/settime", bot.MatchTypePrefix, b.setTimeHandler())
}

func (b *TgBot) sendMessage(id int64, message string) {
	b.bot.SendMessage(b.ctx,
		&bot.SendMessageParams{
			ChatID: id,
			Text:   message,
		})
}

func (b *TgBot) Run() {
	go b.bot.Start(b.ctx)

	for {
		select {
		case msg := <-b.chanData:
			switch msg.(type) {
			case []*user.User:
				log.Println("Distributing results")
				b.distributeResults(msg.([]*user.User))
			}
		}
	}
}

func (b *TgBot) addChannelsHandler() func(ctx context.Context, botInstance *bot.Bot, update *models.Update) {
	return func(ctx context.Context, botInstance *bot.Bot, update *models.Update) {
		channels := strings.Split(update.Message.Text, " ")[1:]
		b.chanData <- channels
		message := "Каналы успешно добавлены: " + strings.Join(channels, ", ")
		b.sendMessage(update.Message.Chat.ID, message)
	}
}
func (b *TgBot) addUserHandler() func(ctx context.Context, botInstance *bot.Bot, update *models.Update) {
	return func(ctx context.Context, botInstance *bot.Bot, update *models.Update) {
		// TODO: добавить функциональность добавления пользователя
		userName := strings.TrimSpace(strings.Split(update.Message.Text, " ")[1])

		message := "Пользователь " + userName + " успешно добавлен."
		b.sendMessage(update.Message.Chat.ID, message)
	}
}
func (b *TgBot) setTimeHandler() func(ctx context.Context, botInstance *bot.Bot, update *models.Update) {
	return func(ctx context.Context, botInstance *bot.Bot, update *models.Update) {
		newTime := strings.TrimSpace(strings.Split(update.Message.Text, " ")[1])

		parsedTime, err := time.Parse("15:04", newTime)
		if err != nil {
			message := "Неверный формат времени. Пожалуйста, введите время в формате HH:MM"
			b.sendMessage(update.Message.Chat.ID, message)
			return
		}
		b.chanData <- parsedTime
		message := "Время успешно установлено: " + newTime
		b.sendMessage(update.Message.Chat.ID, message)
	}
}

func (b *TgBot) distributeResults(users []*user.User) {
	for _, user := range users {
		resultsString := "Сегодня ничего нет :("
		if len(user.Results) > 0 {
			arr := make([]string, 0, len(user.Results))
			for key := range user.Results {
				arr = append(arr, key)
			}
			resultsString = strings.Join(arr, "\n")
		}

		b.sendMessage(int64(user.ID), resultsString)
	}
}
