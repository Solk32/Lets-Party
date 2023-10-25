package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Party struct {
	Name, Email, Phone string
	WillAttend         bool
}

var responses = make([]*Party, 0, 10)

var templates = make(map[string]*template.Template, 3)

func loadTemplates() {
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}
	for index, name := range templateNames {
		t, err := template.ParseFiles("layout.html", name+".html")
		if err == nil {
			templates[name] = t
			fmt.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}

type formData struct {
	*Party
	Errors []string
}

func main() {
	c := NewContainer()
	loadTemplates()

	http.HandleFunc("/", c.welcomeHandler)
	http.HandleFunc("/list", c.listHandler)
	http.HandleFunc("/form", c.formHandler)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
