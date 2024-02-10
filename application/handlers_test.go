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
	"strings"
	"testing"
)

type testLinksStorage struct {
	links   map[int]string
	lastKey int
	maxKey  int
}

func newTestLinkStorage(maxKey int, links map[int]string) *testLinksStorage {
	return &testLinksStorage{
		links:  links,
		maxKey: maxKey,
	}
}

func (t *testLinksStorage) GenerateKey(URL string) (string, error) {
	key := t.lastKey + 1

	if key > t.maxKey {
		return "", generator.ErrLimitReached
	}

	t.links[key] = URL
	t.lastKey = key

	return strconv.Itoa(key), nil
}

func (t *testLinksStorage) GenerateKeys(URLs []string) (map[string]string, error) {
	result := map[string]string{}

	for _, URL := range URLs {
		r, _ := t.GenerateKey(URL)
		result[URL] = r
	}

	return result, nil
}

func (t *testLinksStorage) GetLink(key string) string {
	keyInt, _ := strconv.Atoi(key)

	return t.links[keyInt]
}

func TestIndexHandlerOK(t *testing.T) {
	app := Application{
		Logger: Logger{out: io.Discard},
		Links:  newTestLinkStorage(1, map[int]string{}),
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
		Links:  newTestLinkStorage(1, map[int]string{}),
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
		Links:  newTestLinkStorage(0, map[int]string{}),
	}

	tests := []struct {
		name         string
		request      any
		expectedCode int
		errorMessage string
	}{
		{"Unknown field", envelope{"unknown": "example"}, http.StatusBadRequest, `json: unknown field \"unknown\"`},
		{"Empty body", envelope{}, http.StatusUnprocessableEntity, "URL must be a valid URL string"},
		{"Too big", envelope{"url": strings.Repeat("1", 2001)}, http.StatusUnprocessableEntity, "URL must be maximum 2000 letters long"},
		{"Empty url", envelope{"url": ""}, http.StatusUnprocessableEntity, "URL must be a valid URL string"},
		{"Invalid url #1", envelope{"url": "http"}, http.StatusUnprocessableEntity, "URL must be an absolute URL"},
		{"Invalid url #2", envelope{"url": "http://"}, http.StatusUnprocessableEntity, "URL must be an absolute URL"},
		{"Invalid url #3", envelope{"url": "httpss://exmaple.com"}, http.StatusUnprocessableEntity, "URL must begin with http or https"},
		{"Invalid url #4", envelope{"url": "exmaple.com"}, http.StatusUnprocessableEntity, "URL must be an absolute URL"},
		{"Invalid url #5", envelope{"url": "/exmaple.com"}, http.StatusUnprocessableEntity, "URL must be an absolute URL"},
		{"Over limit", envelope{"url": "http://example.com"}, http.StatusBadRequest, `limit reached`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			body, _ := json.Marshal(tt.request)
			r := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewReader(body))

			app.generateHandler(w, r)

			result := w.Result()

			require.Equal(t, "application/json", result.Header.Get("Content-Type"))
			require.Equal(t, tt.expectedCode, result.StatusCode)

			jsonResponse, err := io.ReadAll(result.Body)

			defer result.Body.Close()

			require.NoError(t, err)
			require.JSONEq(t, `{"error":"`+tt.errorMessage+`"}`+"\n", string(jsonResponse))
		})
	}
}

func TestGoHandlerOK(t *testing.T) {
	app := Application{
		Config: Config{ProjectKeyMaxLength: 5},
		Logger: Logger{out: io.Discard},
		Links: newTestLinkStorage(1, map[int]string{
			1: "https://example.com",
		}),
	}

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
	require.JSONEq(t, `{"link":"https://example.com"}`+"\n", string(jsonResponse))
}

func TestGoHandlerBadRequest(t *testing.T) {
	app := Application{
		Config: Config{ProjectKeyMaxLength: 5},
		Logger: Logger{out: io.Discard},
		Links:  newTestLinkStorage(1, map[int]string{}),
	}

	tests := []struct {
		name         string
		key          string
		errorMessage string
	}{
		{"Empty key", "", "key must be at least 1 letter long"},
		{"Long key", "123456", "key is invalid"},
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

func TestBatchGenerateHandlerOK(t *testing.T) {
	app := Application{
		Logger: Logger{out: io.Discard},
		Links:  newTestLinkStorage(2, map[int]string{}),
	}
	w := httptest.NewRecorder()
	body, _ := json.Marshal([]string{"https://example.org", "https://example2.org"})
	r := httptest.NewRequest(http.MethodPost, "/batch/generate", bytes.NewReader(body))

	app.batchGenerateHandler(w, r)

	result := w.Result()

	require.Equal(t, "application/json", result.Header.Get("Content-Type"))
	require.Equal(t, http.StatusOK, result.StatusCode)

	jsonResponse, err := io.ReadAll(result.Body)

	defer result.Body.Close()

	require.NoError(t, err)
	require.JSONEq(t, `{"links":{"https://example.org":"http://localhost/go/1","https://example2.org":"http://localhost/go/2"}}`+"\n", string(jsonResponse))
}

func TestBatchGoHandlerOK(t *testing.T) {
	app := Application{
		Config: Config{ProjectKeyMaxLength: 5},
		Logger: Logger{out: io.Discard},
		Links: newTestLinkStorage(2, map[int]string{
			1: "https://example.com",
			2: "https://example2.com",
		}),
	}
	w := httptest.NewRecorder()
	body, _ := json.Marshal([]string{"1", "2"})
	r := httptest.NewRequest(http.MethodPost, "/batch/go", bytes.NewReader(body))

	app.batchGoHandler(w, r)

	result := w.Result()

	require.Equal(t, "application/json", result.Header.Get("Content-Type"))
	require.Equal(t, http.StatusOK, result.StatusCode)

	jsonResponse, err := io.ReadAll(result.Body)

	defer result.Body.Close()

	require.NoError(t, err)
	require.JSONEq(t, `{"links":{"1":"https://example.com","2":"https://example2.com"}}`+"\n", string(jsonResponse))
}

func newRequestWithNamedParameter(method, target string, params any) *http.Request {
	r := httptest.NewRequest(method, target, nil)
	ctx := r.Context()
	ctx = context.WithValue(ctx, httprouter.ParamsKey, params)

	r.WithContext(ctx)

	r, _ = http.NewRequestWithContext(ctx, method, target, nil)

	return r
}
