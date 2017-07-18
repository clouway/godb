package datastoretest

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/clouway/godb"
	"github.com/clouway/godb/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/ory-am/dockertest.v3"
)

const (
	name = "testDb"
)

var (
	db        *DB
	instances int32
)

type DB struct {
	godb.Database

	pool     *dockertest.Pool
	resource *dockertest.Resource
}

// NewDatabase is establishing a new database connection using host
// from the environment. The variable name for the host is TEST_DB_HOST.
// Testing database uses random database name to ensure consistency in tests.
// The created database will be dropped after Clean/Drop function is called.
func NewDatabase() *DB {
	host := os.Getenv("TEST_DB_HOST")
	instCopy := atomic.LoadInt32(&instances)

	if instCopy != 0 {
		atomic.AddInt32(&instances, 1)
		return db
	}

	if host != "" {
		return NewDatabaseWithHost(host)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("mongo", "", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err := mgo.Dial(fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp")))
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	db := NewDatabaseWithHost(fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp")))
	db.pool = pool
	db.resource = resource

	return db
}

// NewDatabaseWithHost is establihing a new database connection
// to the provided host
func NewDatabaseWithHost(host string) *DB {
	instCopy := atomic.LoadInt32(&instances)

	if instCopy != 0 {
		atomic.AddInt32(&instances, 1)
		return db
	}

	dbName := name + strconv.Itoa(time.Now().Nanosecond())

	config := &godb.Config{
		Addrs:    []string{host},
		Database: dbName,
	}

	mgoDB, err := mongo.NewDatabase(config)
	if err != nil {
		panic(fmt.Errorf("could not establish connection: %v", err))
	}

	atomic.AddInt32(&instances, 1)

	db = &DB{
		Database: mgoDB,
	}

	return db
}

// Close closes DB connection
func (db *DB) Close() {
	atomic.AddInt32(&instances, -1)

	instCopy := atomic.LoadInt32(&instances)

	if instCopy != 0 {
		return
	}

	db.Clean()
	db.DropDatabase()
	db.Close()

	if db.resource != nil {
		db.pool.Purge(db.resource)
	}
}

// Clean erases all database collections except system.
func (db *DB) Clean() {
	collections, _ := db.Collections()

	for _, c := range collections {
		c.Clean()
	}
}
