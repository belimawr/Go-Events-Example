package main

import (
	"log"
	"net/http"

	events "github.com/GuiaBolso/Go-Events"
	"github.com/belimawr/Go-Events-Example/handlers"
)

func main() {
	eventsMux := events.NewMux()

	eventsMux.Add("UpperCase", 1, handlers.UpperCaseHandler{})

	eventsMux.Add("RuneFinder", 1, handlers.RunaHandler{
		Path: "./UCD.db",
	})

	http.Handle("/events/", eventsMux)

	log.Println("Starting example on: 0.0.0.0:3000")

	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		panic(err.Error())
	}
}
