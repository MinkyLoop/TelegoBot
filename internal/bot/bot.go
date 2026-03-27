package bot

import (
	"context"
	"fmt"
	"study/internal/models"

	tg "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Bot struct {
	bot     *tg.Bot
	handler *th.BotHandler
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewBot(parseFunc func() ([]models.Item, error)) (*Bot, error) {
	botToken := "8091607997:AAEQ7MbsUUi84Kmx62n71y2EBoObBAH5DiQ"
	bot, err := tg.NewBot(botToken, tg.WithDefaultDebugLogger())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())

	updates, err := bot.UpdatesViaLongPolling(ctx, nil)
	if err != nil {
		cancel()
		return nil, err
	}

	handler, err := th.NewBotHandler(bot, updates)
	if err != nil {
		cancel()
		return nil, err
	}
	defer func() { _ = handler.Stop() }()

	handler.Handle(func(ctx *th.Context, update tg.Update) error {
		chatID := update.Message.Chat.ID

		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(chatID),
			"Запускаю парсер",
		),
		)

		items, err := parseFunc()
		if err != nil {
			return err
		}
		if len(items) == 0 {
			return err
		}

		const maxItem = 15
		displayitems := items
		if len(items) > maxItem {
			displayitems = items[:maxItem]
		}

		msg := "Найденные товары\n\n"

		for i, item := range displayitems {
			msg += fmt.Sprintf("%d. %s\n", i+1, item.Title)
			msg += fmt.Sprintf("Цена %.2f рублей\n", float64(item.Prices.Price)/100)
			msg += fmt.Sprintf("Цена без скидки %.2f рублей\n", float64(item.Prices.PriceRegular)/100)
			msg += "\n"
		}

		if len(items) > maxItem {
			msg += fmt.Sprintf("А также %d товаров\n", maxItem-len(items))
		}

		bot.SendMessage(ctx, tu.Message(
			tu.ID(chatID),
			msg,
		))

		return nil
	}, th.CommandEqual("parse"))

	return &Bot{
		bot:     bot,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func (bot *Bot) Start() error {
	defer bot.cancel()
	defer bot.handler.Stop()
	return bot.handler.Start()
}
