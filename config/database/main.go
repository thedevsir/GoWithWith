package database

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2"
)

var Connection *mongodm.Connection

type Composer struct {
	Locals   string
	Addrs    []string
	Database string
	Username string
	Password string
	Source   string
}

func (c Composer) Shoot(models map[string]mongodm.IDocumentBase) {

	file, err := ioutil.ReadFile(c.Locals)

	if err != nil {
		panic(err)
	}

	var localMap map[string]map[string]string
	json.Unmarshal(file, &localMap)

	dbConfig := &mongodm.Config{
		DialInfo: &mgo.DialInfo{
			Addrs:    c.Addrs,
			Timeout:  3 * time.Second,
			Database: c.Database,
			Username: c.Username,
			Password: c.Password,
			Source:   c.Source,
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
