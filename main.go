package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Adrress struct {
	Street string
	House  int
	Name   string
}

type EventViewData struct {
	Title       string
	Description string
	Date        string
	Adrress
}

func main() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			title := r.URL.Query().Get("title")
			msg := r.URL.Query().Get("msg")
			data := EventViewData{
				Title:       title,
				Description: msg,
			}
			tmpl, _ := template.ParseFiles("static/templates/index.html")
			tmpl.Execute(w, data)
		}

		if r.Method == "POST" {
			answer := r.PostForm.Get("answer")
			// todo: обработать согласие
			fmt.Println(answer)
			//response := fmt.Sprintf("Product category=%s id=%s", cat, id)
			fmt.Fprint(w, `{"data":"responce"}`)
		}
	})

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}
