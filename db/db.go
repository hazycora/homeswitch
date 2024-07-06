package db

import (
	"crypto/rand"
	"encoding/binary"

	"git.gay/h/homeswitch/models"
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
	Sync()
}

func Sync() {
	Engine.Sync(new(models.Actor))
}

func RandomId() (id uint64, err error) {
	b := make([]byte, 16)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	id = binary.BigEndian.Uint64(b)
	return
}
