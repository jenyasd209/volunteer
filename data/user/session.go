package user

import (
	"fmt"
	"graduate/data"
	"time"
)

type UserType struct {
	Email string
	ID    int
	// CreatedAt time.Time
	Auth bool
}

//Session struct for save user session
type Session struct {
	ID        int
	Email     string
	UserID    int
	CreatrdAt time.Time
}

var User = &UserType{}

//Create new session in DB, when tableName - session table in DB, fieldName - name constrain key field
func (session *Session) Create(tableName string, fieldName string) (err error) {
	statement := `INSERT INTO ` + tableName + ` (email, ` + fieldName + `, created_at) values
                ($1, $2, $3) returning id, email, ` + fieldName + `, created_at`
	stmt, err := data.Db.Prepare(statement)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stmt.Close()
	err = stmt.QueryRow(User.Email, User.ID, time.Now()).Scan(&session.ID,
		&session.Email, &session.UserID, &session.CreatrdAt)

	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func (user *UserType) Set(email string, id int, auth bool) (err error) {
	User.ID = id
	User.Email = email
	// User.CreatedAt = createdAt
	User.Auth = auth
	return
}
