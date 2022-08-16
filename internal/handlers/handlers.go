package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/fomik2/ticket-system/internal/entities"
)

type Tickets interface {
	Get(id int) (entities.Ticket, error)
	List() ([]entities.Ticket, error)
	Create(entities.Ticket) (entities.Ticket, error)
	Update(entities.Ticket) (entities.Ticket, error)
	Delete(id int) error
}

/*formData передеается в темплейт при вызове editHandler или welcomeHandler*/
type formData struct {
	entities.Ticket
	Errors     []string
	TicketList []entities.Ticket
}

type Handlers struct {
	indexTempl, editorTempl, ticketsPath, counterPath string
}

func New(config map[string]string) *Handlers {
	newHandler := Handlers{}
	newHandler.indexTempl = config["index"]
	newHandler.editorTempl = config["editor"]
	newHandler.ticketsPath = config["tickets"]
	newHandler.counterPath = config["counter"]
	return &newHandler
}

//getTicketID берет реквест и возвращает ID тикета
func getTicketID(writer http.ResponseWriter, r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Panicln("Не могу распарсить ID тикета", err)
	}
	return id
}

//GetTicketForEdit выбрать заявку для редактирования и показать её
func (h *Handlers) GetTicketForEdit(writer http.ResponseWriter, r *http.Request, repo Tickets) {
	createTemplate, err := template.ParseFiles(h.editorTempl)
	if err != nil {
		log.Panicln("Проблема с загрузкой темплейта", err)
	}
	id := getTicketID(writer, r)
	ticket, err := repo.Get(id)
	if err != nil {
		log.Panicln("Пока что она всегда nil")
	}
	createTemplate.Execute(writer, formData{
		Ticket: ticket, Errors: []string{},
	})
}

//EditHandler редактирование заявки
func (h *Handlers) EditHandler(writer http.ResponseWriter, r *http.Request, repo Tickets) {
	r.ParseForm()
	switch r.Form["action"][0] {
	case "Редактировать":
		id := getTicketID(writer, r)
		currentTicket, err := repo.Get(id)
		if err != nil {
			log.Panicln("Пока что она всегда nil")
		}
		currentTicket.Description = r.Form["description"][0]
		currentTicket.Title = r.Form["title"][0]
		currentTicket.Severity = r.Form["severity"][0]
		_, err = repo.Update(currentTicket)
		if err != nil {
			log.Panicln("Пока что она всегда nil")
		}
		http.Redirect(writer, r, "/", http.StatusSeeOther)
	case "Удалить":
		id := getTicketID(writer, r)
		repo.Delete(id)
		http.Redirect(writer, r, "/", http.StatusSeeOther)
	}
}

//CreateTicket создание новой заявки
func (h *Handlers) CreateTicket(writer http.ResponseWriter, r *http.Request, repo Tickets) {
	createTemplate, err := template.ParseFiles(h.indexTempl)
	if err != nil {
		log.Panicln("Проблема с загрузкой темплейта", err)
	}
	r.ParseForm()
	responseData := entities.Ticket{
		Title:       r.Form["title"][0],
		Description: r.Form["description"][0],
		Severity:    r.Form["severity"][0],
		Status:      "Создана",
		CreatedAt:   time.Now().Local(),
	}
	errors := []string{}
	if responseData.Title == "" {
		errors = append(errors, "Введите название заявки")
	}
	if responseData.Description == "" {
		errors = append(errors, "Введите описание")
	}
	if len(errors) > 0 {
		tickets, err := repo.List()
		if err != nil {
			log.Panicln("Пока что она всегда nil", err)
		}
		createTemplate.Execute(writer, formData{Ticket: responseData, Errors: errors, TicketList: tickets})
	} else {
		repo.Create(responseData)
		http.Redirect(writer, r, "/", http.StatusSeeOther)
	}
}

//welcomeHandler отображение формы и списка всех заявок
func (h *Handlers) WelcomeHandler(writer http.ResponseWriter, r *http.Request, repo Tickets) {
	createTemplate, err := template.ParseFiles(h.indexTempl)
	if err != nil {
		log.Panicln("Проблема с загрузкой темплейта", err)
	}
	tickets, err := repo.List()
	if err != nil {
		log.Panicln("Пока что она всегда nil", err)
	}
	createTemplate.Execute(writer, formData{

		Ticket: entities.Ticket{}, Errors: []string{}, TicketList: tickets,
	})
}
