package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type envelope map[string]interface{}

func (app *Application) writeJSON(w http.ResponseWriter, r *http.Request, data interface{}) ([]byte, error) {
	js, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")

	return js, nil
}

func (app *Application) readJSON(w http.ResponseWriter, r *http.Request, destination interface{}) error {
	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	err := decoder.Decode(destination)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshallTypeError *json.UnmarshalTypeError
		var invalidUnmarshallError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshallTypeError):
			if unmarshallTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshallTypeError.Field)
			}

			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshallTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &invalidUnmarshallError):
			panic(err)

		default:
			return err
		}
	}

	err = decoder.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	response, err := app.writeJSON(w, r, envelope{"error": message})

	if err != nil {
		app.Logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(status)
	_, err = w.Write(response)

	if err != nil {
		app.Logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.LogError(err)
	app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
}

func (app *Application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"

	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}
