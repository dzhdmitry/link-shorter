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

	err := app.limitMaxBytes(app.extractGZIP(app.readJSON))(w, r, &data)

	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())

		return
	}

	err = validateURL(data.URL)

	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())

		return
	}

	key, err := app.Links.GenerateKey(data.URL)

	if err != nil {
		if errors.Is(err, generator.ErrLimitReached) {
			app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		} else {
			app.Logger.LogError(err)
			app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		}

		return
	}

	shortLink := app.composeShortLink(key)
	response, err := app.compactGZIP(app.writeJSON)(w, r, http.StatusOK, envelope{"link": shortLink})

	if err != nil {
		app.Logger.LogError(err)
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
	}

	_, err = w.Write(response)

	if err != nil {
		app.Logger.LogError(err)
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
	}
}

func (app *Application) goHandler(w http.ResponseWriter, r *http.Request) {
	key := httprouter.ParamsFromContext(r.Context()).ByName("key")

	if key == "" {
		app.errorResponse(w, r, http.StatusBadRequest, "Key must be at least 1 letter long")

		return
	}

	if len(key) > app.Config.ProjectKeyMaxLength {
		app.errorResponse(w, r, http.StatusBadRequest, "Key is invalid")

		return
	}

	fullLink := app.Links.GetLink(key)

	if fullLink == "" {
		app.errorResponse(w, r, http.StatusNotFound, "Full link not found for key "+key)

		return
	}

	response, err := app.compactGZIP(app.writeJSON)(w, r, http.StatusOK, envelope{"link": fullLink})

	if err != nil {
		app.Logger.LogError(err)
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
	}

	_, err = w.Write(response)

	if err != nil {
		app.Logger.LogError(err)
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
	}
}
