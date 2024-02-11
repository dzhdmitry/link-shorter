package links

import (
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"link-shorter.dzhdmitry.net/db"
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

type SQLStorageSuite struct {
	suite.Suite
	db *sql.DB
}

func (s *SQLStorageSuite) SetupSuite() {
	openDB, err := db.Open(db.PrepareTestDB(), 25, 25, "15m")

	if err != nil {
		panic(err)
	}

	s.db = openDB
}

func (s *SQLStorageSuite) SetupTest() {
	_, _ = s.db.Exec("TRUNCATE links")
}

func (s *SQLStorageSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *SQLStorageSuite) TestStoreKeysURLs() {
	tests := []struct {
		name   string
		keyUrl [][]string
	}{
		{"Empty", [][]string{}},
		{"Single row", [][]string{{"0", "http://example.com"}}},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			storage, err := NewSQLStorage(s.db, 1)

			s.NoError(err)

			err = storage.StoreKeysURLs(tt.keyUrl)

			s.NoError(err)
		})
	}
}

func (s *SQLStorageSuite) TestRestore() {
	_, _ = s.db.Exec("INSERT INTO links (key, url) VALUES ('prev', 'prev-url')")
	_, _ = s.db.Exec("INSERT INTO links (key, url) VALUES ('last', 'url')")

	storage, err := NewSQLStorage(s.db, 1)

	s.NoError(err)

	lastKey, err := storage.GetLastKey()

	s.NoError(err)
	s.Equal("last", lastKey)
}

func (s *SQLStorageSuite) TestGetURL() {
	tests := []struct {
		name        string
		key         string
		expectedURL string
	}{
		{"Empty", "", ""},
		{"Non-existing", "aawd1", ""},
		{"Existing", "1q2w", "https://example.com"},
	}

	_, _ = s.db.Exec("INSERT INTO links (key, url) VALUES ('1q2w', 'https://example.com')")
	storage, err := NewSQLStorage(s.db, 1)

	s.NoError(err)

	for _, tt := range tests {
		s.Run(tt.name, func() {
			url, err := storage.GetURL(tt.key)

			s.NoError(err)
			s.Equal(tt.expectedURL, url)
		})
	}
}

func TestSQLStorage(t *testing.T) {
	suite.Run(t, new(SQLStorageSuite))
}
