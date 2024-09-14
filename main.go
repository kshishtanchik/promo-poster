package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
	"gopkg.in/yaml.v2"
)

/*
const (
	BASE_SERVER_URL = "localhost:8181" //"http://188.120.235.226:8181"
)
*/

// текущие редактируемые сообщения
// [chatId][msgId][field]
var dialogRequest = make(map[int64]map[int]string)

type EventMetadata struct {
	ChatId    string `json:"chatId"`
	MessageId string `json:"messageId"`
}

type Student struct {
	Name          string `json:name`
	Phone         string `json:phone`
	TelegramAlias string `json:telegramAlias`
	EventMetadata
}

type EventData struct {
	Title         string `json:"title" description:"Заголовок события" yaml:"Title"`
	Description   string `json:"description" description:"Краткое описание мероприятия" yaml:"Description,omitempty"`
	EventDate     string `json:"eventDate" description:"Дата мероприятия" yaml:"EventDate"`
	EventTypeName string `json:"eventTypeName" description:"Тип мероприятия" yaml:"EventTypeName"`
	Period        string `json:"period" description:"Период проведения мероприятия в виде 13.00-15.00" yaml:"Period"`
	Address       string `json:"address" description:"Адресс мероприятия" yaml:"Address"`
}

type EventViewData struct {
	EventMetadata
	EventData
}

type Event struct {
	NewEvent `yaml:",newEvent"`
}

type NewEvent struct {
	Title         string `yaml:"Title"`
	Description   string `yaml:"Description,omitempty"`
	EventDate     string `yaml:"EventDate"`
	Address       string `yaml:"Period"`
	Period        string `yaml:"Address"`
	EventTypeName string `yaml:"EventTypeName"`
	ChatId        string `yaml:"ChatId"`
}

const (
	FINISH_TAG           = "#Завершено"
	USER_ERROR           = "Что-то пошло не так.. Уже смотрим что."
	EVENT_FINISH_MESSAGE = "Регистрация на мероприятие закончена"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}
	os.Hostname()
	BOT_TOKEN := os.Getenv("BOT_TOKEN")
	b, err := bot.New(BOT_TOKEN, opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/create", bot.MatchTypeExact, createHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/template", bot.MatchTypeExact, getTemplate)

	re := regexp.MustCompile(`^newEvent:`)

	b.RegisterHandlerRegexp(bot.HandlerTypeMessageText, re, newEventHandler)

	go b.Start(ctx)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			//todo: логировать данные для статистики
			title := r.URL.Query().Get("title")
			msg := r.URL.Query().Get("msg")
			eventTypeName := r.URL.Query().Get("eventTypeName")
			address := r.URL.Query().Get("address")
			period := r.URL.Query().Get("period")
			date := r.URL.Query().Get("date")
			chatId := r.URL.Query().Get("chatId")
			messageId := r.URL.Query().Get("messageId")

			event := EventData{
				Title:         title,
				Description:   msg,
				EventTypeName: eventTypeName,
				Address:       address,
				Period:        period,
				EventDate:     date,
			}
			eventMetadata := EventMetadata{
				MessageId: messageId,
				ChatId:    chatId,
			}
			data := EventViewData{
				EventData:     event,
				EventMetadata: eventMetadata,
			}
			tmpl, _ := template.ParseFiles("static/templates/index.html")
			tmpl.Execute(w, data)
		}

		if r.Method == "POST" {
			//todo: логировать данные для статистики
			var student Student
			var unmarshalErr *json.UnmarshalTypeError

			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			err := decoder.Decode(&student)
			if err != nil {
				if errors.As(err, &unmarshalErr) {
					errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
				} else {
					errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
				}
				return
			}
			messageId, _ := strconv.Atoi(student.MessageId)

			clearPhone := strings.ReplaceAll(student.Phone, " ", "")
			var alias string
			if student.TelegramAlias != "" {
				alias = "@" + student.TelegramAlias + ", "
			}

			studentString := fmt.Sprintf(` %s, %s %s`, student.Name, alias, clearPhone)
			newRm := inline.New(b, inline.NoDeleteAfterClick()).
				Row().
				Button("Пришел", []byte(studentString), present).
				Button("Не пришел", []byte(studentString), absent)

			studentData := fmt.Sprintf(`*[%s](https://t.me/%s)*, %s`, student.Name, student.TelegramAlias, bot.EscapeMarkdown(clearPhone))

			msgText := fmt.Sprintf(`Зарегистрирован: %s`, studentData)
			var registerMessage *models.Message
			registerMessage, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      student.ChatId,
				Text:        msgText,
				ReplyMarkup: newRm,
				ParseMode:   models.ParseModeMarkdown,
				ReplyParameters: &models.ReplyParameters{
					MessageID: messageId,
					ChatID:    student.ChatId,
				},
			})

			if err != nil {
				fmt.Println(err.Error())
				errorResponse(w, USER_ERROR, http.StatusBadRequest)
			}

			if strings.Contains(registerMessage.ReplyToMessage.Text, FINISH_TAG) {
				b.DeleteMessage(ctx, &bot.DeleteMessageParams{
					ChatID:    student.ChatId,
					MessageID: registerMessage.ID,
				})
				errorResponse(w, EVENT_FINISH_MESSAGE, http.StatusNotFound)
			}
			response := fmt.Sprintf(`{"data":"ok","name": "%s" }`, student.Name)
			fmt.Fprint(w, response)
		}
	})

	fmt.Println("Server is listening...")
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8181"
	}
	http.ListenAndServe(":"+PORT, nil)

}

