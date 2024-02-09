package application

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func testExtract(w http.ResponseWriter, r *http.Request, destination interface{}) error {
	return nil
}

func testCompact(w http.ResponseWriter, r *http.Request, data interface{}) ([]byte, error) {
	return []byte("compact"), nil
}

func TestExtractGZIP(t *testing.T) {
	app := Application{}
	w := httptest.NewRecorder()
	file, err := os.ReadFile("./../testdata/compact.gz")

	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/generate", io.NopCloser(bytes.NewReader(file)))
	r.Header.Set("Content-Encoding", "gzip")

	err = app.extractGZIP(testExtract)(w, r, nil)

	require.NoError(t, err)

	body, _ := ioutil.ReadAll(r.Body)

	require.Equal(t, "compact", string(body))
}

func TestCompactGZIP(t *testing.T) {
	app := Application{}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/generate", nil)
	r.Header.Set("Accept-Encoding", "gzip")

	result, err := app.compactGZIP(testCompact)(w, r, nil)

	require.NoError(t, err)

	file, err := os.ReadFile("./../testdata/compact.gz")

	require.NoError(t, err)
	require.Equal(t, file, result)
	require.Equal(t, w.Header().Get("Content-Encoding"), "gzip")
}
