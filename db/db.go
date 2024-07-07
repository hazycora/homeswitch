package db

import (
	"crypto/rand"
	"encoding/hex"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var Engine *xorm.Engine

func init() {
	var err error
	Engine, err = xorm.NewEngine("sqlite3", "./test.db")
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
