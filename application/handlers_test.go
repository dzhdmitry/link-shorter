package application

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
	"io"
	"link-shorter.dzhdmitry.net/generator"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type testLinksStorage struct {
	links   map[int]string
	lastKey int
	maxKey  int
}

func newTestLinkStorage(maxKey int) *testLinksStorage {
	return &testLinksStorage{
		links:  map[int]string{},
		maxKey: maxKey,
	}
}

func (t testLinksStorage) GenerateKey(URL string) (string, error) {
	key := t.lastKey + 1

	if key > t.maxKey {
		return "", generator.ErrLimitReached
	}

	t.links[key] = URL
	t.lastKey = key

	return strconv.Itoa(key), nil
}

func (t testLinksStorage) GetLink(key string) string {
	keyInt, _ := strconv.Atoi(key)

	return t.links[keyInt]
}

func TestIndexHandlerOK(t *testing.T) {
	app := Application{
		Logger: Logger{out: io.Discard},
		Links:  newTestLinkStorage(1),
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", nil)

	app.indexHandler(w, r)

	result := w.Result()

	require.Equal(t, http.StatusOK, result.StatusCode)

	textResponse, err := io.ReadAll(result.Body)

	defer result.Body.Close()

	require.NoError(t, err)
	require.Equal(t, "link-shorter", string(textResponse))
}

func TestGenerateHandlerOK(t *testing.T) {
	app := Application{
		Logger: Logger{out: io.Discard},
		Links:  newTestLinkStorage(1),
	}
	w := httptest.NewRecorder()
	body, _ := json.Marshal(envelope{"url": "https://example.org"})
	r := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewReader(body))

	app.generateHandler(w, r)

	result := w.Result()

	require.Equal(t, "application/json", result.Header.Get("Content-Type"))
	require.Equal(t, http.StatusOK, result.StatusCode)

	jsonResponse, err := io.ReadAll(result.Body)

	defer result.Body.Close()

	require.NoError(t, err)
	require.JSONEq(t, `{"link":"http://localhost/go/1"}`+"\n", string(jsonResponse))
}

func TestGenerateHandlerBadRequest(t *testing.T) {
	app := Application{
		Logger: Logger{out: io.Discard},
		Links:  newTestLinkStorage(0),
	}

	tests := []struct {
		name         string
		request      any
		errorMessage string
	}{
		{"Unknown field", envelope{"unknown": "example"}, `json: unknown field \"unknown\"`},
		{"Empty body", envelope{}, `URL must be present and be at least 8 letters long`},
		{"Invalid url", envelope{"url": "http://"}, `URL must be present and be at least 8 letters long`},
		{"Over limit", envelope{"url": "http://example.com"}, `limit reached`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tt.request)
			r := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewReader(body))

			app.generateHandler(w, r)

			result := w.Result()

			require.Equal(t, "application/json", result.Header.Get("Content-Type"))
			require.Equal(t, http.StatusBadRequest, result.StatusCode)

			jsonResponse, err := io.ReadAll(result.Body)

			defer result.Body.Close()

			require.NoError(t, err)
			require.JSONEq(t, `{"error":"`+tt.errorMessage+`"}`+"\n", string(jsonResponse))
		})
	}
}

func TestGoHandlerOK(t *testing.T) {
	app := Application{
		Config: Config{
			ProjectKeyMaxLength: 5,
		},
		Logger: Logger{out: io.Discard},
		Links:  newTestLinkStorage(1),
	}

	_, _ = app.Links.GenerateKey("http://example.com")
	w := httptest.NewRecorder()
	r := newRequestWithNamedParameter(http.MethodGet, "/go/:key", httprouter.Params{
		{"key", "1"},
	})

	app.goHandler(w, r)

	result := w.Result()

	require.Equal(t, "application/json", result.Header.Get("Content-Type"))
	require.Equal(t, http.StatusOK, result.StatusCode)

	jsonResponse, err := io.ReadAll(result.Body)

	defer result.Body.Close()

	require.NoError(t, err)
	require.JSONEq(t, `{"link":"http://example.com"}`+"\n", string(jsonResponse))
}

func TestGoHandlerBadRequest(t *testing.T) {
	app := Application{
		Config: Config{
			ProjectKeyMaxLength: 5,
		},
		Logger: Logger{out: io.Discard},
		Links:  newTestLinkStorage(1),
	}

	tests := []struct {
		name         string
		key          string
		errorMessage string
	}{
		{"Empty key", "", "Key must be at least 1 letter long"},
		{"Long key", "123456", "Key is invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := newRequestWithNamedParameter(http.MethodGet, "/go/:key", httprouter.Params{
				{"key", tt.key},
			})

			app.goHandler(w, r)

			result := w.Result()

			require.Equal(t, "application/json", result.Header.Get("Content-Type"))
			require.Equal(t, http.StatusBadRequest, result.StatusCode)

			jsonResponse, err := io.ReadAll(result.Body)

			defer result.Body.Close()

			require.NoError(t, err)
			require.JSONEq(t, `{"error":"`+tt.errorMessage+`"}`+"\n", string(jsonResponse))
		})
	}
}

func newRequestWithNamedParameter(method, target string, params any) *http.Request {
	r := httptest.NewRequest(method, target, nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, httprouter.ParamsKey, params)

	r.WithContext(ctx)

	r, _ = http.NewRequestWithContext(ctx, method, target, nil)

	return r
}
