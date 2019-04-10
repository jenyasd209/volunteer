package user

import (
	"fmt"
	"graduate/data"
	"time"
)

type User struct {
	ID        int
	Email     string
	Password  string
	Phone     string
	Facebook  string
	Skype     string
	About     string
	Rait      float32
	CreatedAt time.Time
}

type UserDBHelper interface {
	Create() error
	Delete() error
	Update() error
	CreateSession() (SessionHelper, error)
}

//
// type UserSessionHelper interface {
// 	CreateSession() (SessionHelper, error)
// }
//
// type UserHelper interface {
// 	UserDBHelper
// 	UserSessionHelper
// }

//Create new row from "freelancer" table
func (user *User) Create() (err error) {
	statement := `insert into users (email, password, phone, facebook, about, rait, created_at)
								values ($1, $2, $3, $4, $5, $6, $7) returning id, created_at`
	stmt, err := data.Db.Prepare(statement)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(user.Email, data.Encrypt(user.Password), user.Phone,
		user.Facebook, user.About, user.Rait, time.Now()).Scan(&user.ID, &user.CreatedAt)

	return
}

func (user *User) Update() (err error) {
	return
}

// Delete row from "freelancer" table
func (user *User) Delete() (err error) {
	statement := "DELETE FROM users WHERE id = $1"
	stmt, err := data.Db.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.ID)
	return
}

func (user *User) CreateSession() (session Session, err error) {
	statement := `INSERT INTO session (uuid, email, user_id, created_at) values
	                ($1, $2, $3, $4) returning id, uuid, email, user_id, created_at`
	stmt, err := data.Db.Prepare(statement)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(data.CreateUUID(), user.Email, user.ID, time.Now()).Scan(&session.ID,
		&session.UUID, &session.Email, &session.UserID, &session.CreatrdAt)

	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

//GetAllUsers return all rows from table "freelancer"
func GetAllUsers() (users []User, r error) {
	rows, err := data.Db.Query(`SELECT id, first_name, last_name, email, password,
															phone, facebook, skype, about, rait, created_at FROM users`)
	if err != nil {
		return
	}

	for rows.Next() {
		user := User{}

		if err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Phone,
			&user.Facebook, &user.Skype, &user.About, &user.Rait, &user.CreatedAt); err != nil {
			users = append(users, user)
		}
	}

	rows.Close()
	return
}

//GetUserByEmail return rows with required email
func GetUserByEmail(email string) (user User, err error) {
	err = data.Db.QueryRow(`SELECT id, email, password, phone, facebook, skype,
		 											about, rait, created_at FROM users
													WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.Password, &user.Phone,
		&user.Facebook, &user.Skype, &user.About, &user.Rait, &user.CreatedAt)

	return
}

//GetUserByID return rows with required ID
func GetUserByID(id int) (user User, err error) {
	err = data.Db.QueryRow(`SELECT id, email, password, phone, facebook, skype,
		 											about, rait, created_at FROM users
													WHERE id = $1`, id).Scan(&user.ID, &user.Email, &user.Password, &user.Phone,
		&user.Facebook, &user.Skype, &user.About, &user.Rait, &user.CreatedAt)
	return
}
