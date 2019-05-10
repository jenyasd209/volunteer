package data

import (
	"fmt"
	"log"
	"time"
)

const (
	UserRoleFreelancer = 1
	UserRoleCustomer   = 2
	UserRoleModerator  = 3
)

type HelperUser interface {
	Create() (err error)
	UpdateInformation() (err error)
	UpdateLoginData() (err error)
	UpdatePhoto() (err error)
	Delete() (err error)
	CreateSession() (session Session, err error)
}

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
	Photo     string
	Rait      float32
	RoleID    int
	CreatedAt time.Time
}

//Role struct for "roles" table
type Role struct {
	ID   int
	Name string
}

//Create new row from "freelancer" table
func (user *User) Create() (err error) {
	statement := `insert into users (first_name, last_name, email, password, phone, facebook, skype, about, rait, role_id, created_at)
								values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id, created_at`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(user.FirstName, user.LastName, user.Email, Encrypt(user.Password), user.Phone,
		user.Facebook, user.Skype, user.About, user.Rait, user.RoleID, time.Now()).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		fmt.Println(err)
	}
	return
}

// UpdateInformation row in "freelancer" table
func (user *User) UpdateInformation() (err error) {
	statement := `UPDATE users SET first_name = $1, last_name = $2,	phone = $3,
	 							facebook = $4, skype = $5, about = $6 WHERE id = $7 returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	err = stmt.QueryRow(user.FirstName, user.LastName, user.Phone, user.Facebook, user.Skype, user.About, user.ID).Scan(&user.ID)
	if err != nil {
		fmt.Println(err)
	}
	return
}

// UpdateLoginData row in "freelancer" table
func (user *User) UpdateLoginData() (err error) {
	statement := `UPDATE users SET email = $1, password = $2 WHERE id = $3 returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	err = stmt.QueryRow(user.Email, Encrypt(user.Password), user.ID).Scan(&user.ID)
	if err != nil {
		fmt.Println(err)
	}
	return
}

// UpdatePhoto column photo_url in "freelancer" table
func (user *User) UpdatePhoto() (err error) {
	statement := `UPDATE users SET photo_url = $1 WHERE id = $2 returning id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()
	err = stmt.QueryRow(user.Photo, user.ID).Scan(&user.ID)
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

func (user *User) IsFreelancer()(result bool)  {
	if user.RoleID == UserRoleFreelancer {
		result = true
	}
	return
}

func (user *User) IsCustomer()(result bool)  {
	if user.RoleID == UserRoleCustomer {
		result = true
	}
	return
}

func (user *User) IsModerator()(result bool)  {
	if user.RoleID == UserRoleModerator {
		result = true
	}
	return
}

//GetAllUsers return all rows from table "freelancer"
func GetAllUsers() (users []User, r error) {
	rows, err := Db.Query(`SELECT id, first_name, last_name, email, password, phone, facebook,
       						photo_url, skype, about, rait, created_at FROM users`)
	if err != nil {
		return
	}

	for rows.Next() {
		user := User{}

		if err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone,
			&user.Facebook, &user.Skype, &user.About, &user.Rait, &user.Photo, &user.CreatedAt); err != nil {
			log.Println(err)
			return
		}
		users = append(users, user)
	}

	rows.Close()
	return
}

//GetUserByEmail return rows with required email
func GetUserByEmail(email string) (user User, err error) {
	err = Db.QueryRow(`SELECT id, first_name, last_name, email, password, phone,
		 								 facebook, skype, about, rait, role_id, photo_url, created_at FROM users
										 WHERE email = $1`, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Password, &user.Phone, &user.Facebook, &user.Skype, &user.About, &user.Rait, &user.RoleID, &user.Photo, &user.CreatedAt)

	return
}

//GetUserByID return rows with required ID
func GetUserByID(id int) (user User, err error) {
	err = Db.QueryRow(`SELECT id, first_name, last_name, email, password, phone,
		 								 facebook, skype, about, rait, role_id, photo_url, created_at FROM users
										 WHERE id = $1`, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Password, &user.Phone, &user.Facebook, &user.Skype, &user.About, &user.Rait, &user.RoleID, &user.Photo, &user.CreatedAt)
	return
}

//UserDeleteAll - delete all rows in table "users"
func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = Db.Exec(statement)
	return
}
