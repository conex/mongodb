package mongodb_test

import (
	"os"
	"testing"
	"time"

	"github.com/conex/mongodb"
	"github.com/omeid/conex"
)

func TestMain(m *testing.M) {
	os.Exit(conex.Run(m))
}

func init() {
	mongodb.MongoUpWaitTime = 20 * time.Second
}

func TestPostgres(t *testing.T) {

	sesh, con := mongodb.Box(t, &mongodb.Config{})
	defer con.Drop()

	err := sesh.Ping()
	if err != nil {
		t.Fatal(err)
	}

}
