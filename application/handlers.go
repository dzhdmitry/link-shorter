package application

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"link-shorter.dzhdmitry.net/generator"
	"net/http"
)

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
		app.errorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	if len(data.URL) < 8 {
		app.errorResponse(w, http.StatusBadRequest, "URL must be present and be at least 8 letters long")

		return
	}

	key, err := app.Links.GenerateKey(data.URL)

	if err != nil {
		if errors.Is(err, generator.ErrLimitReached) {
			app.errorResponse(w, http.StatusBadRequest, err.Error())
		} else {
			app.Logger.LogError(err)
			app.errorResponse(w, http.StatusInternalServerError, err.Error())
		}

		return
	}

	shortLink := app.composeShortLink(key)
	err = app.writeJSON(w, http.StatusOK, envelope{"link": shortLink}, nil)

	if err != nil {
		app.Logger.LogError(err)
		app.errorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func (app *Application) goHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	key := params.ByName("key")

	if key == "" {
		app.errorResponse(w, http.StatusBadRequest, "Key must be at least 1 letter long")

		return
	}

	if len(key) > app.Config.ProjectKeyMaxLength {
		app.errorResponse(w, http.StatusBadRequest, "Key is invalid")

		return
	}

	fullLink := app.Links.GetLink(key)

	if fullLink == "" {
		app.errorResponse(w, http.StatusNotFound, "Full link not found for key "+key)

		return
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"link": fullLink}, nil)

	if err != nil {
		app.Logger.LogError(err)
		app.errorResponse(w, http.StatusInternalServerError, err.Error())
	}
}
