package container

import (
	"database/sql"
	"errors"
	"link-shorter.dzhdmitry.net/application"
	"link-shorter.dzhdmitry.net/cache"
	"link-shorter.dzhdmitry.net/db"
	"link-shorter.dzhdmitry.net/links"
)

const storageFilename = "tmp/storage.csv"
const storageTypeFile = "file"
const storageTypePostgres = "postgres"

type Container struct {
	Logger     *application.Logger
	Background *application.Background
}

func (c Container) createFileStorage(async bool) (links.StorageInterface, error) {
	if async {
		return links.NewFileStorageAsync(c.Logger, c.Background, storageFilename)
	}

	return links.NewFileStorage(storageFilename)
}

func (c Container) createStorage(config application.Config) (links.StorageInterface, *sql.DB, error) {
	var storage links.StorageInterface
	var err error
	var dbConn *sql.DB

	if config.ProjectStorageType == storageTypeFile {
		storage, err = c.createFileStorage(config.FileAsync)
	} else if config.ProjectStorageType == storageTypePostgres {
		dbConn, err = db.Open(config.DbDSN, config.DbMaxOpenConns, config.DbMaxIdleConns, config.DbMaxIdleTime)

		if err != nil {
			return nil, dbConn, err
		}

		storage, err = links.NewSQLStorage(dbConn, config.DbTimeout)
	} else {
		return nil, nil, errors.New("unknown storage type: " + config.ProjectStorageType)
	}

	if err != nil {
		return nil, dbConn, err
	}

	return storage, dbConn, nil
}

func (c Container) CreateLinksCollection(config application.Config) (application.LinksCollectionInterface, *sql.DB, error) {
	storage, dbConn, err := c.createStorage(config)

	if err != nil {
		return nil, dbConn, err
	}

	var linksCollection application.LinksCollectionInterface
	linksCollection = links.NewCollection(storage)

	if config.CacheEnabled {
		linksCollection = cache.NewCachedCollection(linksCollection, cache.NewLFUCache(config.CacheCapacity))
	}

	return linksCollection, dbConn, nil
}
