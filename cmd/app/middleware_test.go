package app

import (
	"bytes"
	"github.com/dzhdmitry/link-shorter/internal/utils"
	"github.com/dzhdmitry/link-shorter/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testdata = "./../../test/testdata"

func testExtract(w http.ResponseWriter, r *http.Request, destination interface{}) error {
	return nil
}

func testCompact(w http.ResponseWriter, r *http.Request, data interface{}) ([]byte, error) {
	return []byte("compact"), nil
}

func testLog(w http.ResponseWriter, r *http.Request) {
	//
}

func TestExtractGZIP(t *testing.T) {
	app := Application{}
	w := httptest.NewRecorder()
	file, err := os.ReadFile(testdata + "/compact.gz")

	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/generate", io.NopCloser(bytes.NewReader(file)))
	r.Header.Set("Content-Encoding", "gzip")

	err = app.extractGZIP(testExtract)(w, r, nil)

	require.NoError(t, err)

	body, _ := io.ReadAll(r.Body)

	require.Equal(t, "compact", string(body))
}

func TestCompactGZIP(t *testing.T) {
	app := Application{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/generate", nil)
	r.Header.Set("Accept-Encoding", "gzip")

	result, err := app.compactGZIP(testCompact)(w, r, nil)

	require.NoError(t, err)

	file, err := os.ReadFile(testdata + "/compact.gz")

	require.NoError(t, err)
	require.Equal(t, file, result)
	require.Equal(t, w.Header().Get("Content-Encoding"), "gzip")
}

func TestLogRequest(t *testing.T) {
	w := &test.Writer{}
	logger := utils.NewLogger(w, &test.Clock{})
	app := Application{Logger: logger}
	r := httptest.NewRequest(http.MethodGet, "/some_url", nil)
	r.RemoteAddr = "127.0.0.1:1234"

	app.logRequest(testLog)(httptest.NewRecorder(), r)
	assert.Equal(t, "INFO: [2024-02-07T12:00:00Z] 127.0.0.1:1234 - HTTP/1.1 GET /some_url \n", w.Messages[0])
}
