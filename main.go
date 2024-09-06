package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type EventViewData struct {
	Title         string
	Description   string
	EventDate     string
	EventTypeName string
	Address       string
	Period        string
	ChatId        string
}

func main() {
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

			data := EventViewData{
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
			name := r.PostForm.Get("name")
			// todo: обработать согласие
			fmt.Println(name)
			//response := fmt.Sprintf("Product category=%s id=%s", cat, id)
			response := fmt.Sprintf(`{"data":"ok","name": "%n" }`, name)
			fmt.Fprint(w, response)
		}
	})

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
