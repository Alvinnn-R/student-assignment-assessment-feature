package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"session-17/database"
	"session-17/handler"
	"session-17/repository"
	"session-17/service"
)

func main() {
	var templates = template.Must(template.New("").ParseGlob("views/*.html"))

	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	for _, t := range templates.Templates() {
		fmt.Println("template:", t.Name())
	}

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service, templates)

	r := http.NewServeMux()

	//view
	r.HandleFunc("GET /login", handler.HandlerAuth.LoginView)
	r.HandleFunc("GET /home", handler.HandlerMenu.HomeView)

	//service
	r.HandleFunc("POST /login", handler.HandlerAuth.Login)

	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/", http.StripPrefix("/public/", fs))

	fmt.Println("server running on port 808")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("error server")
	}
}
