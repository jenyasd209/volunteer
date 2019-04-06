package user

type Sessionable interface {
	Check() (bool, error)
	SetUUID(string)
	GetUUID() string
	GetUser() (user Userable, err error)
}

//
// import (
// 	"fmt"
// 	"graduate/data"
// 	"time"
// )
//
// // Session struct for save user session
// type Session struct {
// 	ID        int
// 	UUID      string
// 	Email     string
// 	UserID    int
// 	CreatrdAt time.Time
// }
//
// func (session *Session) Check() (valid bool, err error) {
// 	err = data.Db.QueryRow(`SELECT id, uuid, email, freelancer_id, created_at FROM freelancer_session
// 	                WHERE uuid = $1`, session.UUID).Scan(&session.ID, &session.UUID, &session.Email,
// 		&session.UserID, &session.CreatrdAt)
// 	if err != nil {
// 		valid = false
// 		return
// 	}
// 	if session.ID != 0 {
// 		valid = true
// 	}
// 	return
// }
//
// func (session *Session) GetUser() (freelancer freelancer.Freelancer, err error) {
// 	err = data.Db.QueryRow(`SELECT id, first_name, last_name, email, password, phone,
// 		 											facebook, skype, about, rait, created_at FROM freelancers
// 													WHERE id = $1`, session.UserID).Scan(&freelancer.ID, &freelancer.FirstName,
// 		&freelancer.LastName, &freelancer.Email, &freelancer.Password, &freelancer.Phone,
// 		&freelancer.Facebook, &freelancer.Skype, &freelancer.About, &freelancer.Rait,
// 		&freelancer.CreatedAt)
// 	return
// }
//
// func GetSessionByUUID(uuid string) (session Session) {
// 	statement := `SELECT id, uuid, email, freelancer_id, created_at FROM freelancer_session
// 	                WHERE uuid = $1`
// 	stmt, err := data.Db.Prepare(statement)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
//
// 	defer stmt.Close()
// 	err = stmt.QueryRow(uuid).Scan(&session.ID, &session.UUID, &session.Email,
// 		&session.UserID, &session.CreatrdAt)
//
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
//
// 	return
// }
