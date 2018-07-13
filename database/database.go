package database

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2"
)

var Connection *mongodm.Connection

func Initial(models map[string]mongodm.IDocumentBase, test bool) {

	locals := "locals/locals.json"
	database := os.Getenv("DBName")

	if test {
		locals = "../locals/locals.json"
		database = database + "_test"
	}

	file, err := ioutil.ReadFile(locals)

	if err != nil {
		panic(err)
	}

	var localMap map[string]map[string]string
	json.Unmarshal(file, &localMap)

	dbConfig := &mongodm.Config{
		DialInfo: &mgo.DialInfo{
			Addrs:    []string{os.Getenv("DBAddrs")},
			Timeout:  3 * time.Second,
			Database: database,
			Username: os.Getenv("DBUsername"),
			Password: os.Getenv("DBPassword"),
			Source:   database, // os.Getenv("DBSource")
		},
		Locals: localMap["en-US"],
	}

	Connection, err = mongodm.Connect(dbConfig)

	if err != nil {
		panic(err)
	}

	for k, v := range models {
		Connection.Register(v, k)
	}
}