func createUserText(eventData any) (resultMsg string) {
	ref := reflect.ValueOf(eventData)

	// if its a pointer, resolve its value
	if ref.Kind() == reflect.Ptr {
		ref = reflect.Indirect(ref)
	}

	if ref.Kind() == reflect.Interface {
		ref = ref.Elem()
	}

	// should double check we now have a struct (could still be anything)
	if ref.Kind() != reflect.Struct {
		log.Fatal("unexpected type")
	}

	eventDataType := reflect.TypeOf(eventData)
	fields := reflect.VisibleFields(eventDataType)
	for _, field := range fields {
		fieldVal := ref.FieldByName(field.Name).String()
		resultMsg = resultMsg + "\n" + fieldVal
	}
	return resultMsg
}

func createEditButtons(b *bot.Bot, eventData any) (kb *inline.Keyboard, resultMsg string) {

	ref := reflect.ValueOf(eventData)

	// if its a pointer, resolve its value
	if ref.Kind() == reflect.Ptr {
		ref = reflect.Indirect(ref)
	}

	if ref.Kind() == reflect.Interface {
		ref = ref.Elem()
	}

	// should double check we now have a struct (could still be anything)
	if ref.Kind() != reflect.Struct {
		log.Fatal("unexpected type")
	}

	eventDataType := reflect.TypeOf(eventData)
	fields := reflect.VisibleFields(eventDataType)
	kb = inline.New(b, inline.NoDeleteAfterClick())
	for _, field := range fields {
		displayName := reflect.StructTag.Get(field.Tag, "description")
		buttonLabel := "Редактировать " + displayName
		kb.Row().Button(buttonLabel, []byte(field.Name), setFieldVal)
		fieldVal := ref.FieldByName(field.Name).String()
		resultMsg = resultMsg + "\n" + displayName + ": " + fieldVal
	}
	return kb, resultMsg
}

func createTmplate(eventData any) (resultMsg string) {

	ref := reflect.ValueOf(eventData)

	// if its a pointer, resolve its value
	if ref.Kind() == reflect.Ptr {
		ref = reflect.Indirect(ref)
	}

	if ref.Kind() == reflect.Interface {
		ref = ref.Elem()
	}

	// should double check we now have a struct (could still be anything)
	if ref.Kind() != reflect.Struct {
		log.Fatal("unexpected type")
	}

	eventDataType := reflect.TypeOf(eventData)
	fields := reflect.VisibleFields(eventDataType)
	for _, field := range fields {
		description := reflect.StructTag.Get(field.Tag, "description")
		resultMsg = resultMsg + "\n" + field.Name + ": \"" + description + "\""
	}
	return resultMsg
}

