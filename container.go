package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"partyinvites/db"
)

type Container struct {
	DB        *sql.DB
	templates map[string]*template.Template
}

func NewContainer(templates map[string]*template.Template) *Container {
	c := new(Container)
	c.templates = templates
	c.initDatabase()
	return c
}

func (c *Container) initDatabase() {
	c.DB = db.ConnectDB()
}

func (c *Container) welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	c.templates["welcome"].Execute(writer, nil)
}

type Guest struct {
	ID       int8
	Name     string
	Email    string
	Phone    string
	UserJoin bool
}

func (c *Container) getAttendingGuests() ([]Guest, error) {
	query := "SELECT id, name, email, phone, userjoin FROM guests WHERE userjoin = true"
	rows, err := c.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendingGuests []Guest

	for rows.Next() {
		var guest Guest
		if err := rows.Scan(&guest.ID, &guest.Name, &guest.Email, &guest.Phone, &guest.UserJoin); err != nil {
			return nil, err
		}
		attendingGuests = append(attendingGuests, guest)
	}

	return attendingGuests, nil
}

func (c *Container) listHandler(writer http.ResponseWriter, request *http.Request) {

	attendingGuests, err := c.getAttendingGuests()
	if err != nil {
		http.Error(writer, "Ошибка при получении данных из базы данных", http.StatusInternalServerError)
		return
	}
	party := make([]*Party, 0, len(attendingGuests))
	for _, s := range attendingGuests {
		party = append(party, &Party{
			Name:       s.Name,
			Email:      s.Email,
			Phone:      s.Phone,
			WillAttend: s.UserJoin,
		})
	}

	err = c.templates["list"].Execute(writer, party)
	if err != nil {
		log.Println(err)
	}

}

func (c *Container) formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		c.templates["form"].Execute(writer, formData{
			Party: &Party{}, Errors: []string{},
		})
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		if request.Form["name"][0] != "" && request.Form["email"][0] != "" && request.Form["phone"][0] != "" {
			c.CreateGuest(request.Form["name"][0], request.Form["email"][0], request.Form["phone"][0], request.Form["willattend"][0] == "true")
		}
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
			c.templates["form"].Execute(writer, formData{
				Party: &responseData, Errors: errors,
			})
		} else {
			responses = append(responses, &responseData)
			if responseData.WillAttend {
				c.templates["thanks"].Execute(writer, responseData.Name)
			} else {
				c.templates["sorry"].Execute(writer, responseData.Name)
			}
		}
	}
}

func (c *Container) CreateGuest(name string, email string, phone string, userjoin bool) {
	sqlStatement := `INSERT INTO guests (name, email, phone, userjoin) 
VALUES ($1, $2, $3, $4)`
	var err error
	_, err = c.DB.Exec(sqlStatement, name, email, phone, userjoin)
	if err != nil {
		panic(err)
	}

}
