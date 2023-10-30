package main

import (
	"fmt"
	"net/http"
	"partyinvites/html"
)

type Party struct {
	Name, Email, Phone string
	WillAttend         bool
}

var responses = make([]*Party, 0, 10)

type formData struct {
	*Party
	Errors []string
}

func main() {
	l := html.LoadTemplates()
	c := NewContainer(l)

	http.HandleFunc("/", c.welcomeHandler)
	http.HandleFunc("/list", c.listHandler)
	http.HandleFunc("/form", c.formHandler)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
}
