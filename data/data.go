package data

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

//var Db *sql.DB
//
//func init() {
//	var err error
//	Db, err = sql.Open("postgres", "user=volunteer dbname=volunteer password=qwerty123 sslmode=disable")
//	if err != nil {
//		panic(err)
//	}
//}
//var Db *sql.DB

//func init() {
//	var err error
//	Db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
//	if err != nil {
//		log.Fatalf("Error opening database: %q", err)
//	}
//}
//
//const (
//	host     = "ec2-54-225-89-195.compute-1.amazonaws.com"
//	port     = 5432
//	user     = "cskzypgpkgsquz"
//	password = "67c01a36179e9bcbf616c82c0a1667b147ff53930a67bf1715342268fd878b58"
//	dbname   = "d87sk8skd787p4"
//)
//
//var Db *sql.DB
//
//func init() {
//	var err error
//	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
//		"password=%s dbname=%s sslmode=disable",
//		host, port, user, password, dbname)
//	fmt.Println(psqlInfo)
//	Db, _ = sql.Open("postgres", psqlInfo)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(Db.QueryRow("SELECT * FROM users"))
//
//}

func CreateUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
