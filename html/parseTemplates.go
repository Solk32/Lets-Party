package html

import (
	"fmt"
	"html/template"
	"log"
	"os"
)

var templates = make(map[string]*template.Template, 3)

func LoadTemplates() map[string]*template.Template {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}
	for index, name := range templateNames {
		t, err := template.ParseFiles(wd+"/html/layout.html", wd+"/html/"+name+".html")
		if err == nil {
			templates[name] = t
			fmt.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
	return templates
}
