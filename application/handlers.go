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

	err = app.Validator.validateURL(data.URL)

	if err != nil {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, err.Error())

		return
	}

	key, err := app.Links.GenerateKey(data.URL)

	if err != nil {
		if errors.Is(err, generator.ErrLimitReached) {
			app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		} else {
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	shortLink := app.composeShortLink(key)
	response, err := app.compactGZIP(app.writeJSON)(w, r, envelope{"link": shortLink})

	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) goHandler(w http.ResponseWriter, r *http.Request) {
	key := httprouter.ParamsFromContext(r.Context()).ByName("key")

	if err := app.Validator.validateKey(key); err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())

		return
	}

	fullLink, err := app.Links.GetURL(key)

	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	if fullLink == "" {
		app.errorResponse(w, r, http.StatusNotFound, "Full link not found for key "+key)

		return
	}

	response, err := app.compactGZIP(app.writeJSON)(w, r, envelope{"link": fullLink})

	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) batchGenerateHandler(w http.ResponseWriter, r *http.Request) {
	var data []string

	err := app.limitMaxBytes(app.extractGZIP(app.readJSON))(w, r, &data)

	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())

		return
	}

	err = app.Validator.validateURLs(data)

	if err != nil {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, err.Error())

		return
	}

	keys, err := app.Links.GenerateKeys(data)

	if err != nil {
		if errors.Is(err, generator.ErrLimitReached) {
			app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		} else {
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	links := map[string]string{}

	for URL, key := range keys {
		links[URL] = app.composeShortLink(key)
	}

	response, err := app.compactGZIP(app.writeJSON)(w, r, envelope{"links": links})

	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) batchGoHandler(w http.ResponseWriter, r *http.Request) {
	var data []string

	err := app.limitMaxBytes(app.extractGZIP(app.readJSON))(w, r, &data)

	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())

		return
	}

	err = app.Validator.validateKeys(data)

	if err != nil {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, err.Error())

		return
	}

	fullLinks := make(map[string]interface{}, len(data))

	for _, key := range data {
		fullLinks[key], err = app.Links.GetURL(key)

		if err != nil {
			app.serverErrorResponse(w, r, err)

			return
		}

		if fullLinks[key] == "" {
			app.errorResponse(w, r, http.StatusNotFound, "Full link not found for key "+key)

			return
		}
	}

	response, err := app.compactGZIP(app.writeJSON)(w, r, envelope{"links": fullLinks})

	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
