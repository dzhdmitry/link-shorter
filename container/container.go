package container

import (
	"database/sql"
	"errors"
	"github.com/redis/go-redis/v9"
	"link-shorter.dzhdmitry.net/application"
	"link-shorter.dzhdmitry.net/cache"
	"link-shorter.dzhdmitry.net/db"
	"link-shorter.dzhdmitry.net/links"
)

const storageFilename = "tmp/storage.csv"

type Container struct {
	Logger     *application.Logger
	Background *application.Background
}

func (c *Container) createFileStorage(async bool) (links.StorageInterface, error) {
	if async {
		return links.NewFileStorageAsync(c.Logger, c.Background, storageFilename)
	}

	return links.NewFileStorage(storageFilename)
}

func (c *Container) createStorage(config application.Config) (links.StorageInterface, *sql.DB, error) {
	var storage links.StorageInterface
	var dbConn *sql.DB
	var err error

	if config.ProjectStorageType == application.StorageTypeFile {
		storage, err = c.createFileStorage(config.FileAsync)
	} else if config.ProjectStorageType == application.StorageTypePostgres {
		dbConn, err = db.OpenPostgres(config.DbDSN, config.DbMaxOpenConns, config.DbMaxIdleConns, config.DbMaxIdleTime)

		if err != nil {
			return nil, dbConn, err
		}

		storage = links.NewSQLStorage(dbConn, config.DbTimeout)
	} else {
		return nil, nil, errors.New("unknown storage type: " + config.ProjectStorageType)
	}

	if err != nil {
		return nil, dbConn, err
	}

	return storage, dbConn, nil
}

func (c *Container) createCache(config application.Config) (cache.LinksCacheInterface, *redis.Client, error) {
	var linksCache cache.LinksCacheInterface
	var rdb *redis.Client
	var err error

	if config.CacheType == application.CacheTypeInMemory {
		linksCache = cache.NewLFUCache(config.CacheCapacity)
	} else if config.CacheType == application.CacheTypeRedis {
		rdb, err = db.OpenRedis(config.CacheRedisDSN)

		if err != nil {
			return nil, rdb, err
		}

		linksCache = cache.NewRedisCache(rdb)
	} else if config.CacheType != application.CacheTypeDisabled {
		return nil, nil, errors.New("unknown cache type: " + config.CacheType)
	}

	return linksCache, rdb, nil
}

func (c *Container) CreateLinksCollection(config application.Config) (application.LinksCollectionInterface, *sql.DB, *redis.Client, error) {
	storage, dbConn, err := c.createStorage(config)

	if err != nil {
		return nil, dbConn, nil, err
	}

	var linksCollection application.LinksCollectionInterface
	linksCollection = links.NewCollection(storage)

	if config.CacheType == application.CacheTypeDisabled {
		return linksCollection, dbConn, nil, nil
	}

	linksCache, rdb, err := c.createCache(config)

	if err != nil {
		return nil, dbConn, rdb, err
	}

	linksCollection = cache.NewCachedCollection(linksCollection, linksCache)

	return linksCollection, dbConn, rdb, nil
}
