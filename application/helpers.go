package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type envelope map[string]interface{}

func (app *Application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	js, err := json.Marshal(data)

	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *Application) readJSON(w http.ResponseWriter, r *http.Request, destination interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
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

func (app *Application) errorResponse(w http.ResponseWriter, status int, message string) {
	err := app.writeJSON(w, status, envelope{"error": message}, nil)

	if err != nil {
		app.Logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
