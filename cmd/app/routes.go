package app

import (
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.indexHandler)
	router.HandlerFunc(http.MethodPost, "/generate", app.logRequest(app.generateHandler))
	router.HandlerFunc(http.MethodGet, "/go/:key", app.logRequest(app.goHandler))
	router.HandlerFunc(http.MethodPost, "/batch/generate", app.logRequest(app.batchGenerateHandler))
	router.HandlerFunc(http.MethodPost, "/batch/go", app.logRequest(app.batchGoHandler))

	router.HandlerFunc(http.MethodGet, "/swagger/:any", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	return app.recoverPanic(app.rateLimit(router))
}
