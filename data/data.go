package data

import (
	"crypto/sha1"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=volunteer dbname=volunteering password=qwerty123 sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
