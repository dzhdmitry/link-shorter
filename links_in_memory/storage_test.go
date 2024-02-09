package links_in_memory

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestStore(t *testing.T) {
	_ = os.Remove("./../testdata/results/test_store.csv")
	fs := FileStorage{"./../testdata/results/test_store.csv"}
	err := fs.Store("test-key", "https://example.com")

	require.NoError(t, err)

	data, err := os.ReadFile("./../testdata/results/test_store.csv")

	require.Equal(t, "test-key,https://example.com\n", string(data))
}

func TestRestore(t *testing.T) {
	links := map[string]string{}
	fs := FileStorage{"./../testdata/test_restore.csv"}
	lastKey, err := fs.Restore(links)

	require.NoError(t, err)

	require.Len(t, links, 3)
	require.Equal(t, "https://example1.com", links["test-key"])
	require.Equal(t, "https://example2.com", links["test-key2"])
	require.Equal(t, "https://example3.com", links["test-key3"])
	require.Equal(t, "test-key3", lastKey)
}
