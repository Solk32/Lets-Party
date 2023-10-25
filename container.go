package main

import (
	"database/sql"
	"net/http"
	"partyinvites/db"
)

type Container struct {
	DB *sql.DB
}

func NewContainer() *Container {
	c := new(Container)
	c.initDatabase()
	return c
}

func (c *Container) initDatabase() {
	c.DB = db.ConnectDB()
}

func (c *Container) welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	templates["welcome"].Execute(writer, nil)
}

func (c *Container) listHandler(writer http.ResponseWriter, request *http.Request) {
	templates["list"].Execute(writer, responses)
}

func (c *Container) formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		templates["form"].Execute(writer, formData{
			Party: &Party{}, Errors: []string{},
		})
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		c.CreateGuest(request.Form["name"][0], request.Form["email"][0], request.Form["phone"][0], request.Form["willattend"][0])
		responseData := Party{
			Name:       request.Form["name"][0],
			Email:      request.Form["email"][0],
			Phone:      request.Form["phone"][0],
			WillAttend: request.Form["willattend"][0] == "true",
		}

		errors := []string{}
		if responseData.Name == "" {
			errors = append(errors, "Пожалуйста введите своё имя")
		}
		if responseData.Email == "" {
			errors = append(errors, "Пожалуйста введите Ваш Email")
		}
		if responseData.Phone == "" {
			errors = append(errors, "Пожалуйста введите Ваш телефон")
		}
		if len(errors) > 0 {
			templates["form"].Execute(writer, formData{
				Party: &responseData, Errors: errors,
			})
		} else {
			responses = append(responses, &responseData)
			if responseData.WillAttend {
				templates["thanks"].Execute(writer, responseData.Name)
			} else {
				templates["sorry"].Execute(writer, responseData.Name)
			}
		}
	}
}

func (c *Container) CreateGuest(name string, email string, phone string, join string) {
	sqlStatement := `INSERT INTO guests (name, email, phone, userjoin) 
VALUES ($1, $2, $3, $4)`
	var err error
	_, err = c.DB.Exec(sqlStatement, name, email, phone, join)
	if err != nil {
		panic(err)
	}

}
