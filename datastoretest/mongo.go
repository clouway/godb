package datastoretest

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/globalsign/mgo"
	"github.com/ory/dockertest"

	"github.com/clouway/godb"
	"github.com/clouway/godb/mongo"
)

const (
	name = "testDb"

	// image is pinned because the database is reached through globalsign/mgo,
	// which only ever emits the legacy OP_QUERY opcode. MongoDB removed
	// OP_QUERY commands in 5.1, so an unpinned "mongo:latest" resolves to a
	// server this package cannot talk to at all.
	image = "mongo"
	tag   = "4.4"

	// containerTTL makes the docker daemon kill a container that outlived the
	// test run. Purge below is the normal path, but a test binary that panics
	// or is killed never reaches it, and every leaked container holds on to its
	// data volumes.
	containerTTL = 600
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

	resource, err := pool.Run(image, tag, nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := resource.Expire(containerTTL); err != nil {
		log.Printf("could not set the expiration of the container, it will have to be removed by hand if the run is interrupted: %v", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		sess, err := mgo.Dial(fmt.Sprintf("localhost:%s", resource.GetPort("27017/tcp")))
		if err != nil {
			return err
		}
		defer sess.Close()

		return sess.Ping()
	}); err != nil {
		// Without this the container is left running: Fatalf exits the process
		// before any test can reach Close, and nothing else ever purges it.
		pool.Purge(resource)

		log.Fatalf("Could not connect to the database in the container: %s", err)
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
		Addrs:            []string{host},
		Database:         dbName,
		Timeout:          60 * time.Second,
		MaxRetryAttempts: 5,
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

	// Explicitly the embedded database: an unqualified db.Close() resolves back
	// to this method, which re-enters, drops instances to -1 and returns on the
	// guard above without ever closing the session - and leaves the counter
	// negative, so a later NewDatabase hands back this closed instance.
	db.Database.Close()

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
