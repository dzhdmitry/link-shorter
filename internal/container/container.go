package container

import (
	"database/sql"
	"errors"
	"github.com/redis/go-redis/v9"
	"link-shorter.dzhdmitry.net/cmd/app"
	"link-shorter.dzhdmitry.net/internal/cache"
	"link-shorter.dzhdmitry.net/internal/db"
	"link-shorter.dzhdmitry.net/internal/links"
	"link-shorter.dzhdmitry.net/internal/utils"
)

const storageFilename = "tmp/storage.csv"

type Container struct {
	Logger     *utils.Logger
	Background *utils.Background
}

func (c *Container) createFileStorage(async bool) (links.StorageInterface, error) {
	if async {
		return links.NewFileStorageAsync(c.Logger, c.Background, storageFilename)
	}

	return links.NewFileStorage(storageFilename)
}

func (c *Container) createStorage(config app.Config) (links.StorageInterface, *sql.DB, error) {
	var storage links.StorageInterface
	var dbConn *sql.DB
	var err error

	if config.ProjectStorageType == app.StorageTypeFile {
		storage, err = c.createFileStorage(config.FileAsync)
	} else if config.ProjectStorageType == app.StorageTypePostgres {
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

func (c *Container) createCache(config app.Config) (cache.LinksCacheInterface, *redis.Client, error) {
	var linksCache cache.LinksCacheInterface
	var rdb *redis.Client
	var err error

	if config.CacheType == app.CacheTypeInMemory {
		linksCache = cache.NewLFUCache(config.CacheCapacity)
	} else if config.CacheType == app.CacheTypeRedis {
		rdb, err = db.OpenRedis(config.CacheRedisDSN)

		if err != nil {
			return nil, rdb, err
		}

		linksCache = cache.NewRedisCache(rdb)
	} else if config.CacheType != app.CacheTypeDisabled {
		return nil, nil, errors.New("unknown cache type: " + config.CacheType)
	}

	return linksCache, rdb, nil
}

func (c *Container) CreateLinksCollection(config app.Config) (app.LinksCollectionInterface, *sql.DB, *redis.Client, error) {
	storage, dbConn, err := c.createStorage(config)

	if err != nil {
		return nil, dbConn, nil, err
	}

	var linksCollection app.LinksCollectionInterface
	linksCollection = links.NewCollection(storage)

	if config.CacheType == app.CacheTypeDisabled {
		return linksCollection, dbConn, nil, nil
	}

	linksCache, rdb, err := c.createCache(config)

	if err != nil {
		return nil, dbConn, rdb, err
	}

	linksCollection = cache.NewCachedCollection(linksCollection, linksCache)

	return linksCollection, dbConn, rdb, nil
}
