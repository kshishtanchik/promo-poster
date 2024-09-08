package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"gopkg.in/yaml.v2"
)

const (
	BOT_TOKEN       = "6459602535:AAGwuCpdrZMN130025aNXqRoJwOvXklqaLc"
	BASE_SERVER_URL = "http://188.120.235.226:8181"
)

type Student struct {
	Name          string `json:name`
	Phone         string `json:phone`
	TelegramAlias string `json:telegramAlias`
	ChatId        string `json:"chatId"`
}

type EventData struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	EventDate     string `json:"eventDate"`
	EventTypeName string `json:"eventTypeName"`
	Address       string `json:"address"`
	Period        string `json:"period"`
	ChatId        string `json:"chatId"`
}

type Event struct {
	NewEvent `yaml:",newEvent"`
}

type NewEvent struct {
	Title         string `yaml:"Title"`
	Description   string `yaml:"Description,omitempty"`
	EventDate     string `yaml:"EventDate"`
	EventTypeName string `yaml:"EventTypeName"`
	Address       string `yaml:"Address"`
	Period        string `yaml:"Period"`
	ChatId        string `yaml:"ChatId"`
}

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		//bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(BOT_TOKEN, opts...)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`^newEvent:`)

	b.RegisterHandlerRegexp(bot.HandlerTypeMessageText, re, newEventHandler)

	go b.Start(ctx)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			title := r.URL.Query().Get("title")
			msg := r.URL.Query().Get("msg")
			eventTypeName := r.URL.Query().Get("eventTypeName")
			address := r.URL.Query().Get("address")
			period := r.URL.Query().Get("period")
			date := r.URL.Query().Get("date")
			chatId := r.URL.Query().Get("chatId")

			data := EventData{
				Title:         title,
				Description:   msg,
				EventTypeName: eventTypeName,
				Address:       address,
				Period:        period,
				EventDate:     date,
				ChatId:        chatId,
			}
			tmpl, _ := template.ParseFiles("static/templates/index.html")
			tmpl.Execute(w, data)
		}

		if r.Method == "POST" {
			var e Student
			var unmarshalErr *json.UnmarshalTypeError

			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			err := decoder.Decode(&e)
			if err != nil {
				if errors.As(err, &unmarshalErr) {
					errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
				} else {
					errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
				}
				return
			}

			fmt.Println(e.ChatId)
			kb := inline.New(b, inline.NoDeleteAfterClick()).
				Row().
				Button("Написать в ЛС", []byte("0-1"), onConnect).
				Button("Отметить присутствие", []byte("0-1"), onMarkAttendance)

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      e.ChatId,
				Text:        fmt.Sprintf(`Новый ученик: %s, %s, %s`, e.Name, e.Phone, e.TelegramAlias),
				ReplyMarkup: kb,
			})
			response := fmt.Sprintf(`{"data":"ok","name": "%s" }`, e.Name)
			fmt.Fprint(w, response)
		}
	})

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)

}

func newEventHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	event := NewEvent{}

	err := yaml.Unmarshal([]byte(update.Message.Text), &event)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", event)

	f := inline.New(b, inline.NoDeleteAfterClick())

	link := fmt.Sprintf(
		`%s/event?title=%s&msg=%s&eventTypeName=%s&address=%s&period=%s&date=%s&chatId=%s`,
		BASE_SERVER_URL,
		event.Title,
		event.Description,
		event.EventTypeName,
		event.Address,
		event.Period,
		string(event.EventDate),
		strconv.FormatInt(update.Message.Chat.ID, 10),
	)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        link,
		ReplyMarkup: f,
	})
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := inline.New(b, inline.NoDeleteAfterClick()).
		Row().
		Button("Александр", []byte("0-0"), onInlineKeyboardSelect).
		Button("Отсутствует", []byte("0-1"), onInlineKeyboardSelect).
		Row().
		Button("Мстислав", []byte("1-0"), onInlineKeyboardSelect).
		Button("Отсутствует", []byte("1-1"), onInlineKeyboardSelect).
		Row().
		Button("Дмитрий", []byte("2-0"), onInlineKeyboardSelect).
		Button("Отсутствует", []byte("2-1"), onInlineKeyboardSelect)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        update.Message.Text,
		ReplyMarkup: kb,
	})
}

func onInlineKeyboardSelect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {

	kb := mes.Message.ReplyMarkup.InlineKeyboard

	f := inline.New(b, inline.NoDeleteAfterClick())
	changeID := strings.Split(string(data), "-")
	/*
		switch kb[changeID[0]][changeID[1]].Text {
		case "Отсутствует":
			kb[changeID[0]][changeID[1]].Text = "Присутствует"
		case "Присутствует":
			kb[changeID[0]][changeID[1]].Text = "Отсутствует"
		}
	*/
	for indexr := 0; indexr < len(kb); indexr++ {
		f.Row()
		for indexc := 0; indexc < len(kb[indexr]); indexc++ {
			var textResult string
			if changeID[0] == string(indexr) && changeID[1] == string(indexc) {
				switch kb[indexr][indexc].Text {
				case "Отсутствует":
					textResult = "Присутствует"
				case "Присутствует":
					textResult = "Отсутствует"
				default:
					textResult = kb[indexr][indexc].Text
				}
			} else {
				textResult = kb[indexr][indexc].Text
			}
			f.Button(textResult, []byte(string(indexr)+"-"+string(indexc)), onInlineKeyboardSelect)
		}
	}

	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      mes.Message.Chat.ID,
		MessageID:   mes.Message.ID,
		ReplyMarkup: f,
		Text:        "You selected: " + string(data),
	})
}

func onConnect(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	//todo: переход в личку пользователя
	f := inline.New(b, inline.NoDeleteAfterClick())
	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      mes.Message.Chat.ID,
		MessageID:   mes.Message.ID,
		ReplyMarkup: f,
		Text:        "You selected: " + string(data),
	})
}

func onMarkAttendance(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	//todo: отметить что пришел
	f := inline.New(b, inline.NoDeleteAfterClick())
	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      mes.Message.Chat.ID,
		MessageID:   mes.Message.ID,
		ReplyMarkup: f,
		Text:        "You selected: " + string(data),
	})
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
