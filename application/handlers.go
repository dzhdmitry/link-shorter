package application

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type envelope map[string]interface{}

func (app *Application) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("link-shorter"))
}

func (app *Application) generateHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		URL string
	}{}

	err := app.readJSON(w, r, &data)

	if err != nil {
		app.errorResponse(w, err.Error())

		return
	}

	if len(data.URL) < 8 {
		app.errorResponse(w, "URL must be present and be at least 8 letters long")

		return
	}

	shortLink := "http://localhost/go/" + data.URL // todo
	err = app.writeJSON(w, http.StatusOK, envelope{"link": shortLink}, nil)

	if err != nil {
		return
	}
}

func (app *Application) goHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	key := params.ByName("key")
	fullLink := "http://example.com/" + key // todo
	err := app.writeJSON(w, http.StatusOK, envelope{"link": fullLink}, nil)

	if err != nil {
		return
	}
}

func (app *Application) errorResponse(w http.ResponseWriter, message string) {
	err := app.writeJSON(w, http.StatusBadRequest, envelope{"error": message}, nil)

	if err != nil {
		app.Logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
