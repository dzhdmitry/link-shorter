package app

import (
	"github.com/julienschmidt/httprouter"
	//_ "link-shorter.dzhdmitry.net/docs"
	"net/http"
)

// @Summary      Index
// @Description  Does nothing
// @Tags         Default
// @Accept       html
// @Produce      html
// @Router       / [get]
func (app *Application) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("link-shorter"))

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// generateHandler godoc
// @Summary      Generate short link
// @Description  Provide long link and get short one
// @Tags         Single link
// @Accept       json
// @Produce      json
// @Param        request body object{URL=string} true "Original URL"
// @Success      200  {object}  object{link=string}
// @Failure      400  {object}  object{error=string}
// @Failure      422  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /generate [post]
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
		app.serverErrorResponse(w, r, err)

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

// goHandler godoc
// @Summary      Get short link
// @Description  Go by short link and get original url
// @Tags         Single link
// @Accept       json
// @Produce      json
// @Param        key   path string true "Short key"
// @Success      200  {object}  object{links=string}
// @Failure      400  {object}  object{error=string}
// @Failure      422  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /go/{key} [get]
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

// batchGenerateHandler godoc
// @Summary      Generate short links
// @Description  Provide plenty of links and get short url for each
// @Tags         Multiple links
// @Accept       json
// @Produce      json
// @Param        request body []string true "Original URLs"
// @Success      200  {object}  object{links=object{key=string}}
// @Failure      400  {object}  object{error=string}
// @Failure      422  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /batch/generate [post]
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

	links, err := app.Links.GenerateKeys(data)

	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	for URL, key := range links {
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

// batchGoHandler godoc
// @Summary      Get short links
// @Description  Provide short keys and get original url for each
// @Tags         Multiple links
// @Accept       json
// @Produce      json
// @Param        request body []string true "Short keys"
// @Success      200  {object}  object{links=string}
// @Failure      400  {object}  object{error=string}
// @Failure      422  {object}  object{error=string}
// @Failure      404  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /batch/go [get]
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

	fullLinks, err := app.Links.GetURLs(data)

	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
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
