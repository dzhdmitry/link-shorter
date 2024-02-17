package app

import (
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"net/http/pprof"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.indexHandler)
	router.HandlerFunc(http.MethodPost, "/generate", app.metricsMiddleware(app.logRequest(app.generateHandler)))
	router.HandlerFunc(http.MethodGet, "/go/:key", app.metricsMiddleware(app.logRequest(app.goHandler)))
	router.HandlerFunc(http.MethodPost, "/batch/generate", app.metricsMiddleware(app.logRequest(app.batchGenerateHandler)))
	router.HandlerFunc(http.MethodPost, "/batch/go", app.metricsMiddleware(app.logRequest(app.batchGoHandler)))

	router.HandlerFunc(http.MethodGet, "/swagger/:any", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	for _, v := range []string{"", "allocs", "block", "heap", "threadcreate", "goroutine"} {
		router.HandlerFunc(http.MethodGet, "/debug/pprof/"+v, pprof.Index)
	}

	router.HandlerFunc(http.MethodGet, "/debug/pprof/cmdline", pprof.Cmdline)
	router.HandlerFunc(http.MethodGet, "/debug/pprof/profile", pprof.Profile)
	router.HandlerFunc(http.MethodGet, "/debug/pprof/symbol", pprof.Symbol)
	router.HandlerFunc(http.MethodGet, "/debug/pprof/trace", pprof.Trace)

	router.Handler(http.MethodGet, "/metrics", promhttp.Handler())

	return app.recoverPanic(app.rateLimit(router))
}
