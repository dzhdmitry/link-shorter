package application

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ReaderFunc func(http.ResponseWriter, *http.Request, interface{}) error

func (f ReaderFunc) Run(w http.ResponseWriter, r *http.Request, destination interface{}) error {
	return f(w, r, destination)
}

type WriterFunc func(http.ResponseWriter, *http.Request, int, interface{}) ([]byte, error)

func (f WriterFunc) Run(w http.ResponseWriter, r *http.Request, status int, data interface{}) ([]byte, error) {
	return f(w, r, status, data)
}

func (app *Application) logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		app.Logger.LogInfo(message)
		next.ServeHTTP(w, r)
	}
}

func (app *Application) limitMaxBytes(next ReaderFunc) ReaderFunc {
	return func(w http.ResponseWriter, r *http.Request, destination interface{}) error {
		maxBytes := 1_048_576
		r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

		return next.Run(w, r, destination)
	}
}

func (app *Application) extractGZIP(next ReaderFunc) ReaderFunc {
	return func(w http.ResponseWriter, r *http.Request, destination interface{}) error {
		if r.Header.Get("Content-Encoding") != "gzip" {
			return next.Run(w, r, destination)
		}

		var data bytes.Buffer

		reader, err := gzip.NewReader(r.Body)

		if err != nil {
			return err
		}

		_, err = io.Copy(&data, reader)

		defer reader.Close()

		if err != nil {
			return err
		}

		r.Body = io.NopCloser(bytes.NewReader(data.Bytes()))

		return next.Run(w, r, destination)
	}
}

func (app *Application) compactGZIP(next WriterFunc) WriterFunc {
	return func(w http.ResponseWriter, r *http.Request, status int, data interface{}) ([]byte, error) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			return next.Run(w, r, status, data)
		}

		response, err := next.Run(w, r, status, data)
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		gzipWriter := gzip.NewWriter(&buf)
		_, err = gzipWriter.Write(response)

		if err != nil {
			return nil, err
		}

		err = gzipWriter.Close()

		if err != nil {
			return nil, err
		}

		compactResponse := buf.Bytes()

		w.Header().Set("Content-Encoding", "gzip")

		return compactResponse, nil
	}
}
