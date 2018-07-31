package database

import (
	"os"
	"testing"

	"github.com/Gommunity/GoWithWith/services/utility"
	"github.com/stretchr/testify/assert"
	"github.com/zebresel-com/mongodm"
)

type testModel struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`
	fake                 string `json:"fake" bson:"fake"`
}

func TestDatabaseConnection(t *testing.T) {

	utility.LoadEnvironmentVariables("../../.env")

	t.Run("LookingForLocalFile", func(t *testing.T) {

		db := &Composer{
			Locals:   "locals.json",
			Addrs:    []string{"fakeAddrs"},
			Database: "",
			Username: "",
			Password: "",
			Source:   "",
		}

		assert.Panics(t, func() { db.Shoot(map[string]mongodm.IDocumentBase{}) })
	})

	t.Run("WrongConnection", func(t *testing.T) {

		db := &Composer{
			Locals:   "../../resource/locals/locals.json",
			Addrs:    []string{"fakeAddrs"},
			Database: "",
			Username: "",
			Password: "",
			Source:   "",
		}

		assert.Panics(t, func() { db.Shoot(map[string]mongodm.IDocumentBase{}) })
	})

	t.Run("SuccessConnectionWithModel", func(t *testing.T) {

		db := &Composer{
			Locals:   "../../resource/locals/locals.json",
			Addrs:    []string{os.Getenv("DBAddrs")},
			Database: os.Getenv("DBName"),
			Username: os.Getenv("DBUsername"),
			Password: os.Getenv("DBPassword"),
			Source:   os.Getenv("DBSource"),
		}

		assert.NotPanics(t, func() {
			db.Shoot(map[string]mongodm.IDocumentBase{
				"test": &testModel{},
			})
		})
	})
}
