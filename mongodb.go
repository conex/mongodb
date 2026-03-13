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
	Image = "mongo:5"
	// Port used for connecting to MongoDB server.
	Port = "27017"

	// MongoUpWaitTime dictates how long we should wait for MongoDB to accept connections on {{Port}}.
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

// Box returns an mgo.Session and the container running the MongoDB
// instance. It will call t.Fatal on errors.
func Box(t testing.TB, config *Config) (*mgo.Session, conex.Container) {
	if config == nil {
		config = &Config{}
	}

	c := conex.Box(t, &conex.Config{
		Image:  Image,
		Expose: []string{Port},
	})

	config.host = c.Address()
	config.port = Port

	t.Log("Waiting for MongoDB to accept connections")

	if err := c.Wait(Port, MongoUpWaitTime); err != nil {
		c.Drop()
		t.Fatal("MongoDB failed to start:", err)
	}

	t.Log("MongoDB is now accepting connections")

	// Retry connection as MongoDB may need additional time after the port is open
	var db *mgo.Session
	var err error
	for i := 0; i < 10; i++ {
		db, err = mgo.DialWithTimeout(config.url(), 5*time.Second)
		if err == nil {
			break
		}
		t.Logf("MongoDB connection attempt %d failed: %v, retrying...", i+1, err)
		time.Sleep(time.Second)
	}

	if err != nil {
		c.Drop()
		t.Fatal(err)
	}

	return db, c
}
