package application

import (
	"fmt"
	"net/http"
)

func (app *Application) logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		app.Logger.LogInfo(message)
		next.ServeHTTP(w, r)
	}
}