func getTemplate(ctx context.Context, b *bot.Bot, update *models.Update) {
	newEvent := EventData{}
	resultMsg := bot.EscapeMarkdown("Скопируйте выделенный фрагмент. Замените текст полей в кавыйчках и отправте ответным сообщением. Бот выдаст ссылку для регистрации. Разместите ее в соцсетях")
	template := createTmplate(newEvent)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      resultMsg + "\n```yaml\nnewEvent:" + template + "\n```",
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func createHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	newEvent := EventData{}
	kb, resultMsg := createEditButtons(b, newEvent)

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        resultMsg,
		ReplyMarkup: kb,
	})
}

// формирование ссылки на событие
func newEventHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	event := NewEvent{}

	err := yaml.Unmarshal([]byte(update.Message.Text), &event)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", event)

	q := fmt.Sprintf(
		"title=%s&msg=%s&eventTypeName=%s&address=%s&period=%s&date=%s&chatId=%s&messageId=%s",
		event.Title,
		event.Description,
		event.EventTypeName,
		event.Address,
		event.Period,
		string(event.EventDate),
		strconv.FormatInt(update.Message.Chat.ID, 10),
		strconv.FormatInt(int64(update.Message.ID)+1, 10),
	)
	//todo: брать из настроек
	BASE_SERVER_URL := os.Getenv("BASE_SERVER_URL")
	PORT := os.Getenv("PORT")
	if PORT == "" && BASE_SERVER_URL == "" {
		PORT = "8181"
		BASE_SERVER_URL = "localhost"
	}
	link := fmt.Sprintf(
		"%s/event?%s",
		BASE_SERVER_URL+":"+PORT,
		url.PathEscape(q),
	)
	userText := createUserText(event)
	msg := fmt.Sprintf(
		"%s\nlink: ```%s```\n*%s*",
		bot.EscapeMarkdown(userText),
		link,
		bot.EscapeMarkdown("[Участники]"),
	)
	f := inline.New(b, inline.NoDeleteAfterClick())
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        msg,
		ReplyMarkup: f,
		ParseMode:   models.ParseModeMarkdown,
	})
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	// todo: обработка ввода ответов на диалоги
	chatId := update.Message.Chat.ID
	//todo: в dialogRequest[chatId] хранить модель
	dialog := dialogRequest[chatId]
	if dialog != nil {

	}
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
		Text:        mes.Message.Text,
	})
}

func setFieldVal(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	chatId := mes.Message.Chat.ID
	//todo: panic: assignment to entry in nil map
	// инициировать map
	dialogRequest[chatId][mes.Message.ID] = string(data)
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   "Введите " + string(data),
	})
}

func present(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {

	sourceMessageText := mes.Message.ReplyToMessage.Text
	newStudentRecord := fmt.Sprintf("\n + | %s", string(data))

	f := inline.New(b, inline.NoDeleteAfterClick()).
		Row().
		Button("Завершить", []byte("#Завершено"), finishEvent)

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      mes.Message.Chat.ID,
		MessageID:   mes.Message.ReplyToMessage.ID,
		ReplyMarkup: f,
		Text:        sourceMessageText + newStudentRecord,
	})

	if err != nil {
		fmt.Printf(err.Error())
	}

	// удалить mes.Message.
	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    mes.Message.Chat.ID,
		MessageID: mes.Message.ID,
	})
}

func finishEvent(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {

	var re = regexp.MustCompile(`link\:.+\n`)
	fixEvent := re.ReplaceAllString(mes.Message.Text, "\n")

	f := inline.New(b, inline.NoDeleteAfterClick())
	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      mes.Message.Chat.ID,
		MessageID:   mes.Message.ID,
		ReplyMarkup: f,
		Text:        string(data) + "\n" + fixEvent,
	})
}

func absent(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, data []byte) {
	sourceMessageText := mes.Message.ReplyToMessage.Text
	newStudentRecord := fmt.Sprintf("\n- | %s", string(data))

	f := inline.New(b, inline.NoDeleteAfterClick()).
		Row().
		Button("Завершить", []byte("#Завершено"), finishEvent)

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      mes.Message.Chat.ID,
		MessageID:   mes.Message.ReplyToMessage.ID,
		ReplyMarkup: f,
		Text:        sourceMessageText + newStudentRecord,
	})

	if err != nil {
		fmt.Printf(err.Error())
	}

	// удалить mes.Message.
	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    mes.Message.Chat.ID,
		MessageID: mes.Message.ID,
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
