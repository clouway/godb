package datastoretest

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/mgo.v2"

	"github.com/clouway/godb"
	"github.com/clouway/godb/mongo"
)

const (
	name = "testDb"
)

type DB struct {
	godb.Database
	database *mgo.Database
}

// NewDatabase is establishing a new database connection using host from the environment.
// The variable name for the host is TEST_DB_HOST.
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
	config := &godb.Config{
		Addrs:    []string{host},
		Database: name,
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
	database := sess.DB(name)

	return &DB{mgoDB, database}
}

// Close closes DB connection
func (db *DB) Close() {
	db.database.Session.Close()
}

// Clean erases all database collections except system.
func (db *DB) Clean() {
	col, err := db.database.CollectionNames()
	if err != nil {
		panic(fmt.Errorf("could not get collection names: %v", err))
	}

	for _, v := range col {
		if !strings.Contains(v, "system.") {
			err := db.database.C(v).DropCollection()

			if err != nil {
				panic(fmt.Errorf("could not drop collection '%s': %v", v, err))
			}
		}
	}
}

func connect(host string) (*mgo.Session, error) {
	sess, err := mgo.Dial(host)

	if err != nil {
		return nil, err
	}

	return sess, nil
}
