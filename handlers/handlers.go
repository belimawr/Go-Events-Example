package handlers

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	events "github.com/GuiaBolso/Go-Events"
	"github.com/ThoughtWorksInc/runas"
)

// UpperCaseHandler - Event Handler that makes a string uppercase
type UpperCaseHandler struct{}

// Serve - Handler for the UpperCase Event
func (h UpperCaseHandler) Serve(ctx context.Context, event events.Event) (events.Event, error) {
	payload := struct {
		Text string `json:"string"`
	}{}

	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return events.NewError(event.FlowID, err.Error()), err
	}

	response := struct {
		UpperCase string `json:"upper_case"`
	}{}

	response.UpperCase = strings.ToUpper(payload.Text)

	return events.NewResponse(event, response)
}

// RunaHandler - Handler to find unicode characteres
type RunaHandler struct {
	Path string
}

// Serve - Event Handler for "RuneFinder" event
func (h RunaHandler) Serve(ctx context.Context, event events.Event) (events.Event, error) {
	payload := struct {
		Consulta string `json:"consulta"`
	}{}

	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return events.NewError(event.FlowID, err.Error()), err
	}

	ucd, err := runas.AbrirUCD(h.Path)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer ucd.Close()

	consulta := payload.Consulta
	output := runas.Listar(ucd, strings.ToUpper(consulta))

	return events.NewResponse(event, output)
}
