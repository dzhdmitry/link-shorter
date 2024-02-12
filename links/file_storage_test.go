package links

import (
	"github.com/stretchr/testify/require"
	"link-shorter.dzhdmitry.net/application"
	"os"
	"testing"
)

func TestStoreKeysURLs(t *testing.T) {
	_ = os.Remove("./../testdata/results/test_store.csv")
	s, err := NewFileStorage("./../testdata/results/test_store.csv")

	require.NoError(t, err)

	err = s.StoreKeysURLs([][]string{{"test-key", "https://example.com"}})

	require.NoError(t, err)

	data, err := os.ReadFile("./../testdata/results/test_store.csv")

	require.Equal(t, "test-key,https://example.com\n", string(data))

}

func TestRestore(t *testing.T) {
	tests := []struct {
		name            string
		filepath        string
		expectedLastKey string
		expectedLinks   map[string]string
	}{
		{"Non-existed file", "./../testdata/non-existing.csv", "", map[string]string{}},
		{"Regular file", "./../testdata/test_restore.csv", "test-key3", map[string]string{
			"test-key":  "https://example1.com",
			"test-key2": "https://example2.com",
			"test-key3": "https://example3.com",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewFileStorage(tt.filepath)

			require.NoError(t, err)

			err = s.Restore()

			require.NoError(t, err)
			require.Equal(t, tt.expectedLastKey, s.lastKey)
			require.Equal(t, tt.expectedLinks, s.links)
		})
	}

}

func TestGetURL(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"Empty", "", ""},
		{"Non-existed", "1d32g", ""},
		{"Regular", "test-key2", "https://example2.com"},
	}

	s, _ := NewFileStorage("./../testdata/test_restore.csv")
	_ = s.Restore()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := s.GetURL(tt.key)

			require.NoError(t, err)
			require.Equal(t, tt.expected, url)
		})
	}
}

func TestGetLastKey(t *testing.T) {
	s, _ := NewFileStorage("./../testdata/test_restore.csv")
	_ = s.Restore()

	url, err := s.GetLastKey()

	require.NoError(t, err)
	require.Equal(t, "test-key3", url)
}

func TestAsyncStoreKeysURLs(t *testing.T) {
	_ = os.Remove("./../testdata/results/test_store.csv")
	background := &application.Background{}
	s, err := NewFileStorageAsync(application.Logger{}, background, "./../testdata/results/test_store.csv")

	require.NoError(t, err)

	err = s.StoreKeysURLs([][]string{{"test-key", "https://example.com"}})

	require.NoError(t, err)

	background.Wait()

	data, err := os.ReadFile("./../testdata/results/test_store.csv")

	require.Equal(t, "test-key,https://example.com\n", string(data))

}
