package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"session-17/database"
	"session-17/handler"
	"session-17/repository"
	"session-17/router"
	"session-17/service"
)

func main() {
	var templates = template.Must(template.New("").ParseGlob("views/**/*.html"))

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

	r := router.NewRouter(handler)

	fmt.Println("server running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("error server")
	}
}
