package user

import (
	"fmt"
	"graduate/data"
	"time"
)

type Session struct {
	ID        int
	UUID      string
	Email     string
	UserID    int
	CreatrdAt time.Time
}

type SessionHelper interface {
	SetUUID(string)
	GetUUID() string
	GetUserID() int
	User() (User, error)
	Delete() error
	Check() (bool, error)
}

func (session *Session) Check() (valid bool, err error) {
	err = data.Db.QueryRow(`SELECT id, uuid, email, user_id, created_at FROM session
	                WHERE uuid = $1`, session.UUID).Scan(&session.ID, &session.UUID, &session.Email,
		&session.UserID, &session.CreatrdAt)
	if err != nil {
		valid = false
		return
	}
	if session.ID != 0 {
		valid = true
	}
	return
}

func (session *Session) User() (user User, err error) {
	err = data.Db.QueryRow(`SELECT id, email, password, phone,
		 											facebook, skype, about, rait, created_at FROM users
													WHERE id = $1`, session.UserID).Scan(&user.ID, &user.Email,
		&user.Password, &user.Phone, &user.Facebook, &user.Skype, &user.About, &user.Rait,
		&user.CreatedAt)
	return
}

func (session *Session) GetUserID() (id int) {
	id = session.UserID
	return
}

func (session *Session) SetUUID(uuid string) {
	session.UUID = uuid
	return
}

func (session *Session) GetUUID() (uuid string) {
	uuid = session.UUID
	return
}

func (session *Session) Delete() (err error) {
	return
}

func GetSessionByUUID(uuid string) (session Session) {
	statement := `SELECT id, uuid, email, user_id, created_at FROM user_session
	                WHERE uuid = $1`
	stmt, err := data.Db.Prepare(statement)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(uuid).Scan(&session.ID, &session.UUID, &session.Email,
		&session.UserID, &session.CreatrdAt)

	if err != nil {
		fmt.Println(err)
		return
	}

	return
}
