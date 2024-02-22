package links

import (
	"github.com/dzhdmitry/link-shorter/internal/utils"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

var testdata = "./../../test/testdata"

func TestStoreURLs(t *testing.T) {
	_ = os.Remove(testdata + "/results/test_store.csv")
	s, err := NewFileStorage(testdata + "/results/test_store.csv")

	require.NoError(t, err)

	URLs, err := s.StoreURLs([]string{"https://example.com"})

	require.Equal(t, map[string]string{"https://example.com": "1"}, URLs)
	require.NoError(t, err)

	data, err := os.ReadFile(testdata + "/results/test_store.csv")

	require.NoError(t, err)
	require.Equal(t, "1,https://example.com\n", string(data))
}

func TestRestore(t *testing.T) {
	tests := []struct {
		name            string
		filepath        string
		expectedLastKey int64
		expectedLinks   map[int64]string
	}{
		{"Non-existed file", testdata + "/non-existing.csv", 0, map[int64]string{}},
		{"Regular file", testdata + "/test_restore.csv", 3, map[int64]string{
			1: "https://example1.com",
			2: "https://example2.com",
			3: "https://example3.com",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewFileStorage(tt.filepath)

			require.NoError(t, err)

			err = s.Restore()

			require.NoError(t, err)
			require.Equal(t, tt.expectedLastKey, s.lastNumber)
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
		{"Regular", "2", "https://example2.com"},
	}

	s, _ := NewFileStorage(testdata + "/test_restore.csv")
	_ = s.Restore()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := s.GetURL(tt.key)

			require.NoError(t, err)
			require.Equal(t, tt.expected, url)
		})
	}
}

func TestGetURLs(t *testing.T) {
	tests := []struct {
		name         string
		keys         []string
		expectedURLs map[string]string
	}{
		{"Empty", []string{"", ""}, map[string]string{}},
		{"Non-existing", []string{"aawd1"}, map[string]string{}},
		{"Existing", []string{"2", "3"}, map[string]string{
			"2": "https://example2.com",
			"3": "https://example3.com",
		}},
	}

	s, _ := NewFileStorage(testdata + "/test_restore.csv")
	_ = s.Restore()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := s.GetURLs(tt.keys)

			require.NoError(t, err)
			require.Equal(t, tt.expectedURLs, url)
		})
	}
}

func TestAsyncStoreURLs(t *testing.T) {
	_ = os.Remove(testdata + "/results/test_store.csv")
	background := &utils.Background{}
	s, err := NewFileStorageAsync(
		utils.NewLogger(io.Discard, &utils.Clock{}),
		background, testdata+"/results/test_store.csv",
	)

	require.NoError(t, err)

	URLs, err := s.StoreURLs([]string{"https://example.com"})

	require.NoError(t, err)
	require.Equal(t, map[string]string{"https://example.com": "1"}, URLs)

	background.Wait()

	data, err := os.ReadFile(testdata + "/results/test_store.csv")

	require.NoError(t, err)
	require.Equal(t, "1,https://example.com\n", string(data))
}
