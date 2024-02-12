package links

import (
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"
	"link-shorter.dzhdmitry.net/db"
	"testing"
)

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

func (s *SQLStorageSuite) TestGetURLs() {
	tests := []struct {
		name         string
		keys         []string
		expectedURLs map[string]string
	}{
		{"Empty", []string{"", ""}, map[string]string{}},
		{"Non-existing", []string{"aawd1"}, map[string]string{}},
		{"Existing", []string{"1q2w", "4hfc8"}, map[string]string{
			"1q2w":  "https://example.com",
			"4hfc8": "https://example2.com",
		}},
	}

	_, _ = s.db.Exec("INSERT INTO links (key, url) VALUES ('1q2w', 'https://example.com')")
	_, _ = s.db.Exec("INSERT INTO links (key, url) VALUES ('4hfc8', 'https://example2.com')")
	storage, err := NewSQLStorage(s.db, 1)

	s.NoError(err)

	for _, tt := range tests {
		s.Run(tt.name, func() {
			url, err := storage.GetURLs(tt.keys)

			s.NoError(err)
			s.Equal(tt.expectedURLs, url)
		})
	}
}

func TestSQLStorage(t *testing.T) {
	suite.Run(t, new(SQLStorageSuite))
}
