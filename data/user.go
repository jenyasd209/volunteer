package data

import (
	"fmt"
	"time"
)

//User struct for "users" table
type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Password  string
	Phone     string
	Facebook  string
	Skype     string
	About     string
	Rait      float32
	CreatedAt time.Time
}

//Create new row from "freelancer" table
func (user *User) Create() (err error) {
	statement := `insert into users (first_name, last_name, email, password, phone, facebook, skype, about, rait, created_at)
								values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id, created_at`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(user.FirstName, user.LastName, user.Email, Encrypt(user.Password), user.Phone,
		user.Facebook, user.Skype, user.About, user.Rait, time.Now()).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}
	return
}

// Update row in "freelancer" table
func (user *User) Update() (err error) {
	statement := `UPDATE users SET first_name = $1, last_name = $2, email = $3, password = $4,
	 							phone = $5, facebook = $6, skype = $7, about = $8 WHERE id = $9 returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	err = stmt.QueryRow(user.FirstName, user.LastName, user.Email, Encrypt(user.Password), user.Phone, user.Facebook, user.Skype, user.About, user.ID).Scan()
	if err != nil {
		fmt.Println(err)
	}
	return
}

// Delete row from "freelancer" table
func (user *User) Delete() (err error) {
	statement := "DELETE FROM users WHERE id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	_, err = stmt.Exec(user.ID)
	if err != nil {
		fmt.Println(err)
	}
	return
}

// CreateSession - new session for User
func (user *User) CreateSession() (session Session, err error) {
	statement := `INSERT INTO session (uuid, email, user_id, created_at) values
	                ($1, $2, $3, $4) returning id, uuid, email, user_id, created_at`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(CreateUUID(), user.Email, user.ID, time.Now()).Scan(&session.ID,
		&session.UUID, &session.Email, &session.UserID, &session.CreatrdAt)

	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

//GetAllUsers return all rows from table "freelancer"
func GetAllUsers() (users []User, r error) {
	rows, err := Db.Query(`SELECT id, first_name, last_name, email, password,
															phone, facebook, skype, about, rait, created_at FROM users`)
	if err != nil {
		return
	}

	for rows.Next() {
		user := User{}

		if err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone,
			&user.Facebook, &user.Skype, &user.About, &user.Rait, &user.CreatedAt); err != nil {
			users = append(users, user)
		}
	}

	rows.Close()
	return
}

//GetUserByEmail return rows with required email
func GetUserByEmail(email string) (user User, err error) {
	err = Db.QueryRow(`SELECT id, first_name, last_name, email, password, phone,
		 								 facebook, skype, about, rait, created_at FROM users
										 WHERE email = $1`, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Password, &user.Phone, &user.Facebook, &user.Skype, &user.About, &user.Rait, &user.CreatedAt)

	return
}

//GetUserByID return rows with required ID
func GetUserByID(id int) (user User, err error) {
	err = Db.QueryRow(`SELECT id, first_name, last_name, email, password, phone,
		 								 facebook, skype, about, rait, created_at FROM users
										 WHERE id = $1`, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Password, &user.Phone, &user.Facebook, &user.Skype, &user.About, &user.Rait, &user.CreatedAt)
	return
}
