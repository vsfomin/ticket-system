package app

import (
	"net/http"

	"github.com/fomik2/ticket-system/internal/handlers"

	"github.com/gorilla/mux"
)

func Run(index, layout, editor, auth, user_create, http_port, css_path, database string, repo handlers.RepoInterface) error {
	r := mux.NewRouter()
	handler, err := handlers.New(index, layout, editor, auth, user_create, repo)
	if err != nil {
		return err
	}

	r.HandleFunc("/login", handler.Login).Methods("GET")
	r.HandleFunc("/login", handler.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", handler.LogoutHandler).Methods("GET")

	r.HandleFunc("/api/create", handler.JWTAuthMiddleWare(handler.APICreateTicket)).Methods("POST")
	r.HandleFunc("/api/tickets/{id:[0-9]+}", handler.JWTAuthMiddleWare(handler.APIGetTicket)).Methods("GET")
	r.HandleFunc("/api/signin", handler.APISignin).Methods("POST")

	r.HandleFunc("/", handler.Authentication(handler.WelcomeHandler)).Methods("GET")
	r.HandleFunc("/", handler.Authentication(handler.CreateTicket)).Methods("POST")
	r.HandleFunc("/tickets/{id:[0-9]+}", handler.Authentication(handler.EditHandler)).Methods("POST")
	r.HandleFunc("/tickets/{id:[0-9]+}/delete/", handler.Authentication(handler.DeleteHandler)).Methods("POST")
	r.HandleFunc("/tickets/{id:[0-9]+}", handler.Authentication(handler.GetTicketForEdit)).Methods("GET")
	r.HandleFunc("/user_create/", handler.CreateUserGet).Methods("GET")
	r.HandleFunc("/user_create/", handler.CreateUser).Methods("POST")
	http.Handle("/", r)
	fs := http.FileServer(http.Dir(css_path))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	err = http.ListenAndServe(http_port, nil)
	if err != nil {
		return err
	}
	return nil
}
