package admin

import (
	"github.com/zebresel-com/mongodm"
)

const Collection = "Admin"

type Admin struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`

	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
