package app

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/felixge/httpsnoop"
	"golang.org/x/time/rate"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ReaderFunc func(http.ResponseWriter, *http.Request, interface{}) error

func (f ReaderFunc) Run(w http.ResponseWriter, r *http.Request, destination interface{}) error {
	return f(w, r, destination)
}

type WriterFunc func(http.ResponseWriter, *http.Request, interface{}) ([]byte, error)

func (f WriterFunc) Run(w http.ResponseWriter, r *http.Request, data interface{}) ([]byte, error) {
	return f(w, r, data)
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
	return func(w http.ResponseWriter, r *http.Request, data interface{}) ([]byte, error) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			return next.Run(w, r, data)
		}

		response, err := next.Run(w, r, data)
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

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()

			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.Config.LimiterEnabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)

			if err != nil {
				app.serverErrorResponse(w, r, err)

				return
			}

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(
						rate.Limit(app.Config.LimiterRPS),
						app.Config.LimiterBurst,
					),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)

				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Application) metricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := httpsnoop.CaptureMetrics(next, w, r)

		MetricRequestsTotal.WithLabelValues(strconv.Itoa(metrics.Code)).Observe(metrics.Duration.Seconds())
	}
}
