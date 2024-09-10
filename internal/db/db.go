package db

import (
	"crypto/rand"
	"encoding/hex"

	"git.gay/h/homeswitch/internal/config"
	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

var Engine *xorm.Engine

func init() {
	var err error
	Engine, err = xorm.NewEngine("postgres", config.DBConnectionUri)
	if err != nil {
		panic(err)
	}
}

func RandomId() (id string, err error) {
	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	id = hex.EncodeToString(b)
	return
}
