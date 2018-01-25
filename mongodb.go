package mongodb

import (
	"fmt"
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/omeid/conex"
)

var (
	// Image to use for the box.
	Image = "mongo:3"
	// Port used for connect to postgres server.
	Port = "27017"

	// MongoUpWaitTime dectiates how long we should wait for post Postgresql to accept connections on {{Port}}.
	MongoUpWaitTime = 10 * time.Second
)

func init() {
	conex.Require(func() string { return Image })
}

// Config used to connect to the database.
type Config struct {
	Database string

	host string
	port string
}

func (c *Config) url() string {

	url := fmt.Sprintf(
		"mongodb://%s:%s",
		c.host,
		c.port,
	)

	if c.Database != "" {
		return fmt.Sprintf("%s/%s", url, c.Database)
	}

	return url
}

// Box returns an sql.DB connection and the container running the Postgresql
// instance. It will call t.Fatal on errors.
func Box(t testing.TB, config *Config) (*mgo.Session, conex.Container) {
	c := conex.Box(t, &conex.Config{
		Image:  Image,
		Expose: []string{Port},
	})

	config.host = c.Address()
	config.port = Port

	t.Logf("Waiting for MongoDB to accept connections")

	err := c.Wait(Port, MongoUpWaitTime)

	if err != nil {
		c.Drop() // return the container
		t.Fatal("MongoDB failed to start.", err)
	}

	t.Log("MongoDB is now accepting connections")
	db, err := mgo.Dial(config.url())

	if err != nil {
		c.Drop() // return the container
		t.Fatal(err)
	}

	return db, c
}
