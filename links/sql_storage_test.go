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
	_, _ = s.db.Exec("ALTER SEQUENCE links_id_seq RESTART")
}

func (s *SQLStorageSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *SQLStorageSuite) TestStoreURLs() {
	tests := []struct {
		name     string
		urls     []string
		expected map[string]string
	}{
		{"Empty", []string{}, map[string]string{}},
		{"Single row", []string{"https://example.com"}, map[string]string{"https://example.com": "1"}},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			storage := NewSQLStorage(s.db, 1)
			data, err := storage.StoreURLs(tt.urls)

			s.NoError(err)
			s.Equal(tt.expected, data)
		})
	}
}

func (s *SQLStorageSuite) TestGetURL() {
	tests := []struct {
		name        string
		key         string
		expectedURL string
	}{
		{"Empty", "", ""},
		{"Non-existing", "aawd1", ""},
		{"Existing", "1", "https://example.com"},
	}

	_, _ = s.db.Exec("INSERT INTO links (url) VALUES ('https://example.com')")
	storage := NewSQLStorage(s.db, 1)

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
		{"Existing", []string{"1", "2"}, map[string]string{
			"1": "https://example.com",
			"2": "https://example2.com",
		}},
	}

	_, _ = s.db.Exec("INSERT INTO links (url) VALUES ('https://example.com')")
	_, _ = s.db.Exec("INSERT INTO links (url) VALUES ('https://example2.com')")
	storage := NewSQLStorage(s.db, 1)

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
