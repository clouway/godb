package datastoretest

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/clouway/godb"
	"github.com/clouway/godb/mongo"
	"gopkg.in/mgo.v2"
)

const (
	name = "testDb"
)

type DB struct {
	godb.Database
	database *mgo.Database
}

// NewDatabase is establishing a new database connection using host
// from the environment. The variable name for the host is TEST_DB_HOST.
// Testing database uses random database name to ensure consistency in tests.
// The created database will be dropped after Clean/Drop function is called.
func NewDatabase() *DB {
	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		host = "localhost:27017"
	}

	return NewDatabaseWithHost(host)
}

// NewDatabaseWithHost is establihing a new database connection
// to the provided host
func NewDatabaseWithHost(host string) *DB {
	t := time.Now().Nanosecond()
	dbName := name + strconv.Itoa(t)
	config := &godb.Config{
		Addrs:    []string{host},
		Database: dbName,
	}

	mgoDB, err := mongo.NewDatabase(config)
	if err != nil {
		panic(fmt.Errorf("could not establish connection: %v", err))
	}

	sess, err := connect(host)
	if err != nil {
		panic(fmt.Errorf("could not establish connection: %v", err))
	}

	sess.SetMode(mgo.Strong, true)
	sess.SetSocketTimeout(10 * time.Second)

	database := sess.DB(dbName)

	return &DB{mgoDB, database}
}

// Close closes DB connection
func (db *DB) Close() {
	db.Clean()
	db.database.DropDatabase()
	db.database.Session.Close()
}

// Clean erases all database collections except system.
func (db *DB) Clean() {
	cnames, _ := db.database.CollectionNames()

	for _, cname := range cnames {
		db.database.C(cname).RemoveAll(nil)
	}
}

func connect(host string) (*mgo.Session, error) {
	sess, err := mgo.Dial(host)

	if err != nil {
		return nil, err
	}

	return sess, nil
}
