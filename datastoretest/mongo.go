package datastoretest

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/clouway/godb"
	"github.com/clouway/godb/mongo"
	"gopkg.in/mgo.v2"
)

const (
	name       = "testDb"
	mongoImage = "mongo"
)

type DB struct {
	godb.Database

	database    *mgo.Database
	containerID string
}

// NewDatabase is establishing a new database connection using host
// from the environment. The variable name for the host is TEST_DB_HOST.
// Testing database uses random database name to ensure consistency in tests.
// The created database will be dropped after Clean/Drop function is called.
func NewDatabase() *DB {
	host := os.Getenv("TEST_DB_HOST")

	if host != "" {
		return NewDatabaseWithHost(host)
	}

	if _, err := exec.LookPath("docker"); err != nil {
		log.Panicln("skipping without docker available in path")
	}

	if ok, err := dockerHaveImage(mongoImage); !ok || err != nil {
		if err != nil {
			log.Panicf("Error running docker to check for %s: %v", mongoImage, err)
		}
		log.Printf("Pulling docker image %s ...", mongoImage)
		if err := dockerPull(mongoImage); err != nil {
			log.Panicf("Error pulling %s: %v", mongoImage, err)
		}
	}

	containerID, err := dockerRun("-d", mongoImage, "--smallfiles")
	if err != nil {
		log.Panicf("docker run: %v", err)
	}

	ip, err := dockerIP(containerID)

	if err != nil {
		log.Panicf("Error getting container IP: %v", err)
	}

	db := NewDatabaseWithHost(ip)
	db.containerID = containerID

	return db
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

	return &DB{mgoDB, database, ""}
}

// Close closes DB connection
func (db *DB) Close() {
	db.Clean()
	db.database.DropDatabase()
	db.database.Session.Close()

	if db.containerID != "" {
		if err := dockerKillContainer(db.containerID); err != nil {
			log.Panicf("Error killing container %v: %v", db.containerID, err)
		}
	}
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
