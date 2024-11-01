package bot

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/kshishtanchik/promo-poster/cmd/bot/handlers"
)

var TgBot *bot.Bot
var Ctx *context.Context

func New(token string) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandler),
	}
	if token == "" {
		panic("No bot token!")
		os.Exit(2)
	}
	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
		os.Exit(2)
	}
	TgBot = b
	Ctx = &ctx
	go TgBot.Start(*Ctx)
}
